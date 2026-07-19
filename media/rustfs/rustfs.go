package rustfs

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/dankobg/fluffly/config"
	"github.com/dankobg/fluffly/media"
	minio "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const StorageKind = "rustfs"

type RustfsUploader struct {
	client     *minio.Client
	bucketName string
	apiURL     string
}

func NewRustfsUploader(cfg config.RustfsConfig) (*RustfsUploader, error) {
	endpoint := net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Address))

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, cfg.Token),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to init rustfs client: %w", err)
	}

	mupl := &RustfsUploader{
		client:     client,
		bucketName: cfg.DefaultBucket,
		apiURL:     client.EndpointURL().String(),
	}
	if err := mupl.ensureBucket(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ensure bucket: %w", err)
	}

	if err := mupl.setPolicy(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to set policy: %w", err)
	}

	return mupl, nil
}

func (mu *RustfsUploader) Kind() string {
	return StorageKind
}

func (mu *RustfsUploader) Upload(ctx context.Context, filename string, r io.Reader, size int64) (string, error) {
	contentType, rdr, err := media.DetectContentType(r)
	if err != nil {
		contentType = "application/octet-stream"
	}

	objectName := filepath.Clean(strings.TrimSpace(filename))
	if _, err := mu.client.PutObject(ctx, mu.bucketName, objectName, rdr, size, minio.PutObjectOptions{
		ContentType: contentType,
	}); err != nil {
		return "", fmt.Errorf("failed to upload object: %w", err)
	}

	return objectName, nil
}

func (mu *RustfsUploader) Delete(ctx context.Context, filename string) error {
	objectName := filepath.Clean(strings.TrimSpace(filename))
	if err := mu.client.RemoveObject(ctx, mu.bucketName, objectName, minio.RemoveObjectOptions{}); err != nil {
		return fmt.Errorf("failed to delete object %s: %w", objectName, err)
	}

	return nil
}

func (mu *RustfsUploader) URL(name string, kind string) (string, error) {
	if kind == media.StorageKindExternal {
		return name, nil
	}

	url, err := url.JoinPath(mu.apiURL, mu.bucketName, name)
	if err != nil {
		return "", fmt.Errorf("failed to resolve file url - %s: %w", name, err)
	}

	return url, nil
}

func (mu *RustfsUploader) setPolicy(ctx context.Context) error {
	policy := `{
	"Version": "2012-10-17",
	"Statement": [
		{
			"Effect": "Allow",
			"Principal": {"AWS": ["*"]},
			"Action": ["s3:GetObject"],
			"Resource": ["arn:aws:s3:::fluffly/*"]
		}
	]
}`

	return mu.client.SetBucketPolicy(ctx, mu.bucketName, policy)
}

func (mu *RustfsUploader) ensureBucket(ctx context.Context) error {
	if err := mu.client.MakeBucket(ctx, mu.bucketName, minio.MakeBucketOptions{}); err != nil {
		exists, errBucketExists := mu.client.BucketExists(ctx, mu.bucketName)
		if errBucketExists != nil || !exists {
			return fmt.Errorf("failed to ensure bucket %q: %w", mu.bucketName, err)
		}
	}

	return nil
}
