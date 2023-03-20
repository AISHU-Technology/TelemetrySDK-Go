package encoder

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"reflect"
	"strconv"
	"time"
	"unsafe"

	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/span/exporter"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/span/field"
)

var (
	maxWaitExporterTime = 20 * time.Second

	_quotationSafe   = []byte("\\\"")
	_reverseSafe     = []byte("\\\\")
	_backspaceSafe   = []byte("\\b")
	_formfeedSafe    = []byte("\\f")
	_horizontalSafe  = []byte("\\t")
	_lineFeedSafe    = []byte("\\n")
	_carriageSafe    = []byte("\\r")
	_quotation       = []byte("\"")
	_leftBracket     = []byte("[")
	_rightBracket    = []byte("]")
	_seperator       = []byte(",")
	_lineFeed        = []byte("\n")
	_leftBigBracket  = []byte("{")
	_rightBigBracket = []byte("}")
	_colon           = []byte(": ")
)

type Encoder interface {
	Write(field.Field) error
	Close() error
}

// SyncEncoder 同步模式专用，只能上报到一个地址。
type SyncEncoder interface {
	Write(field.Field) error
	Close() error
	sync()
}

// JsonEncoder JSON编码器。
type JsonEncoder struct {
	w            io.Writer
	logExporters map[string]exporter.LogExporter
	buf          io.Writer
	End          []byte
	bufReal      *bytes.Buffer
	ctx          context.Context
	cancelFunc   context.CancelFunc
}

// NewJsonEncoder 创建2.5.0及之前的版本的JSON编码器。
func NewJsonEncoder(w io.Writer) Encoder {
	ctx, cancel := context.WithCancel(context.Background())
	res := &JsonEncoder{
		w:            w,
		logExporters: nil,
		End:          _lineFeed,
		bufReal:      bytes.NewBuffer(make([]byte, 0, 4096)),
		ctx:          ctx,
		cancelFunc:   cancel,
	}
	res.buf = res.bufReal
	return res
}

// NewJsonEncoderBench 创建2.5.0及之前的版本的JSON编码器。
func NewJsonEncoderBench(w io.Writer) Encoder {
	ctx, cancel := context.WithCancel(context.Background())
	res := &JsonEncoder{
		w:            w,
		logExporters: nil,
		End:          _lineFeed,
		bufReal:      bytes.NewBuffer(make([]byte, 0, 4096)),
		buf:          ioutil.Discard,
		ctx:          ctx,
		cancelFunc:   cancel,
	}
	return res
}

// NewJsonEncoderWithExporters 创建2.6.0后的异步发送模式JSON编码器。
func NewJsonEncoderWithExporters(exporters ...exporter.LogExporter) Encoder {
	ctx, cancel := context.WithCancel(context.Background())
	eps := make(map[string]exporter.LogExporter)
	for _, e := range exporters {
		eps[e.Name()] = e
	}
	res := &JsonEncoder{
		w:            nil,
		logExporters: eps,
		End:          _lineFeed,
		bufReal:      bytes.NewBuffer(make([]byte, 0, 4096)),
		ctx:          ctx,
		cancelFunc:   cancel,
	}
	res.buf = res.bufReal
	return res
}

// NewSyncEncoder 创建2.6.0后的同步发送模式JSON编码器。
func NewSyncEncoder(exporter_ exporter.LogExporter) SyncEncoder {
	ctx, cancel := context.WithCancel(context.Background())
	eps := make(map[string]exporter.LogExporter, 1)
	eps[exporter_.Name()] = exporter_
	res := &JsonEncoder{
		w:            nil,
		logExporters: eps,
		End:          _lineFeed,
		bufReal:      bytes.NewBuffer(make([]byte, 0, 4096)),
		ctx:          ctx,
		cancelFunc:   cancel,
	}
	res.buf = res.bufReal
	return res
}

