package eventsdk

import (
	"reflect"
	"testing"
)

func TestBoolArray(t *testing.T) {
	type args struct {
		v []bool
	}
	tests := []struct {
		name string
		args args
		want Value
	}{
		{
			"",
			args{
				[]bool{true, true},
			},
			value{
				Type: "BOOLARRAY",
				Data: []bool{true, true},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BoolArray(tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BoolArray() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBoolValue(t *testing.T) {
	type args struct {
		v bool
	}
	tests := []struct {
		name string
		args args
		want Value
	}{
		{
			"",
			args{
				true,
			},
			value{
				Type: "BOOL",
				Data: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BoolValue(tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BoolValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFloatArray(t *testing.T) {
	type args struct {
		v []float64
	}
	tests := []struct {
		name string
		args args
		want Value
	}{
		{
			"",
			args{
				[]float64{1.0, 2.00},
			},
			value{
				Type: "FLOATARRAY",
				Data: []float64{1.0, 2.00},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FloatArray(tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FloatArray() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFloatValue(t *testing.T) {
	type args struct {
		v float64
	}
	tests := []struct {
		name string
		args args
		want Value
	}{
		{
			"",
			args{
				3.3,
			},
			value{
				Type: "FLOAT",
				Data: 3.3,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FloatValue(tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FloatValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntArray(t *testing.T) {
	type args struct {
		v []int
	}
	tests := []struct {
		name string
		args args
		want Value
	}{
		{
			"",
			args{
				[]int{1, 2, 3},
			},
			value{
				Type: "INTARRAY",
				Data: []int{1, 2, 3},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IntArray(tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IntArray() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntValue(t *testing.T) {
	type args struct {
		v int
	}
	tests := []struct {
		name string
		args args
		want Value
	}{
		{
			"",
			args{
				4,
			},
			value{
				Type: "INT",
				Data: 4,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IntValue(tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IntValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringArray(t *testing.T) {
	type args struct {
		v []string
	}
	tests := []struct {
		name string
		args args
		want Value
	}{
		{
			"",
			args{
				[]string{"a", "b", "c"},
			},
			value{
				Type: "STRINGARRAY",
				Data: []string{"a", "b", "c"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringArray(tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StringArray() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringValue(t *testing.T) {
	type args struct {
		v string
	}
	tests := []struct {
		name string
		args args
		want Value
	}{
		{
			"",
			args{
				"string",
			},
			value{
				Type: "STRING",
				Data: "string",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringValue(tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StringValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValueGetData(t *testing.T) {
	type fields struct {
		Type string
		Data interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   interface{}
	}{
		{
			"",
			fields{
				Type: "INT",
				Data: 78,
			},
			78,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := value{
				Type: tt.fields.Type,
				Data: tt.fields.Data,
			}
			if got := v.GetData(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValueGetType(t *testing.T) {
	type fields struct {
		Type string
		Data interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"",
			fields{
				Type: "INT",
				Data: nil,
			},
			"INT",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := value{
				Type: tt.fields.Type,
				Data: tt.fields.Data,
			}
			if got := v.GetType(); got != tt.want {
				t.Errorf("GetType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValuePrivate(t *testing.T) {
	type fields struct {
		Type string
		Data interface{}
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			"",
			fields{
				Type: "",
				Data: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := value{
				Type: tt.fields.Type,
				Data: tt.fields.Data,
			}
			v.private()
		})
	}
}
