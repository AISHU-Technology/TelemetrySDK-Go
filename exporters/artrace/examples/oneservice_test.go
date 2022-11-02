package examples

import (
	"testing"
)

func TestStdoutExample(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			"Stdout",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			StdoutExample()
		})
	}
}

func TestHTTPExample(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			"Test发送到AnyRobot Trace数据接收器",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			HTTPExample()
		})
	}
}

func TestHTTPSExample(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			"Test发送到AnyRobot Trace数据接收器",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			HTTPSExample()
		})
	}
}

func TestWithAllExample(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			"Test调用全部接口的示例",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			WithAllExample()
		})
	}
}

func TestExample(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			"Test原始的业务系统",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Example()
		})
	}
}
