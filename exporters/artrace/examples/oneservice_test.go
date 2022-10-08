package examples

import (
	"context"
	"testing"
)

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

func TestHTTPSExample(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			"HTTPSExample",
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
			"TestWithAllExample",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			WithAllExample()
		})
	}
}

func Test_Add(t *testing.T) {
	type args struct {
		ctx context.Context
		x   int64
		y   int64
	}
	tests := []struct {
		name  string
		args  args
		want  context.Context
		want1 int64
	}{
		{
			name: "add",
			args: args{
				context.Background(),
				3,
				4,
			},
			want:  context.Background(),
			want1: 7,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, got1 := add(tt.args.ctx, tt.args.x, tt.args.y)
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("add() got = %v, want %v", got, tt.want)
			//}
			if got1 != tt.want1 {
				t.Errorf("add() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_Multiply(t *testing.T) {
	type args struct {
		ctx context.Context
		x   int64
		y   int64
	}
	tests := []struct {
		name  string
		args  args
		want  context.Context
		want1 int64
	}{
		{
			name: "multiply",
			args: args{
				context.Background(),
				8,
				2,
			},
			want:  context.Background(),
			want1: 16,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, got1 := multiply(tt.args.ctx, tt.args.x, tt.args.y)
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("multiply() got = %v, want %v", got, tt.want)
			//}
			if got1 != tt.want1 {
				t.Errorf("multiply() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
