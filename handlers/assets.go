package handlers

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
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

func getVideoAspectRatio(filePath string) (string, error) {
	command := fmt.Sprintf("ffprobe -v error -print_format json -show_streams %s", filePath)
	args := strings.Fields(command)
	cmd := exec.Command(args[0], args[1:]...)
	var stdOut, stdErr bytes.Buffer
	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("ffprobe error: %w\nStdErr: %s", err, stdErr.String())
	}
	var results struct {
		Streams []struct {
			Width              int    `json:"width,omitempty"`
			Height             int    `json:"height,omitempty"`
			CodedWidth         int    `json:"coded_width,omitempty"`
			CodedHeight        int    `json:"coded_height,omitempty"`
			SampleAspectRatio  string `json:"sample_aspect_ratio,omitempty"`
			DisplayAspectRatio string `json:"display_aspect_ratio,omitempty"`
		} `json:"streams"`
	}

	err = json.Unmarshal(stdOut.Bytes(), &results)
	if err != nil {
		return "", err
	}

	return results.Streams[0].DisplayAspectRatio, nil
}

func processVideoForFastStart(filePath string) (string, error) {
	outPath := filePath + ".processing"
	command := fmt.Sprintf("ffmpeg -i %s -c copy -movflags faststart -f mp4 %s", filePath, outPath)
	log.Println(command)
	args := strings.Fields(command)
	cmd := exec.Command(args[0], args[1:]...)

	var stdOut, stdErr bytes.Buffer
	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("ffprobe error: %w\nStdErr: %s", err, stdErr.String())
	}

	return outPath, nil
}
