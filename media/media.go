package media

import (
	"bytes"
	"context"
	"io"
	"net/http"
)

const StorageKindExternal = "external"

type Uploader interface {
	Kind() string
	Upload(ctx context.Context, filename string, r io.Reader, size int64) (string, error)
	Delete(ctx context.Context, filename string) error
	URL(name string, kind string) (string, error)
}

func DetectContentType(r io.Reader) (string, io.Reader, error) {
	buf := make([]byte, 512)

	n, err := r.Read(buf)
	if err != nil && err != io.EOF {
		return "", r, err
	}

	contentType := http.DetectContentType(buf[:n])
	r = io.MultiReader(bytes.NewReader(buf[:n]), r)

	return contentType, r, nil
}
