package network

import (
	"fmt"
	"log/slog"

	"github.com/hibare/headscale-client-go"
)

func ConnectToHeadscale(serverURL string, apiKey string, logger *slog.Logger) (headscale.HeadscaleClientInterface, error) {
	client, err := headscale.NewClient(serverURL, apiKey, headscale.HeadscaleClientOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to create Headscale client: %w", err)
	}

	return client, nil
}
