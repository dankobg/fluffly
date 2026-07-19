package petfinder

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/dankobg/fluffly/config"
	"github.com/dankobg/fluffly/httpserver"
)

type DownloadOrganizationsCmd struct {
	Dir   string `required:"" type:"path" help:"Download dir"`
	Min   int    `default:"1" help:"Starting page (inclusive)"`
	Max   int    `default:"10" help:"Ending page (inclusive)"`
	Limit int    `default:"10" help:"Concurrency limit"`
}

func (dc *DownloadOrganizationsCmd) Run() error {
	cfg, _, err := config.New()
	if err != nil {
		slog.Error("failed to initialize config", slog.Any("error", err))
		return err
	}

	httpc := httpserver.NewHttpClient()

	pf, err := newPetfinderClient(httpc, cfg.Petfinder)
	if err != nil {
		return fmt.Errorf("failed to create petfinder client: %w", err)
	}

	client := newRateLimitedClient(pf)

	missingPages, e := getMissingPageNumbers(dc.Dir, dc.Min, dc.Max)
	if e != nil {
		return fmt.Errorf("failed to get missing pages: %w", e)
	}

	if err := downloadOrganizationPages(context.Background(), client, dc.Dir, missingPages, dc.Limit); err != nil {
		return fmt.Errorf("failed to download organization pages: %w", err)
	}

	return nil
}
