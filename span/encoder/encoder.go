package encoder

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"reflect"
	"strconv"
	"time"
	"unsafe"

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

type JsonEncoder struct {
	w       io.Writer
	buf     io.Writer
	End     []byte
	bufReal *bytes.Buffer
	Ctx     context.Context
	Cancel  context.CancelFunc
}

func NewJsonEncoder(w io.Writer) Encoder {
	ctx, cancel := context.WithCancel(context.Background())
	res := &JsonEncoder{
		w:       w,
		End:     _lineFeed,
		bufReal: bytes.NewBuffer(make([]byte, 0, 4096)),
		Ctx:     ctx,
		Cancel:  cancel,
	}
	res.buf = res.bufReal
	return res
}

func NewJsonEncoderBench(w io.Writer) Encoder {
	ctx, cancel := context.WithCancel(context.Background())
	res := &JsonEncoder{
		w:       w,
		End:     _lineFeed,
		bufReal: bytes.NewBuffer(make([]byte, 0, 4096)),
		buf:     ioutil.Discard,
		Ctx:     ctx,
		Cancel:  cancel,
	}
	return res
}

// buffer by encoder for output
func (js *JsonEncoder) Write(f field.Field) error {
	js.write(f)
	// _, res := js.WriteBytes(js.End)
	js.WriteBytes(js.End)
	return js.flush()

}

func (js *JsonEncoder) flush() error {
	_, res := js.w.Write(js.bufReal.Bytes())
	if res != nil {
		panic(res)
	}
	js.bufReal.Reset()
	return res
}

func (js *JsonEncoder) getContext() context.Context {
	return js.Ctx
}

func (js *JsonEncoder) Close() error {
	if js.bufReal.Len() > 0 {
		return js.flush()
	}
	go func() {
		t := time.NewTimer(maxWaitExporterTime)
		defer t.Stop()
		select {
		case <-t.C:
			js.Cancel()
		}
	}()
	return nil
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
		w.WriteBytes(_quotation)
		js.safeWriteString(v)
		_, res := w.WriteBytes(_quotation)
		return res
	case field.TimeType:
		v := strconv.FormatInt(int64(time.Time(f.(field.TimeField)).UnixNano()), 10)
		bytes := js.string2Bytes(v)
		_, res := w.WriteBytes(bytes)
		return res
	case field.ArrayType:
		v := f.(*field.ArrayField)
		w.WriteBytes(_leftBracket)
		i := 0
		for ; i < len(*v)-1; i += 1 {
			js.write((*v)[i])
			w.WriteBytes(_seperator)
		}
		if i < len(*v) {
			js.write((*v)[i])
		}
		_, res := w.WriteBytes(_rightBracket)
		return res
	case field.StructType:
		fs := f.(*field.StructField)
		w.WriteBytes(_leftBigBracket)
		i := 0
		for ; i < fs.Length()-1; i += 1 {
			k, v, _ := fs.At(i)
			w.WriteBytes(_quotation)
			js.safeWriteString(k)
			w.WriteBytes(_quotation)
			w.WriteBytes(_colon)
			js.write(v)
			w.WriteBytes(_seperator)
		}
		if k, v, err := fs.At(i); err == nil {
			w.WriteBytes(_quotation)
			js.safeWriteString(k)
			w.WriteBytes(_quotation)
			w.WriteBytes(_colon)
			js.write(v)
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
	result := make([]byte, sh.Len, sh.Len)
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
			w.WriteBytes(bytes[left:k])
			left = k + 1
			_, res = w.WriteBytes(_quotationSafe)
		case '\\':
			w.WriteBytes(bytes[left:k])
			left = k + 1
			_, res = w.WriteBytes(_reverseSafe)
		case '\b':
			w.WriteBytes(bytes[left:k])
			left = k + 1
			_, res = w.WriteBytes(_backspaceSafe)
		case '\f':
			w.WriteBytes(bytes[left:k])
			left = k + 1
			_, res = w.WriteBytes(_formfeedSafe)
		case '\t':
			w.WriteBytes(bytes[left:k])
			left = k + 1
			_, res = w.WriteBytes(_horizontalSafe)
		case '\n':
			w.WriteBytes(bytes[left:k])
			left = k + 1
			_, res = w.WriteBytes(_lineFeedSafe)
		case '\r':
			w.WriteBytes(bytes[left:k])
			left = k + 1
			_, res = w.WriteBytes(_carriageSafe)
		}
		if res != nil {
			return left, res
		}
	}

	if left < len(bytes) {
		_, res = w.WriteBytes(bytes[left:])
	}
	return left, res
}

func (js *JsonEncoder) WriteBytes(c []byte) (int, error) {
	return js.buf.Write(c)
}
