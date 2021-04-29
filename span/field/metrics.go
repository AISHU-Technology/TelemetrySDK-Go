package field

import "time"

type Mmetric struct {
	Attrs  StructField
	Labels ArrayField
	Name   StringField
	Value  Float64Field
	Time   TimeField
}

func (m *Mmetric) Set(name string, value float64) {
	m.Name = StringField(name)
	m.Value = Float64Field(value)
	m.Time = TimeField(time.Now())
}

func (m *Mmetric) AddAttribute(key, value string) {
	m.Attrs.Set(key, StringField(value))
}

func (m *Mmetric) AddLabel(l string) {
	m.Labels = append(m.Labels, StringField(l))
}
