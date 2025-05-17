package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

func (cfg apiConfig) ensureAssetsDir() error {
	if _, err := os.Stat(cfg.assetsRoot); os.IsNotExist(err) {
		return os.Mkdir(cfg.assetsRoot, 0755)
	}
	return nil
}

func getAssetName(fileID uuid.UUID, extension string) string {
	return fmt.Sprintf("%s.%s", fileID.String(), extension)
}

func getAssetExtension(mediaType string) (string, error) {
	i := strings.Split(mediaType, "/")
	if len(i) != 2 {
		return "", fmt.Errorf("incorrect media type;")
	}
	return i[1], nil
}

func getAssetPath(filename, assetsRoot string) string {
	return "." + string(os.PathSeparator) + filepath.Join(assetsRoot, filename)
}

func getAssetURL(filename, port string) string {
	return fmt.Sprintf("http://localhost:%s/assets/%s", port, filename)
}
