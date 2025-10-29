package petfinder

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/dankobg/fluffly/config"
	"github.com/dankobg/fluffly/dbg"
	"github.com/dankobg/fluffly/httpserver"
)

type AuthCmd struct{}

func (dc *AuthCmd) Run() error {
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

	auth, err := client.authenticate(ctx)
	if err != nil {
		return fmt.Errorf("failed to get access token: %w", err)
	}

	dbg.PrintJSON(auth)

	return nil
}
