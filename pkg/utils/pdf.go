package utils

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/ledongthuc/pdf"
)

func ExtractText(path string) (string, int, error) {
	f, r, err := pdf.Open(path)
	if err != nil {
		return "", 0, err
	}
	defer f.Close()

	var sb strings.Builder
	totalPage := r.NumPage()

	for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
		p := r.Page(pageIndex)
		if p.V.IsNull() {
			continue
		}

		content, err := p.GetPlainText(nil)
		if err != nil {
			continue
		}

		sb.WriteString(content)
	}

	return sb.String(), totalPage, nil
}

func IsReadableText(s string) bool {
	if len(s) == 0 {
		return false
	}

	var printable, letters int

	for _, r := range s {
		if unicode.IsPrint(r) {
			printable++
		}
		if unicode.IsLetter(r) || unicode.IsSpace(r) {
			letters++
		}
	}

	printableRatio := float64(printable) / float64(len([]rune(s)))
	letterRatio := float64(letters) / float64(len([]rune(s)))

	return printableRatio > 0.85 && letterRatio > 0.6
}

func ExtractTextOCR(pdfPath string) (string, int, error) {
	tmpDir, err := os.MkdirTemp("", "pdf-ocr-*")
	if err != nil {
		return "", 0, err
	}
	defer os.RemoveAll(tmpDir) // cleanup always

	prefix := filepath.Join(tmpDir, "page")

	cmd := exec.Command(
		"pdftoppm",
		"-png",
		pdfPath,
		prefix,
	)
	if err := cmd.Run(); err != nil {
		return "", 0, fmt.Errorf("pdf to image failed: %w", err)
	}

	images, err := filepath.Glob(prefix + "*.png")
	if err != nil {
		return "", 0, err
	}

	var result bytes.Buffer

	for _, img := range images {
		ocrCmd := exec.Command("tesseract", img, "stdout")
		out, err := ocrCmd.Output()
		if err != nil {
			continue
		}
		result.Write(out)
		result.WriteByte('\n')
	}

	return result.String(), len(images), nil
}
