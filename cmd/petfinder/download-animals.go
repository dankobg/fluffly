package petfinder

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/dankobg/fluffly/config"
	"github.com/dankobg/fluffly/httpserver"
)

type DownloadAnimalsCmd struct {
	Dir   string `required:"" type:"path"`
	Min   int    `default:"1"`
	Max   int    `default:"10"`
	Limit int    `default:"10"`
}

func (dc *DownloadAnimalsCmd) Run() error {
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

	if err := downloadAnimalPages(context.Background(), client, dc.Dir, missingPages, dc.Limit); err != nil {
		return fmt.Errorf("failed to download animal pages: %w", err)
	}

	return nil
}
