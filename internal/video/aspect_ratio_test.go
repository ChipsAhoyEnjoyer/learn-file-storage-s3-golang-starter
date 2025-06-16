package video

import (
	"testing"
)

func Test_getVideoAspectRatio(t *testing.T) {
	type args struct {
		filePath string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Get Aspect Ratio",
			args: args{
				filePath: "./internal/video/test_video.mp4",
			},
			want:    "16:9",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getVideoAspectRatio(tt.args.filePath)
			if (err != nil) != tt.wantErr {
				t.Log(err)
				t.Errorf("getVideoAspectRatio() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getVideoAspectRatio() = %v, want %v", got, tt.want)
			}
		})
	}
}
