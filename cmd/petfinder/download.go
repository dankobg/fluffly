package petfinder

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"sync"

	"golang.org/x/time/rate"
)

const (
	rateLimitReqsPerDay    = 1_000 // pf free tier: 1_000 per day
	rateLimitReqsPerSecond = 35    // pf free tier: 50 per day
	rateLimitBurst         = 5
)

type rateLimitedClient struct {
	pf              *petfinderClient
	limiter         *rate.Limiter
	limitReqsPerDay int
}

func newRateLimitedClient(pf *petfinderClient) *rateLimitedClient {
	return &rateLimitedClient{
		pf:              pf,
		limiter:         rate.NewLimiter(rate.Limit(rateLimitReqsPerSecond), rateLimitBurst),
		limitReqsPerDay: rateLimitReqsPerDay,
	}
}

func (client *rateLimitedClient) downloadOrganizationPage(ctx context.Context, filePath string, params listOrganizationsQueryParams) error {
	if err := client.limiter.Wait(ctx); err != nil {
		return fmt.Errorf("rate limiter wait: %w", err)
	}

	if err := downloadOrganizationPage(ctx, client.pf, filePath, params); err != nil {
		return fmt.Errorf("download organization page: %w", err)
	}

	return nil
}

func (client *rateLimitedClient) downloadAnimalPage(ctx context.Context, filePath string, params listAnimalsQueryParams) error {
	if err := client.limiter.Wait(ctx); err != nil {
		return fmt.Errorf("rate limiter wait: %w", err)
	}

	if err := downloadAnimalPage(ctx, client.pf, filePath, params); err != nil {
		return fmt.Errorf("download animal page: %w", err)
	}

	return nil
}

var rePage = regexp.MustCompile(`^page_(\d+)\.json$`)

// getMissingPageNumbers returns missing page numbers in a given dir with
// `page_{number}.json` pattern from min to max page
func getMissingPageNumbers(dir string, minPage, maxPage int) ([]int, error) {
	f, err := os.Open(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to open organizations directory: %w", err)
	}

	defer func() { _ = f.Close() }()

	fileNames, err := f.Readdirnames(-1)
	if err != nil {
		return nil, fmt.Errorf("failed to read dir names: %w", err)
	}

	var pages []int

	for _, name := range fileNames {
		m := rePage.FindStringSubmatch(name)
		if len(m) == 2 {
			num, err := strconv.Atoi(m[1])
			if err != nil {
				continue
			}

			pages = append(pages, num)
		}
	}

	if len(pages) > 2 {
		slices.Sort(pages)
	}

	existing := make(map[int]bool)
	for _, p := range pages {
		existing[p] = true
	}

	var missing []int

	for i := minPage; i <= maxPage; i++ {
		if !existing[i] {
			missing = append(missing, i)
		}
	}

	return missing, nil
}

// writeToFile writes response to a file
func writeToFile(resp *http.Response, filePath string) error {
	out, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}

	defer func() { _ = out.Close() }()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	return nil
}

// downloadOrganizationPage fetches organizations page data and downloads it to a file
func downloadOrganizationPage(ctx context.Context, pf *petfinderClient, filePath string, params listOrganizationsQueryParams) error {
	resp, err := pf.listOrganizations(ctx, params)
	if err != nil {
		return fmt.Errorf("list organizations: %w", err)
	}

	defer func() { _ = resp.Body.Close() }()

	return writeToFile(resp, filePath)
}

// downloadAnimalPage fetches animals page data and downloads it to a file
func downloadAnimalPage(ctx context.Context, pf *petfinderClient, filePath string, params listAnimalsQueryParams) error {
	resp, err := pf.listAnimals(ctx, params)
	if err != nil {
		return fmt.Errorf("list animals: %w", err)
	}

	defer func() { _ = resp.Body.Close() }()

	return writeToFile(resp, filePath)
}

// downloadPages downloads pages concurrently
func downloadsPages(
	ctx context.Context,
	dir string,
	pages []int,
	concurrencyLimit int,
	fn func(ctx context.Context, fpath string, page, limit int) error,
) error {
	if len(pages) == 0 {
		fmt.Println("no pages to download")
		return nil
	}

	fmt.Println("pages to download: ", pages)

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("failed to create output dir: %s: %w", dir, err)
	}

	sem := make(chan struct{}, concurrencyLimit)

	var (
		wg   sync.WaitGroup
		mu   sync.Mutex
		errs []error
	)

	for _, page := range pages {
		wg.Add(1)

		sem <- struct{}{}

		go func(page int) {
			defer func() {
				<-sem
				wg.Done()
			}()

			filePath := filepath.Join(dir, fmt.Sprintf("page_%d.json", page))

			err := fn(ctx, filePath, page, 100)
			if err != nil {
				mu.Lock()

				errs = append(errs, fmt.Errorf("page %d: %w", page, err))
				mu.Unlock()

				return
			}

			fmt.Printf("saved page: %d\n", page)
		}(page)
	}

	wg.Wait()

	if err := errors.Join(errs...); err != nil {
		return fmt.Errorf("some downloads failed: %w", err)
	}

	return nil
}

func downloadOrganizationPages(ctx context.Context, client *rateLimitedClient, dir string, pages []int, concurrencyLimit int) error {
	return downloadsPages(ctx, dir, pages, concurrencyLimit, func(ctx context.Context, filePath string, page, limit int) error {
		return client.downloadOrganizationPage(ctx, filePath, listOrganizationsQueryParams{
			Page:  &page,
			Limit: new(limit),
		})
	})
}

func downloadAnimalPages(ctx context.Context, client *rateLimitedClient, dir string, pages []int, concurrencyLimit int) error {
	return downloadsPages(ctx, dir, pages, concurrencyLimit, func(ctx context.Context, filePath string, page, limit int) error {
		return client.downloadAnimalPage(ctx, filePath, listAnimalsQueryParams{
			Page:  &page,
			Limit: new(limit),
		})
	})
}
