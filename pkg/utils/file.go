package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

const (
	MaxResumeSizeBytes = 10 * 1024 * 1024
	MaxResumeSize      = 3 * 1024 * 1024
)

func GenerateUUID() string {
	return uuid.NewString()
}

func SaveContentToFile(folder, filename, ext string, content []byte) (string, error) {
	if err := os.MkdirAll(folder, 0755); err != nil {
		return "", err
	}

	base := filepath.Base(filename)
	base = strings.TrimSuffix(base, filepath.Ext(base))

	ext = strings.TrimPrefix(ext, ".")

	fullPath := filepath.Join(folder, base+"."+ext)

	if err := os.WriteFile(fullPath, content, 0644); err != nil {
		return "", err
	}

	return fullPath, nil
}

func IsLargeFile(size int64) bool {
	return size > MaxResumeSizeBytes
}

func HashFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func HashMultipartFile(fileHeader *multipart.FileHeader) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}
