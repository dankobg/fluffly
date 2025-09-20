package media

import (
	"context"
	"io"
)

const StorageKindExternal = "external"

type Uploader interface {
	Kind() string
	Upload(ctx context.Context, filename string, r io.Reader, size int64) (string, error)
	Delete(ctx context.Context, filename string) error
	URL(name string, kind string) (string, error)
}
