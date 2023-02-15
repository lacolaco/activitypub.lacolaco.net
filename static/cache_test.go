package static

import "testing"

func TestDetectCacheControl(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "index.html",
			args: args{
				filename: "index.html",
			},
			want: "no-cache",
		},
		{
			name: "production build css",
			args: args{
				filename: "main.0123456789abcdef.css",
			},
			want: "max-age=31536000, public",
		},
		{
			name: "production build js",
			args: args{
				filename: "main.0123456789abcdef.js",
			},
			want: "max-age=31536000, public",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DetectCacheControl(tt.args.filename); got != tt.want {
				t.Errorf("DetectCacheControl() = %v, want %v", got, tt.want)
			}
		})
	}
}