// Write 区分本地输出和上报AnyRobot，本地输出仅提取关键信息输出，上报AnyRobot为数组形式，有返回值error判断是否发送成功。
func (js *JsonEncoder) Write(f field.Field) error {
	//判断用户是否使用原来2.5.0之前的io.Writer实现输出，如果是的话由于批量发送，为保证与原来数据模型一致将数组遍历单个输出
	if js.w != nil {
		//将数组遍历使用io.Writer单个输出
		js.dealIoWriter(f)
		return nil
	}
	//判断用户是否使用2.6.0之后的exporter实现输出
	if js.logExporters != nil && len(js.logExporters) != 0 {
		//判断用户是否使用新增的StdoutExporter实现标准输出，如果是的话由于批量发送，为保证与原来数据模型一致将数组遍历单个输出
		stdoutExporter, ok := js.logExporters["StdoutExporter"]
		if ok {
			//将数组遍历使用StdoutExporter单个输出，并提取关键信息输出。
			js.dealStdoutExporter(stdoutExporter, f)
		}
		//不是StdoutExporter的剩余exporter比如arExporter将输出整个数组，下面将整个数组转为byte
		err := js.write(f)
		if err != nil {
			log.Println(field.GenerateSpecificError(err))
		}
		// _, res := js.WriteBytes(js.End)
		_, writeBytesErr := js.WriteBytes(js.End)
		if writeBytesErr != nil {
			log.Println(field.GenerateSpecificError(writeBytesErr))
		}
		//调用不是StdoutExporter的剩余exporter比如arExporter的输出方法将整个数组进行输出
		flushWithExportersErr := js.flushWithExporters()
		if flushWithExportersErr != nil {
			log.Println(field.GenerateSpecificError(flushWithExportersErr))
			return flushWithExportersErr
		}
	}
	return nil
}

// Close 关闭JSON编码器。
func (js *JsonEncoder) Close() error {
	if js.bufReal.Len() > 0 {
		js.flush()
		js.flushWithExporters() //nolint
		return nil
	}
	go func() {
		t := time.NewTimer(maxWaitExporterTime)
		defer t.Stop()
		<-t.C
		js.cancelFunc()
		if js.logExporters != nil && len(js.logExporters) != 0 {
			for _, exporter := range js.logExporters {
				exporter.Shutdown(js.ctx) //nolint
			}
		}
	}()
	return nil
}

func (js *JsonEncoder) sync() {}

func (js *JsonEncoder) dealIoWriter(f field.Field) {
	//断言为数组形式
	fieldArr, ok := f.(*field.ArrayField)
	if ok {
		for i := 0; i < fieldArr.Length(); i++ {
			oneField, errArr := fieldArr.At(i)
			if errArr != nil {
				log.Println(field.GenerateSpecificError(errArr))
			}
			//将单个log转为byte数组
			err := js.write(oneField)
			if err != nil {
				log.Println(field.GenerateSpecificError(err))
			}
			// _, res := js.WriteBytes(js.End)
			_, writeBytesErr := js.WriteBytes(js.End)
			if writeBytesErr != nil {
				log.Println(field.GenerateSpecificError(writeBytesErr))
			}
			//输出
			js.flush()
		}
	}
}

// dealStdoutExporter 处理控制台输出的特殊逻辑，控制台输出要求单条日志按时间顺序输出，并且只展示关键信息。
func (js *JsonEncoder) dealStdoutExporter(stdoutExporter exporter.LogExporter, logs field.Field) {
	//断言为数组形式
	fieldArr, fieldArrOk := logs.(*field.ArrayField)
	if fieldArrOk {
		for i := 0; i < fieldArr.Length(); i++ {
			// 取出单个Log。
			oneField, errArr := fieldArr.At(i)
			if errArr != nil {
				log.Println(field.GenerateSpecificError(errArr))
			}
			//将单个log转为byte数组
			err := js.write(oneField)
			if err != nil {
				log.Println(field.GenerateSpecificError(err))
			}
			// _, res := js.WriteBytes(js.End)
			_, writeBytesErr := js.WriteBytes(js.End)
			if writeBytesErr != nil {
				log.Println(field.GenerateSpecificError(writeBytesErr))
			}
			//调用stdoutExporter的输出方法将单个log输出
			exportLogsErr := stdoutExporter.ExportLogs(js.ctx, js.bufReal.Bytes())
			if exportLogsErr != nil {
				log.Println(field.GenerateSpecificError(exportLogsErr))
			}
			js.bufReal.Reset()
		}
	}
}

func (js *JsonEncoder) flush() error {
	if js.w != nil {
		_, res := js.w.Write(js.bufReal.Bytes())
		if res != nil {
			log.Println(field.GenerateSpecificError(res))
		}
		js.bufReal.Reset()
	}
	return nil
}

func (js *JsonEncoder) flushWithExporters() error {
	var returnErr error = nil
	if js.logExporters != nil && len(js.logExporters) != 0 {
		for _, e := range js.logExporters {
			//过滤掉已经输出的StdoutExporter，其他exporter正常输出
			if e.Name() == "StdoutExporter" {
				continue
			}
			if err := e.ExportLogs(js.ctx, js.bufReal.Bytes()); err != nil {
				// 如果错误则记日志。
				log.Println(field.GenerateSpecificError(err))
				returnErr = err
			}
		}
		js.bufReal.Reset()
	}
	return returnErr
}

