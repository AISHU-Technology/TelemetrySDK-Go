package examples

import "testing"

func TestStdoutExporterExample(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			"StdoutExporterExample",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			StdoutExporterExample()
		})
	}
}
func TestHTTPExample(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			"HTTPExample",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			HTTPExample()
		})
	}
}

func TestDefaultExporterExample(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			"DefaultExporterExample",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			DefaultExporterExample()
		})
	}
}

func TestStdoutExample(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			"StdoutExample",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			StdoutExample()
		})
	}
}

func TestExample(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Example()
		})
	}
}
