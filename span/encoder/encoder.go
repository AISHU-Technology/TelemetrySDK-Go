package encoder

import (
	"bytes"
	"io"
	"io/ioutil"
	"reflect"
	"span/field"
	"strconv"
	"time"
	"unsafe"
)

var (
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
}

func NewJsonEncoder(w io.Writer) Encoder {
	res := &JsonEncoder{
		w:       w,
		End:     _lineFeed,
		bufReal: bytes.NewBuffer(make([]byte, 0, 4096)),
	}
	res.buf = res.bufReal
	return res
}

func NewJsonEncoderBench(w io.Writer) Encoder {
	res := &JsonEncoder{
		w:       w,
		End:     _lineFeed,
		bufReal: bytes.NewBuffer(make([]byte, 0, 4096)),
		buf:     ioutil.Discard,
	}
	return res
}

// buffer by encoder for output
func (js *JsonEncoder) Write(f field.Field) error {
	js.write(f)
	// _, res := js.WriteBytes(js.End)
	js.WriteBytes(js.End)

	if js.bufReal.Len() >= 4096 {
		return js.flush()
	}
	return nil
}

func (js *JsonEncoder) flush() error {
	_, res := js.w.Write(js.bufReal.Bytes())
    if res != nil {
        panic(res)
    }
	js.bufReal.Reset()
	return res
}

func (js *JsonEncoder) Close() error {
	if js.bufReal.Len() > 0 {
		return js.flush()
	}
	return nil
}

// type JsonEncoder struct {
// 	buf io.Writer
// 	End []byte
// }

// func NewJsonEncoder(w io.Writer) Encoder {
// 	return &JsonEncoder{
// 		buf: w,
// 		End: _lineFeed,
// 	}
// }

// func (js *JsonEncoder) Write(f field.Field) error {
// 	js.write(f)
// 	_, res := js.WriteBytes(js.End)
// 	return res
// }

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
		v := strconv.FormatInt(time.Time(f.(field.TimeField)).Unix(), 10)
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
	case field.MetricType:
		m := f.(*field.Mmetric)
		w.WriteBytes(_leftBigBracket)
		js.write(m.Name)
		w.WriteBytes(_colon)
		js.write(m.Value)
		w.WriteBytes(_seperator)

		w.WriteBytes(_quotation)
		js.safeWriteString("Attributes")
		w.WriteBytes(_quotation)
		w.WriteBytes(_colon)
		js.write(&m.Attrs)
		w.WriteBytes(_seperator)

		w.WriteBytes(_quotation)
		js.safeWriteString("Labels")
		w.WriteBytes(_quotation)
		w.WriteBytes(_colon)
		js.write(&m.Labels)

		_, res := w.WriteBytes(_rightBigBracket)
		return res
	case field.ExternalSpanType:
		span := f.(*field.ExternalSpanField)
		w.WriteBytes(_leftBigBracket)

		w.WriteBytes(_quotation)
		js.safeWriteString("TraceId")
		w.WriteBytes(_quotation)
		w.WriteBytes(_colon)
		w.WriteBytes(_quotation)
		js.WriteString(span.TraceID())
		w.WriteBytes(_quotation)
		w.WriteBytes(_seperator)

		w.WriteBytes(_quotation)
		js.safeWriteString("ParentId")
		w.WriteBytes(_quotation)
		w.WriteBytes(_colon)
		w.WriteBytes(_quotation)
		js.WriteString(span.ParentID())
		w.WriteBytes(_quotation)
		w.WriteBytes(_seperator)

		w.WriteBytes(_quotation)
		js.safeWriteString("InternalParentId")
		w.WriteBytes(_quotation)
		w.WriteBytes(_colon)
		w.WriteBytes(_quotation)
		js.WriteString(span.InternalParentID())
		w.WriteBytes(_quotation)
		w.WriteBytes(_seperator)

		w.WriteBytes(_quotation)
		js.safeWriteString("SpanId")
		w.WriteBytes(_quotation)
		w.WriteBytes(_colon)
		w.WriteBytes(_quotation)
		js.WriteString(span.ID())
		w.WriteBytes(_quotation)
		w.WriteBytes(_seperator)

		w.WriteBytes(_quotation)
		js.safeWriteString("StartTime")
		w.WriteBytes(_quotation)
		w.WriteBytes(_colon)
		js.write(field.TimeField(span.StartTime))
		w.WriteBytes(_seperator)

		w.WriteBytes(_quotation)
		js.safeWriteString("EndTime")
		w.WriteBytes(_quotation)
		w.WriteBytes(_colon)
		js.write(field.TimeField(span.EndTime))
		w.WriteBytes(_seperator)

		w.WriteBytes(_quotation)
		js.safeWriteString("Attributes")
		w.WriteBytes(_quotation)
		w.WriteBytes(_colon)
		js.write(&span.Attributes)

		_, res := w.WriteBytes(_rightBigBracket)
		return res

	}
}

// String2Bytes unsafe convert string to []byte, they point to the same memory
func (js *JsonEncoder) string2Bytes(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&bh))
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

func (js *JsonEncoder) WriteString(c string) (int, error) {
	b := js.string2Bytes(c)
	n, err := js.buf.Write(b)
	return n, err
}

func (js *JsonEncoder) WriteBytes(c []byte) (int, error) {
	return js.buf.Write(c)
}
