package server

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
	"time"

	openapi_types "github.com/oapi-codegen/runtime/types"
)

type uploadResult struct {
	Name string
	Err  error
}

type uploadSource interface {
	Name() string
	Open() (r io.Reader, size int64, contentType string, e error)
}

type multipartSource struct {
	fh *multipart.FileHeader
}

func (m multipartSource) Name() string {
	return m.fh.Filename
}
func (m multipartSource) Open() (io.Reader, int64, string, error) {
	var f openapi_types.File
	f.InitFromMultipart(m.fh)
	rdr, err := f.Reader()
	if err != nil {
		return nil, 0, "", err
	}
	defer rdr.Close()
	return rdr, f.FileSize(), "", nil
}

type urlSource struct {
	c   *http.Client
	url string
}

func (u urlSource) Name() string {
	return filepath.Clean(filepath.Base(strings.TrimSpace(u.url)))
}
func (u urlSource) Open() (io.Reader, int64, string, error) {
	resp, err := u.c.Get(u.url)
	if err != nil {
		return nil, 0, "", err
	}
	defer resp.Body.Close()
	return resp.Body, resp.ContentLength, resp.Header.Get("Content-Type"), nil
}

func (a *ApiHandler) uploadOrganizationFiles(ctx context.Context, sources []uploadSource, workers int) []uploadResult {
	jobs := make(chan uploadSource)
	results := make(chan uploadResult)
	var wg sync.WaitGroup

	for range workers {
		wg.Go(func() {
			for src := range jobs {
				rdr, size, _, err := src.Open()
				if err != nil {
					results <- uploadResult{Err: fmt.Errorf("file %s: %w", src.Name(), err)}
					continue
				}
				url, err := a.uploadOrganizationFile(ctx, src.Name(), rdr, size)
				if err != nil {
					results <- uploadResult{Err: fmt.Errorf("upload %s: %w", src.Name(), err)}
				} else {
					results <- uploadResult{Name: url}
				}
			}
		})
	}

	go func() {
		for _, src := range sources {
			jobs <- src
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
