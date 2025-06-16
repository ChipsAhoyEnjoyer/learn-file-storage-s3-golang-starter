package video

import (
	"bytes"
	"encoding/json"
	"log"
	"math"
	"os/exec"
)

func getVideoAspectRatio(filePath string) (string, error) {
	command := "ffprobe"
	cmd := exec.Command(command, "-v", "error", "-print_format", "json", "-show_streams", filePath)
	stdout := new(bytes.Buffer)
	cmd.Stdout = stdout

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	data := &struct {
		Height int `json:"height"`
		Width  int `json:"width"`
	}{}
	decoder := json.NewDecoder(stdout)
	if err = decoder.Decode(data); err != nil {
		return "", err
	}
	log.Println(stdout)
	ratio := float64(data.Width) / float64(data.Height)
	epsilon := 0.1
	aspectRatio := "other"
	if math.Abs((16.0/9.0)-ratio) < epsilon {
		aspectRatio = "16:9"
	} else if math.Abs((9.0/16.0)-ratio) < epsilon {
		aspectRatio = "9:16"
	}
	return aspectRatio, nil
}
