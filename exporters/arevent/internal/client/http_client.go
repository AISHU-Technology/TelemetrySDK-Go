package client

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporters/arevent/internal/common"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporters/arevent/internal/config"
	customErrors "devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporters/arevent/internal/errors"
)

// HttpClient 客户端结构体。
type HttpClient struct {
	cfg       config.HTTPConfig
	retryFunc config.RetryFunc
	client    *http.Client
	stopCh    chan struct{}
}

// Stop 关闭发送器。
func (d *HttpClient) Stop(ctx context.Context) error {
	close(d.stopCh)
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	return nil
}

// UploadTraces 批量发送Trace数据。
func (d *HttpClient) UploadEvents(ctx context.Context, events []*common.AREvent) error {
	//退出逻辑：
	ctx, cancel := d.contextWithStop(ctx)
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	defer cancel()

	//编码Trace数据。
	file := bytes.NewBuffer([]byte{})
	encoder := json.NewEncoder(file)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "\t")
	_ = encoder.Encode(events)
	request, err := d.newRequest(file.Bytes())
	if err != nil {
		return err
	}

	//发送HTTP请求。
	requestFunc := d.retryFunc(ctx, func(ctx context.Context) error {
		// 外部驱动终止时，此处退出发送数据。
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		request.reset(ctx)
		//真正发送http.Request的地方
		resp, err := send(d, request.Request)
		//resp, err := d.client.Do(request.Request)
		if err != nil {
			return err
		}

		var rErr error
		switch resp.StatusCode {
		// 发送成功，不重发。
		case http.StatusNoContent, http.StatusOK:
		// 格式校验不通过，不重发。
		case http.StatusBadRequest:
			rErr = errors.New(customErrors.AnyRobotEventExporter_InvalidFormat)
		// 接收器地址不正确，不重发。
		case http.StatusNotFound:
			rErr = errors.New(customErrors.AnyRobotEventExporter_JobIdNotFound)
		// Trace太长超过5MB，不重发。
		case http.StatusRequestEntityTooLarge:
			rErr = errors.New(customErrors.AnyRobotEventExporter_PayloadTooLarge)
		// 网络错误，使用~可重发错误~来管理重发机制。
		case http.StatusTooManyRequests, http.StatusInternalServerError, http.StatusServiceUnavailable:
			rErr = newResponseError(resp.Header)
			if _, err := io.Copy(ioutil.Discard, resp.Body); err != nil {
				_ = resp.Body.Close()
				return err
			}
		default:
			rErr = errors.New(customErrors.AnyRobotEventExporter_Unsent)
		}
		if err := resp.Body.Close(); err != nil {
			return err
		}
		return rErr
	})
	return requestFunc
}

const ENCODING = "Content-Encoding"

// newRequest POST /api/feed_ingester/v1/jobs/{job_id}/events。
func (d *HttpClient) newRequest(body []byte) (arRequest, error) {
	u := url.URL{Scheme: d.getScheme(), Host: d.cfg.Endpoint, Path: d.cfg.Path}
	r, err := http.NewRequest(http.MethodPost, u.String(), nil)
	if err != nil {
		return arRequest{Request: r}, err
	}

	for k, v := range d.cfg.Headers {
		r.Header.Set(k, v)
	}
	//设置来源记录。
	r.Header.Set("Service-Language", "Golang")

	req := arRequest{Request: r}
	//是否使用压缩。
	switch d.cfg.Compression {
	case config.NoCompression:
		r.Header.Set(ENCODING, "json")
		r.ContentLength = (int64)(len(body))
		req.bodyReader = bodyReader(body)
	case config.GzipCompression:
		// 使用Gzip压缩关闭ContentLength。
		r.ContentLength = -1
		r.Header.Set(ENCODING, "gzip")

		gz := gzPool.Get().(*gzip.Writer)
		defer gzPool.Put(gz)

		var b bytes.Buffer
		gz.Reset(&b)

		if _, err := gz.Write(body); err != nil {
			r.Header.Set(ENCODING, "json")
			return req, err
		}
		if err := gz.Close(); err != nil {
			r.Header.Set(ENCODING, "json")
			return req, err
		}

		req.bodyReader = bodyReader(b.Bytes())
	}

	return req, nil
}

func send(d *HttpClient, req *http.Request) (*http.Response, error) {
	//resp, err := d.client.Do(request.Request)
	return d.client.Do(req)
}

// gzPool Gzip压缩流。
var gzPool = sync.Pool{
	New: func() interface{} {
		w := gzip.NewWriter(ioutil.Discard)
		return w
	},
}

