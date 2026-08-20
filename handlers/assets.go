package handlers

import (
	"crypto/rand"
	"encoding/base64"
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

func getAssetName(mediaType string) string {
	key := make([]byte, 32)
	rand.Read(key)
	name := base64.RawURLEncoding.EncodeToString(key)

	ext := mediaTypeToExt(mediaType)
	return fmt.Sprintf("%s%s", name, ext)
}

func (cfg ApiConfig) getAssertDiskPath(assetName string) string {
	return filepath.Join(cfg.AssetsRoot, assetName)
}

func (cfg ApiConfig) getObjectUrl(key string) string {
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", cfg.S3Bucket, cfg.S3Region, key)
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
