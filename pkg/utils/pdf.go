package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/anshu4sharma/resume_ats/internal/config"
	"github.com/anshu4sharma/resume_ats/internal/structs"
)

func ExtractTextFromPdfBox(src io.Reader, filename string) (*structs.JavaExtractionResponse, error) {
	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)

	errChan := make(chan error, 1)

	go func() {
		defer pw.Close()
		defer writer.Close()
		part, err := writer.CreateFormFile("file", filename)
		if err != nil {
			errChan <- err
			return
		}
		_, err = io.Copy(part, src)
		errChan <- err
	}()

	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
		},
	}
	resp, err := client.Post(config.Load().PDFBOX_URL, writer.FormDataContentType(), pr)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	select {
	case err := <-errChan:
		if err != nil {
			return nil, fmt.Errorf("streaming error: %w", err)
		}
	default:
	}

	var result structs.JavaExtractionResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}
