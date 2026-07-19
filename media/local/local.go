package local

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/dankobg/fluffly/media"
)

const (
	StorageKind = "local"
)

type LocalUploader struct {
	uploadDir string
	baseURL   string
}

func NewLocalUploader(baseURL, uploadDir string) (*LocalUploader, error) {
	if baseURL == "" {
		return nil, errors.New("base URL not provided")
	}

	if uploadDir == "" {
		return nil, errors.New("upload dir not provided")
	}

	absUploadDir, err := filepath.Abs(uploadDir)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve upload dir: %w", err)
	}

	return &LocalUploader{
		baseURL:   baseURL,
		uploadDir: absUploadDir,
	}, nil
}

func (lu *LocalUploader) Kind() string {
	return StorageKind
}

func (lu *LocalUploader) Upload(ctx context.Context, filename string, r io.Reader, size int64) (string, error) {
	if rc, ok := r.(io.ReadCloser); ok {
		defer func() { _ = rc.Close() }()
	}

	objectName := filepath.Clean(strings.TrimSpace(filename))
	destFileName := filepath.Join(lu.uploadDir, objectName)

	destDir := filepath.Dir(destFileName)
	if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create upload dir: %w", err)
	}

	destFile, err := os.Create(destFileName)
	if err != nil {
		return "", fmt.Errorf("failed to create dest file - %s: %w", filename, err)
	}

	defer func() { _ = destFile.Close() }()

	if _, err := io.Copy(destFile, r); err != nil {
		return "", fmt.Errorf("failed to copy file contents: %w", err)
	}

	return objectName, nil
}

func (lu *LocalUploader) Delete(ctx context.Context, filename string) error {
	objectName := filepath.Clean(strings.TrimSpace(filename))

	filePath := filepath.Join(lu.uploadDir, objectName)
	if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete file %s: %w", filePath, err)
	}

	return nil
}

func (lu *LocalUploader) URL(name string, kind string) (string, error) {
	if kind == media.StorageKindExternal {
		return name, nil
	}

	url, err := url.JoinPath(lu.baseURL, name)
	if err != nil {
		return "", fmt.Errorf("failed to resolve file url - %s: %w", name, err)
	}

	return url, nil
}
