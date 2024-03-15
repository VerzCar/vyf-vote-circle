package utils

import "testing"

func TestIsImageMimeType(t *testing.T) {
	type args struct {
		mimeType string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "should return true for correct mime type",
			args: struct{ mimeType string }{mimeType: "image/png"},
			want: true,
		},
		{
			name: "should return false for incorrect mime type",
			args: struct{ mimeType string }{mimeType: "text/plain"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := IsImageMimeType(tt.args.mimeType); got != tt.want {
					t.Errorf("IsImageMimeType() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
