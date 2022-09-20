package config

import (
	"crypto/tls"
	"time"
)

// Compression 定义了Trace压缩方式。
type Compression int

const (
	// NoCompression 0：不压缩。
	NoCompression Compression = iota
	// GzipCompression 1：Gzip压缩。
	GzipCompression
)

type HTTPConfig struct {
	Insecure    bool
	Endpoint    string
	Path        string
	Compression Compression
	Timeout     time.Duration
	Headers     map[string]string
	TLSCfg      *tls.Config
}

// HTTPOption HTTP client专用的配置项结构体。
type HTTPOption struct {
	fn func(Config) Config
}

func (h *HTTPOption) ApplyHTTPOption(cfg Config) Config {
	return h.fn(cfg)
}
func newHTTPOption(fn func(cfg Config) Config) HTTPOption {
	return HTTPOption{fn: fn}
}

// 发送器配置项目。

// WithScheme 选择http/https。
func WithScheme(scheme string) HTTPOption {
	if scheme == "http" {
		return withInsecure()
	}
	return withSecure()
}

// WithEndpoint 设置ip:port。
func WithEndpoint(endpoint string) HTTPOption {
	return newHTTPOption(func(cfg Config) Config {
		cfg.HTTPConfig.Endpoint = endpoint
		return cfg
	})
}

// WithPath 设置/path。
func WithPath(urlPath string) HTTPOption {
	return newHTTPOption(func(cfg Config) Config {
		cfg.HTTPConfig.Path = urlPath
		return cfg
	})
}

// WithCompression 设置压缩方式。
func WithCompression(compression Compression) HTTPOption {
	return newHTTPOption(func(cfg Config) Config {
		cfg.HTTPConfig.Compression = compression
		return cfg
	})
}

// WithRetry 设置重发。
func WithRetry(rc RetryConfig) HTTPOption {
	return newHTTPOption(func(cfg Config) Config {
		cfg.RetryConfig = rc
		return cfg
	})
}

// WithTLSClientConfig 设置TLS连接。
func WithTLSClientConfig(tlsCfg *tls.Config) HTTPOption {
	return newHTTPOption(func(cfg Config) Config {
		cfg.HTTPConfig.TLSCfg = tlsCfg.Clone()
		return cfg
	})
}

// withInsecure 设置为http。
func withInsecure() HTTPOption {
	return newHTTPOption(func(cfg Config) Config {
		cfg.HTTPConfig.Insecure = true
		return cfg
	})
}

// withSecure 设置为https。
func withSecure() HTTPOption {
	return newHTTPOption(func(cfg Config) Config {
		cfg.HTTPConfig.Insecure = false
		return cfg
	})
}

// WithHeader 设置请求头。
func WithHeader(headers map[string]string) HTTPOption {
	return newHTTPOption(func(cfg Config) Config {
		cfg.HTTPConfig.Headers = headers
		return cfg
	})
}

// WithTimeout 设置连接超时时间。
func WithTimeout(duration time.Duration) HTTPOption {
	return newHTTPOption(func(cfg Config) Config {
		cfg.HTTPConfig.Timeout = duration
		return cfg
	})
}
