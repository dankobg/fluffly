package petfinder

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"sync"

	"github.com/dankobg/fluffly/config"
	"github.com/dankobg/fluffly/httpserver"
	"github.com/dankobg/fluffly/ptr"
)

type DownloadOrganizationsCmd struct {
	Dir string `required:"" type:"path"`
	Min int    `default:"1"`
	Max int    `default:"10"`
}

func (dc *DownloadOrganizationsCmd) Run() error {
	cfg, _, err := config.New()
	if err != nil {
		slog.Error("failed to initialize config", slog.Any("error", err))
		return err
	}

	httpc := httpserver.NewHttpClient()
	client, err := newPetfinderClient(httpc, cfg.Petfinder)
	if err != nil {
		return fmt.Errorf("failed to create petfinder client: %w", err)
	}

	ctx := context.Background()

	missingPages, e := getMissingPages(dc.Dir, dc.Min, dc.Max)
	if e != nil {
		return fmt.Errorf("failed to get missing pages: %w", e)
	}

	if err := downloadOrganizationPages(ctx, client, dc.Dir, missingPages, 10); err != nil {
		return fmt.Errorf("failed to download organization pages: %w", err)
	}

	return nil
}

var rePage = regexp.MustCompile(`^page_(\d+)\.json$`)

// getMissingPages returns missing pages in a dir with
// `page_{number}.json` pattern from min to max page
func getMissingPages(dir string, minPage, maxPage int) ([]int, error) {
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

func downloadOrganizationPage(ctx context.Context, client *petfinderClient, filePath string, params listOrganizationsQueryParams) error {
	resp, err := client.listOrganizations(ctx, params)
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

// download missing pages concurrently
func downloadOrganizationPages(ctx context.Context, client *petfinderClient, dir string, missing []int, limit int) error {
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
			err := downloadOrganizationPage(ctx, client, filePath, listOrganizationsQueryParams{
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
