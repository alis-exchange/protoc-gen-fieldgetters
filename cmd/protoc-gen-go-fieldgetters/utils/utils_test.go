package utils

import "testing"

func TestToLowerFirst(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test All Uppercase",
			args: args{s: "HELLO"},
			want: "hELLO",
		},
		{
			name: "Test All Lowercase",
			args: args{s: "hello"},
			want: "hello",
		},
		{
			name: "Test Empty String",
			args: args{s: ""},
			want: "",
		},
		{
			name: "Test Single Letter",
			args: args{s: "H"},
			want: "h",
		},
		{
			name: "Test Single Letter Lowercase",
			args: args{s: "h"},
			want: "h",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToLowerFirst(tt.args.s); got != tt.want {
				t.Errorf("ToLowerFirst() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSnakeCaseToCamelCase(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test Snake Case",
			args: args{s: "hello_world"},
			want: "helloWorld",
		},
		{
			name: "Test Camel Case",
			args: args{s: "helloWorld"},
			want: "helloWorld",
		},
		{
			name: "Test Empty String",
			args: args{s: ""},
			want: "",
		},
		{
			name: "Test Single Word",
			args: args{s: "Hello"},
			want: "hello",
		},
		{
			name: "Test Snake Case All Uppercase",
			args: args{s: "HELLO_WORLD"},
			want: "helloWorld",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SnakeCaseToCamelCase(tt.args.s); got != tt.want {
				t.Errorf("SnakeCaseToCamelCase() = %v, want %v", got, tt.want)
			}
		})
	}
}
