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

func TestOldStdoutExample(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			"OldStdoutExample",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			OldStdoutExample()
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
