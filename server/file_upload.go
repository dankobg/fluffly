package server

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"sync"
	"time"

	openapi_types "github.com/oapi-codegen/runtime/types"
)

type uploadResult struct {
	Name string
	Err  error
}

func (a *ApiHandler) uploadOrganizationMultipartFiles(ctx context.Context, fhs []*multipart.FileHeader, workers int) []uploadResult {
	jobs := make(chan *multipart.FileHeader)
	results := make(chan uploadResult)
	var wg sync.WaitGroup

	for range workers {
		wg.Go(func() {
			for fh := range jobs {
				var f openapi_types.File
				f.InitFromMultipart(fh)
				rdr, err := f.Reader()
				if err != nil {
					results <- uploadResult{Err: fmt.Errorf("file %s: %w", fh.Filename, err)}
					continue
				}
				url, err := a.uploadOrganizationFile(ctx, f.Filename(), rdr, f.FileSize())
				rdr.Close()
				if err != nil {
					results <- uploadResult{Err: fmt.Errorf("upload %s: %w", f.Filename(), err)}
				} else {
					results <- uploadResult{Name: url}
				}
			}
		})
	}

	go func() {
		for _, fh := range fhs {
			jobs <- fh
		}
		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	out := make([]uploadResult, 0)
	for res := range results {
		out = append(out, res)
	}
	return out
}

type deleteResult struct {
	Name string
	Err  error
}

func (a *ApiHandler) deleteOrganizationFiles(ctx context.Context, filenames []string, workers int) []deleteResult {
	jobs := make(chan string)
	results := make(chan deleteResult)
	var wg sync.WaitGroup

	for range workers {
		wg.Go(func() {
			for name := range jobs {
				err := a.uploader.Delete(ctx, name)
				if err != nil {
					results <- deleteResult{Err: fmt.Errorf("delete %s: %w", name, err)}
				} else {
					results <- deleteResult{Name: name}
				}
			}
		})
	}

	go func() {
		for _, name := range filenames {
			jobs <- name
		}
		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	out := make([]deleteResult, 0)
	for res := range results {
		out = append(out, res)
	}
	return out
}

func (a *ApiHandler) uploadOrganizationFile(ctx context.Context, filename string, r io.Reader, size int64) (string, error) {
	name := filepath.Join("organizations", fmt.Sprintf("%d-%s", time.Now().UnixNano(), filepath.Base(filename)))
	return a.uploader.Upload(ctx, name, r, size)
}

func (a *ApiHandler) uploadAnimalFile(ctx context.Context, filename string, r io.Reader, size int64) (string, error) {
	name := filepath.Join("animals", fmt.Sprintf("%d-%s", time.Now().UnixNano(), filepath.Base(filename)))
	return a.uploader.Upload(ctx, name, r, size)
}
