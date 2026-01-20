package utils

import (
	"os"
	"path/filepath"
	"strings"
)

const (
	MaxResumeSizeBytes = 5 * 1024 * 1024
)

func SaveContentToFile(folder, filename, ext string, content []byte) (string, error) {
	if err := os.MkdirAll(folder, 0755); err != nil {
		return "", err
	}

	base := filepath.Base(filename)
	base = strings.TrimSuffix(base, filepath.Ext(base))

	if strings.HasPrefix(ext, ".") {
		ext = ext[1:]
	}

	fullPath := filepath.Join(folder, base+"."+ext)

	if err := os.WriteFile(fullPath, content, 0644); err != nil {
		return "", err
	}

	return fullPath, nil
}

func IsLargeFile(size int64) bool {
	return size > MaxResumeSizeBytes
}
