package petfinder

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"sync"
	"time"

	"github.com/dankobg/fluffly/ptr"
	"golang.org/x/time/rate"
)

const (
	rateLimitReqsPerDay    = 1_000 // pf free tier: 1_000 per day
	rateLimitReqsPerSecond = 40    // pf free tier: 50 per day
	rateLimitBurst         = 8
)

type rateLimitedClient struct {
	pf        *petfinderClient
	limiter   *rate.Limiter
	maxPerDay int
	mu        sync.Mutex
	requests  int
	lastReset time.Time
}

func newRateLimitedClient(pf *petfinderClient) *rateLimitedClient {
	return &rateLimitedClient{
		pf:        pf,
		limiter:   rate.NewLimiter(rate.Limit(rateLimitReqsPerSecond), rateLimitBurst),
		maxPerDay: rateLimitReqsPerDay,
		lastReset: time.Now(),
	}
}

func (client *rateLimitedClient) downloadOrganizationPage(ctx context.Context, filePath string, params listOrganizationsQueryParams) error {
	client.mu.Lock()
	if time.Since(client.lastReset) > 24*time.Hour {
		client.requests = 0
		client.lastReset = time.Now()
	}
	if client.requests >= client.maxPerDay {
		client.mu.Unlock()
		return fmt.Errorf("daily request limit (%d) reached", client.maxPerDay)
	}
	client.requests++
	client.mu.Unlock()

	if err := client.limiter.Wait(ctx); err != nil {
		return fmt.Errorf("rate limiter wait: %w", err)
	}

	if err := downloadOrganizationPage(ctx, client.pf, filePath, params); err != nil {
		return fmt.Errorf("download organization page: %w", err)
	}

	return nil
}

func (client *rateLimitedClient) downloadAnimalPage(ctx context.Context, filePath string, params listAnimalsQueryParams) error {
	client.mu.Lock()
	if time.Since(client.lastReset) > 24*time.Hour {
		client.requests = 0
		client.lastReset = time.Now()
	}
	if client.requests >= client.maxPerDay {
		client.mu.Unlock()
		return fmt.Errorf("daily request limit (%d) reached", client.maxPerDay)
	}
	client.requests++
	client.mu.Unlock()

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
	defer f.Close()

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

// downloadOrganizationPage fetches organizations page data and downloads it to a file `page_{number}.json`
func downloadOrganizationPage(ctx context.Context, pf *petfinderClient, filePath string, params listOrganizationsQueryParams) error {
	resp, err := pf.listOrganizations(ctx, params)
	if err != nil {
		return fmt.Errorf("list organizations: %w", err)
	}
	defer resp.Body.Close()

	out, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	return nil
}

// downloadOrganizationPages downloads missing organization pages concurrently
func downloadOrganizationPages(ctx context.Context, client *rateLimitedClient, dir string, missing []int, limit int) error {
	if len(missing) == 0 {
		fmt.Println("no pages to download")
		return nil
	}
	fmt.Println("pages to download: ", missing)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create output dir: %s: %w", dir, err)
	}

	sem := make(chan struct{}, limit)
	var wg sync.WaitGroup
	var mu sync.Mutex
	var errs []error

	for _, page := range missing {
		wg.Add(1)
		sem <- struct{}{}

		go func(page int) {
			defer wg.Done()
			defer func() { <-sem }()

			filePath := filepath.Join(dir, fmt.Sprintf("page_%d.json", page))
			err := client.downloadOrganizationPage(ctx, filePath, listOrganizationsQueryParams{
				Page:  &page,
				Limit: ptr.Of(100),
			})
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

// downloadAnimalPage fetches animals page data and downloads it to a file `page_{number}.json`
func downloadAnimalPage(ctx context.Context, pf *petfinderClient, filePath string, params listAnimalsQueryParams) error {
	resp, err := pf.listAnimals(ctx, params)
	if err != nil {
		return fmt.Errorf("list animals: %w", err)
	}
	defer resp.Body.Close()

	out, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	return nil
}

// downloadAnimalPages downloads missing animal pages concurrently
func downloadAnimalPages(ctx context.Context, client *rateLimitedClient, dir string, missing []int, limit int) error {
	if len(missing) == 0 {
		fmt.Println("no pages to download")
		return nil
	}
	fmt.Println("pages to download: ", missing)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create output dir: %s: %w", dir, err)
	}

	sem := make(chan struct{}, limit)
	var wg sync.WaitGroup
	var mu sync.Mutex
	var errs []error

	for _, page := range missing {
		wg.Add(1)
		sem <- struct{}{}

		go func(page int) {
			defer wg.Done()
			defer func() { <-sem }()

			filePath := filepath.Join(dir, fmt.Sprintf("page_%d.json", page))
			err := client.downloadAnimalPage(ctx, filePath, listAnimalsQueryParams{
				Page:  &page,
				Limit: ptr.Of(100),
			})
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
