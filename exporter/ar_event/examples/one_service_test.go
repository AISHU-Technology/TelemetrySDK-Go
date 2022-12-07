package examples

import "testing"

func TestStdoutExample(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			StdoutExample()
		})
	}
}

func TestWithAllExample(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			HTTPExample()
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