func (js *JsonEncoder) write(f field.Field) error {
	w := js
	switch f.Type() {
	default:
		return nil
	case field.IntType:
		v := strconv.Itoa(int(f.(field.IntField)))
		bytes := js.string2Bytes(v)
		_, res := w.WriteBytes(bytes)
		return res
	case field.Float64Type:
		v := strconv.FormatFloat(float64(f.(field.Float64Field)), 'f', -1, 64)
		bytes := js.string2Bytes(v)
		_, res := w.WriteBytes(bytes)
		return res
	case field.StringType:
		v := string(f.(field.StringField))
		w.WriteBytes(_quotation) //nolint
		js.safeWriteString(v)    //nolint
		_, res := w.WriteBytes(_quotation)
		return res
	case field.TimeType:
		v := strconv.FormatInt(int64(time.Time(f.(field.TimeField)).UnixNano()), 10)
		bytes := js.string2Bytes(v)
		_, res := w.WriteBytes(bytes)
		return res
	case field.ArrayType:
		v := f.(*field.ArrayField)
		w.WriteBytes(_leftBracket) //nolint
		i := 0
		for ; i < len(*v)-1; i += 1 {
			js.write((*v)[i])        //nolint
			w.WriteBytes(_seperator) //nolint
		}
		if i < len(*v) {
			js.write((*v)[i]) //nolint
		}
		_, res := w.WriteBytes(_rightBracket) //nolint
		return res
	case field.StructType:
		fs := f.(*field.StructField)
		w.WriteBytes(_leftBigBracket) //nolint
		i := 0
		for ; i < fs.Length()-1; i += 1 {
			k, v, _ := fs.At(i)
			w.WriteBytes(_quotation) //nolint
			js.safeWriteString(k)    //nolint
			w.WriteBytes(_quotation) //nolint
			w.WriteBytes(_colon)     //nolint
			js.write(v)              //nolint
			w.WriteBytes(_seperator) //nolint
		}
		if k, v, err := fs.At(i); err == nil {
			w.WriteBytes(_quotation) //nolint
			js.safeWriteString(k)    //nolint
			w.WriteBytes(_quotation) //nolint
			w.WriteBytes(_colon)     //nolint
			js.write(v)              //nolint
		}
		_, res := w.WriteBytes(_rightBigBracket)
		return res

	case field.JsonType:
		j := f.(*field.JsonFiled)
		b, err := json.Marshal(j.Data)

		if err != nil {
			return err
		}
		_, res := w.WriteBytes(b)
		return res

	case field.MapType:
		b, err := json.Marshal(f)

		if err != nil {
			return err
		}
		_, res := w.WriteBytes(b)
		return res
	}

}

// String2Bytes unsafe convert string to []byte, they point to the same memory
func (js *JsonEncoder) string2Bytes(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	result := make([]byte, sh.Len)
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&result))
	bh.Data = sh.Data
	return result
	// bh := reflect.SliceHeader{
	// 	Data: sh.Data,
	// 	Len:  sh.Len,
	// 	Cap:  sh.Len,
	// }
	// return *(*[]byte)(unsafe.Pointer(&bh))
}

func (js *JsonEncoder) safeWriteString(s string) (int, error) {
	left := 0
	w := js
	bytes := js.string2Bytes(s)
	var res error
	for k, i := range s {
		switch i {
		default:
		case '"':
			w.WriteBytes(bytes[left:k]) //nolint
			left = k + 1
			_, res = w.WriteBytes(_quotationSafe) //nolint
		case '\\':
			w.WriteBytes(bytes[left:k]) //nolint
			left = k + 1
			_, res = w.WriteBytes(_reverseSafe) //nolint
		case '\b':
			w.WriteBytes(bytes[left:k]) //nolint
			left = k + 1
			_, res = w.WriteBytes(_backspaceSafe) //nolint
		case '\f':
			w.WriteBytes(bytes[left:k]) //nolint
			left = k + 1
			_, res = w.WriteBytes(_formfeedSafe) //nolint
		case '\t':
			w.WriteBytes(bytes[left:k]) //nolint
			left = k + 1
			_, res = w.WriteBytes(_horizontalSafe) //nolint
		case '\n':
			w.WriteBytes(bytes[left:k]) //nolint
			left = k + 1
			_, res = w.WriteBytes(_lineFeedSafe) //nolint
		case '\r':
			w.WriteBytes(bytes[left:k]) //nolint
			left = k + 1
			_, res = w.WriteBytes(_carriageSafe) //nolint
		}
		if res != nil {
			return left, res
		}
	}

	if left < len(bytes) {
		_, res = w.WriteBytes(bytes[left:]) //nolint
	}
	return left, res
}

func (js *JsonEncoder) WriteBytes(c []byte) (int, error) {
	return js.buf.Write(c)
}
