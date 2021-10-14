package benchmarks

import (
	"bytes"
	"gitlab.aishu.cn/anyrobot/observability/telemetrysdk/telemetry-go/span/encoder"
	"gitlab.aishu.cn/anyrobot/observability/telemetrysdk/telemetry-go/span/field"
	"sync"
	"testing"
	"time"
)

// var (
// 	_spanMessages = fakeSpanStringFields(1000)

// )

// func fakeSpanStringFields(n int) []span.StringField {
// 	res := make([]span.StringField, n)

// 	for i := 0; i < n; i++ {
// 		res[i].Init(fmt.Sprintf("Test logging, but use a somewhat realistic message length. (#%v)", i))
// 	}
// 	return res
// }

var (
	_oneSapnUser  = &field.StructField{}
	_tenSpanUsers = &field.ArrayField{}
)
var bufferpool = sync.Pool{
	New: func() interface{} {
		return bytes.NewBuffer(nil)
	},
}

func init() {
	_oneSapnUser.Set("Name", field.StringField("Jane Doe"))
	_oneSapnUser.Set("Email", field.StringField("jane@test.com"))
	_oneSapnUser.Set("CreateAt", field.TimeField(time.Date(1980, 1, 1, 12, 0, 0, 0, time.UTC)))

	for i := 0; i < 10; i++ {
		_tenSpanUsers.Append(_oneSapnUser)
	}
}

// func getSpanMessage(i int) field.Field {
// 	return field.StringField(getMessage(i))
// }

func fakeSpanStructField() field.Field {
	res := field.MallocStructField(9)

	res.Set("int", field.IntField(_tenInts[0]))
	ints := field.MallocArrayField(10)
	for _, i := range _tenInts {
		ints.Append(field.IntField(i))
	}
	res.Set("ints", ints)

	res.Set("string", field.StringField(_tenStrings[0]))
	strings := field.MallocArrayField(10)
	for _, i := range _tenInts {
		strings.Append(field.IntField(_tenInts[i]))
	}
	res.Set("strings", strings)

	res.Set("time", field.TimeField(_tenTimes[0]))
	times := field.MallocArrayField(10)
	for _, i := range _tenInts {
		times.Append(field.TimeField(_tenTimes[i]))
	}
	res.Set("times", times)

	res.Set("user1", _oneSapnUser)
	res.Set("user2", _oneSapnUser)
	res.Set("users", _tenSpanUsers)

	return res
}

func fakeStruct() fakeJson {
	t := fakeJson{
		Int:     _tenInts[0],
		Ints:    _tenInts,
		String:  _tenStrings[0],
		Strings: _tenStrings,
		Time:    _tenTimes[0],
		Times:   _tenTimes,
		User1:   _oneUser,
		User2:   _oneUser,
		User:    _tenUsers,
	}
	return t
}

func fakeMap() map[string]interface{} {
	t := map[string]interface{}{
		"int":     _tenInts[0],
		"ints":    _tenInts,
		"string":  _tenStrings[0],
		"strings": _tenStrings,
		"time":    _tenTimes[0],
		"times":   _tenTimes,
		"user1":   _oneUser,
		"user2":   _oneUser,
		"user":    _tenUsers,
	}
	return t
}

type fakeJson struct {
	Int     int
	Ints    []int
	String  string
	Strings []string
	Time    time.Time
	Times   []time.Time
	User1   *user
	User2   *user
	User    []*user
}

func getTestStructField() field.Field {
	res := field.MallocStructField(2)
	res.Set("name", field.StringField("123"))
	res.Set("age", field.IntField(12))
	return res
}

func getTestJson() field.Field {

	type A struct {
		Name string
		Age  int
	}
	var a = &A{Name: "123", Age: 12}

	return field.MallocJsonField(a)
}

func BenchmarkJsonEncoder_Write_Struct(b *testing.B) {
	a := getTestStructField()
	buf := bytes.NewBuffer(nil)
	enc := encoder.NewJsonEncoder(buf)
	enc.Write(a)
}

func BenchmarkJsonEncoder_Write_Json(b *testing.B) {
	a := getTestJson()
	buf := bytes.NewBuffer(nil)
	enc := encoder.NewJsonEncoder(buf)
	enc.Write(a)
}
