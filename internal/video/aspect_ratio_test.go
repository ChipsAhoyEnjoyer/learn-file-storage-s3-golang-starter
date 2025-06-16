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
			name: "Get Aspect Ratio 16:9",
			args: args{
				filePath: "test_video_horizontal.mp4",
			},
			want:    "16:9",
			wantErr: false,
		},
		{
			name: "Get Aspect Ratio 9:16",
			args: args{
				filePath: "test_video_vertical.mp4",
			},
			want:    "9:16",
			wantErr: false,
		},
		{
			name: "Get Aspect Ratio No Video",
			args: args{
				filePath: "non_existent_video.mp4",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getVideoAspectRatio(tt.args.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("getVideoAspectRatio() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getVideoAspectRatio() = %v, want %v", got, tt.want)
			}
		})
	}
}
