package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"math/rand"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

const (
	MaxResumeSizeBytes = 10 * 1024 * 1024
	MaxResumeSize      = 3 * 1024 * 1024
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateID(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
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

func HashBytes(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
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