// contextWithStop 把上下文停止信号传递给客户端，驱动Exporter停止。
func (d *HttpClient) contextWithStop(ctx context.Context) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	go func(ctx context.Context, cancel context.CancelFunc) {
		select {
		case <-ctx.Done():
		case <-d.stopCh:
			cancel()
		}
	}(ctx, cancel)
	return ctx, cancel
}

// ourTransport 根据net/http自定义连接方式。
var ourTransport = &http.Transport{
	Proxy: http.ProxyFromEnvironment,
	DialContext: (&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}).DialContext,
	ForceAttemptHTTP2:     true,
	MaxIdleConns:          100,
	IdleConnTimeout:       90 * time.Second,
	TLSHandshakeTimeout:   10 * time.Second,
	TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
	ExpectContinueTimeout: 1 * time.Second,
}

// getScheme 决定通过http或者https发送。
func (d *HttpClient) getScheme() string {
	if d.cfg.Insecure {
		return "http"
	}
	return "https"
}

// arRequest 包了一层可重置的body reader。
type arRequest struct {
	*http.Request

	// bodyReader 发送同一请求，用于重发机制。
	bodyReader func() io.ReadCloser
}

// reset 重置请求参数。
func (r *arRequest) reset(ctx context.Context) {
	r.Body = r.bodyReader()
	r.Request = r.Request.WithContext(ctx)
}

// newResponseError 返回一个retryableError。
func newResponseError(header http.Header) error {
	var rErr config.RetryableError
	if s, ok := header["Retry-After"]; ok {
		if t, err := strconv.ParseInt(s[0], 10, 64); err == nil {
			rErr.Throttle = t
		}
	}
	return rErr
}

// bodyReader 返回字节流的读写体。
func bodyReader(buf []byte) func() io.ReadCloser {
	return func() io.ReadCloser {
		return ioutil.NopCloser(bytes.NewReader(buf))
	}
}

// NewHTTPClient 创建Exporter的HTTP客户端。
func NewHTTPClient(opts ...config.Option) Client {
	cfg := config.NewConfig(opts...)

	client := &http.Client{
		Transport: ourTransport,
		Timeout:   cfg.HTTPConfig.Timeout,
	}

	return &HttpClient{
		cfg:       cfg.HTTPConfig,
		retryFunc: cfg.RetryConfig.RetryFunc(),
		stopCh:    make(chan struct{}),
		client:    client,
	}
}

// UploadTraces 批量发送Trace数据。
func (d *HttpClient) UploadEvent(ctx context.Context, event *common.AREvent) error {
	//退出逻辑：
	ctx, cancel := d.contextWithStop(ctx)
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	defer cancel()

	//编码Trace数据。
	file := bytes.NewBuffer([]byte{})
	encoder := json.NewEncoder(file)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "\t")
	_ = encoder.Encode(event)
	request, err := d.newRequest(file.Bytes())
	if err != nil {
		return err
	}

	//发送HTTP请求。
	requestFunc := d.retryFunc(ctx, func(ctx context.Context) error {
		// 外部驱动终止时，此处退出发送数据。
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		request.reset(ctx)
		//真正发送http.Request的地方
		resp, err := send(d, request.Request)
		//resp, err := d.client.Do(request.Request)
		if err != nil {
			return err
		}

		var rErr error
		switch resp.StatusCode {
		// 发送成功，不重发。
		case http.StatusNoContent, http.StatusOK:
		// 格式校验不通过，不重发。
		case http.StatusBadRequest:
			rErr = errors.New(customErrors.AnyRobotEventExporter_InvalidFormat)
		// 接收器地址不正确，不重发。
		case http.StatusNotFound:
			rErr = errors.New(customErrors.AnyRobotEventExporter_JobIdNotFound)
		// Trace太长超过5MB，不重发。
		case http.StatusRequestEntityTooLarge:
			rErr = errors.New(customErrors.AnyRobotEventExporter_PayloadTooLarge)
		// 网络错误，使用~可重发错误~来管理重发机制。
		case http.StatusTooManyRequests, http.StatusInternalServerError, http.StatusServiceUnavailable:
			rErr = newResponseError(resp.Header)
			if _, err := io.Copy(ioutil.Discard, resp.Body); err != nil {
				_ = resp.Body.Close()
				return err
			}
		default:
			rErr = errors.New(customErrors.AnyRobotEventExporter_Unsent)
		}
		if err := resp.Body.Close(); err != nil {
			return err
		}
		return rErr
	})
	return requestFunc
}
