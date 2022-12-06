package eventsdk

import (
	"github.com/shirou/gopsutil/v3/host"
	"reflect"
	"testing"
)

func TestNewAttribute(t *testing.T) {
	type args struct {
		key string
		v   Value
	}
	tests := []struct {
		name string
		args args
		want Attribute
	}{
		{
			"",
			args{
				key: "key",
				v:   StringValue("123"),
			},
			&attribute{
				Key: "key",
				Value: value{
					Type: "STRING",
					Data: "123",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAttribute(tt.args.key, tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAttribute() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAttributeGetKey(t *testing.T) {
	type fields struct {
		Key   string
		Value value
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"",
			fields{
				Key: "key",
				Value: value{
					Type: "STRING",
					Data: "123",
				},
			},
			"key",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &attribute{
				Key:   tt.fields.Key,
				Value: tt.fields.Value,
			}
			if got := a.GetKey(); got != tt.want {
				t.Errorf("GetKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAttributeGetValue(t *testing.T) {
	type fields struct {
		Key   string
		Value value
	}
	tests := []struct {
		name   string
		fields fields
		want   Value
	}{
		{
			"",
			fields{
				Key: "key",
				Value: value{
					Type: "STRING",
					Data: "123",
				},
			},
			value{
				Type: "STRING",
				Data: "123",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &attribute{
				Key:   tt.fields.Key,
				Value: tt.fields.Value,
			}
			if got := a.GetValue(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAttributeValid(t *testing.T) {
	type fields struct {
		Key   string
		Value value
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			"",
			fields{
				Key: "key",
				Value: value{
					Type: "STRING",
					Data: "123",
				},
			},
			true,
		}, {
			"",
			fields{
				Key: "",
				Value: value{
					Type: "STRING",
					Data: "123",
				},
			},
			false,
		}, {
			"",
			fields{
				Key: "os",
				Value: value{
					Type: "STRING",
					Data: "123",
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &attribute{
				Key:   tt.fields.Key,
				Value: tt.fields.Value,
			}
			if got := a.Valid(); got != tt.want {
				t.Errorf("Valid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAttributeKeyNotCollide(t *testing.T) {
	type fields struct {
		Key   string
		Value value
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			"",
			fields{
				Key: "key",
				Value: value{
					Type: "STRING",
					Data: "123",
				},
			},
			true,
		}, {
			"",
			fields{
				Key: "host",
				Value: value{
					Type: "STRING",
					Data: "123",
				},
			},
			false,
		}, {
			"",
			fields{
				Key: "os",
				Value: value{
					Type: "STRING",
					Data: "123",
				},
			},
			false,
		}, {
			"",
			fields{
				Key: "telemetry",
				Value: value{
					Type: "STRING",
					Data: "123",
				},
			},
			false,
		}, {
			"",
			fields{
				Key: "service",
				Value: value{
					Type: "STRING",
					Data: "123",
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &attribute{
				Key:   tt.fields.Key,
				Value: tt.fields.Value,
			}
			if got := a.keyNotCollide(); got != tt.want {
				t.Errorf("keyNotCollide() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAttributeKeyNotNil(t *testing.T) {
	type fields struct {
		Key   string
		Value value
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			"",
			fields{
				Key: "key",
				Value: value{
					Type: "STRING",
					Data: "123",
				},
			},
			true,
		}, {
			"",
			fields{
				Key: "",
				Value: value{
					Type: "STRING",
					Data: "123",
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &attribute{
				Key:   tt.fields.Key,
				Value: tt.fields.Value,
			}
			if got := a.keyNotNil(); got != tt.want {
				t.Errorf("keyNotNil() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAttributePrivate(t *testing.T) {
	type fields struct {
		Key   string
		Value value
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			"",
			fields{
				Key: "key",
				Value: value{
					Type: "STRING",
					Data: "123",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &attribute{
				Key:   tt.fields.Key,
				Value: tt.fields.Value,
			}
			a.private()
		})
	}
}

func TestGetDefaultAttributes(t *testing.T) {
	tests := []struct {
		name string
		want map[string]interface{}
	}{
		{
			"",
			getDefaultAttributes(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getDefaultAttributes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getDefaultAttributes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetHostIP(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			"",
			getHostIP(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getHostIP(); got != tt.want {
				t.Errorf("getHostIP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetHostInfo(t *testing.T) {
	tests := []struct {
		name string
		want *host.InfoStat
	}{
		{
			"",
			getHostInfo(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getHostInfo(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getHostInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}
