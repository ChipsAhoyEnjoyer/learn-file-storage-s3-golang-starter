package video

import (
	"bytes"
	"fmt"
	"os/exec"
)

func ProcessVideoForFastStart(filePath string) (string, error) {
	out := filePath + ".processing"
	cmd := exec.Command(
		"ffmpeg",
		"-i", filePath,
		"-c", "copy",
		"-movflags", "faststart",
		"-f", "mp4",
		out,
	)
	var e bytes.Buffer
	cmd.Stderr = &e
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("ffmpeg error: %v - std err: %s", err, e.String())
	}
	return out, nil
}
