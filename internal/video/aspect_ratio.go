package video

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"os/exec"
)

func GetVideoAspectRatio(filePath string) (string, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-print_format", "json", "-show_streams", filePath)
	var out bytes.Buffer
	var e bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &e
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("ffprobe error: %v - std err: %v", err, e.String())
	}

	data := &struct {
		Streams []struct {
			Height int `json:"height"`
			Width  int `json:"width"`
		} `json:"streams"`
	}{}

	decoder := json.NewDecoder(&out)
	if err := decoder.Decode(data); err != nil {
		return "", err
	}
	ratio := float64(data.Streams[0].Width) / float64(data.Streams[0].Height)
	epsilon := 0.1
	aspectRatio := "other"
	if math.Abs((16.0/9.0)-ratio) < epsilon {
		aspectRatio = "16:9"
	} else if math.Abs((9.0/16.0)-ratio) < epsilon {
		aspectRatio = "9:16"
	}
	return aspectRatio, nil
}
