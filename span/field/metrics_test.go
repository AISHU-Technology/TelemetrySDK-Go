package field

import (
	"testing"

	"gotest.tools/assert"
)

func TestMmetricSet(t *testing.T) {
	m := &Mmetric{}
	m.Set("test", 0.0)
	assert.Equal(t, string(m.Name), "test")
	assert.Equal(t, float64(m.Value), 0.0)
}

func TestMmetricAddAttribute(t *testing.T) {
	m := &Mmetric{}
	m.AddAttribute("test0", "test0")
	m.AddAttribute("test1", "test1")
	assert.Equal(t, string(m.Attrs.keys[0]), "test0")
	assert.Equal(t, m.Attrs.values[0].(StringField), StringField("test0"))

	assert.Equal(t, string(m.Attrs.keys[1]), "test1")
	assert.Equal(t, m.Attrs.values[1].(StringField), StringField("test1"))
}

func TestMmetricAddLabel(t *testing.T) {
	m := &Mmetric{}
	m.AddLabel("test0")
	m.AddLabel("test1")
	assert.Equal(t, string(m.Labels[0].(StringField)), "test0")
	assert.Equal(t, string(m.Labels[1].(StringField)), "test1")
}

// func BenchmarkMemetrics(b *testing.B) {
// 	mm := Mmetric{}
// 	mm.Set("test", 0)
// 	key := "test"
// 	for i := 0; i < 10; i++ {
// 		mm.AddAttribute(key, key)
// 	}

// 	b.Run("metrics/json", func(b *testing.B) {
// 		b.ResetTimer()
// 		b.RunParallel(func(p *testing.PB) {
// 			for p.Next() {
// 				json.Marshal(mm)
// 			}
// 		})
// 	})

// 	buf := bytes.NewBuffer(nil)
// 	w := &writer.JsonWriter{
// 		W: buf,
// 	}
// 	enc := encoder.NewJsonEncoder(w)
// 	b.Run("metrics/self", func(b *testing.B) {
// 		b.ResetTimer()
// 		b.RunParallel(func(p *testing.PB) {
// 			for p.Next() {
// 				enc.Write(mm)
// 				buf.Reset()
// 			}
// 		})
// 	})

// }
