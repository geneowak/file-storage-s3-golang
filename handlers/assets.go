package handlers

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func (cfg ApiConfig) ensureAssetsDir() error {
	if _, err := os.Stat(cfg.AssetsRoot); os.IsNotExist(err) {
		return os.Mkdir(cfg.AssetsRoot, 0755)
	}
	return nil
}

func getAssetName(fileName string, mediaType string) string {
	ext := mediaTypeToExt(mediaType)
	return fmt.Sprintf("%s%s", fileName, ext)
}

func (cfg ApiConfig) getAssertDiskPath(assetName string) string {
	return filepath.Join(cfg.AssetsRoot, assetName)
}

func (cfg ApiConfig) getAssertURL(assertName string) string {
	return fmt.Sprintf("http://localhost:%s/%s", cfg.Port, assertName)
}

func mediaTypeToExt(mediaType string) string {
	parts := strings.Split(mediaType, "/")
	if len(parts) != 2 {
		return ".bin"
	}
	return "." + parts[1]
}
