package wireguard

import (
	"context"
	"database/sql"
	"fmt"
	backend "glider"
	"log/slog"
	"os"
	"os/exec"

	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

type WireGuardService struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewWireGuardService(db *sql.DB, logger *slog.Logger) *WireGuardService {
	return &WireGuardService{
		db:     db,
		logger: logger,
	}
}

func (wgs *WireGuardService) Create(ctx context.Context, network *backend.Network) (*backend.Network, error) {
	wgs.logger.Info("Creating WireGuard network", "network", network)

	wgConfig, err := createWireguradConfig(network)
	if err != nil {
		wgs.logger.Error("Failed to create WireGuard config", "error", err)
		return nil, err
	}

	err = storeConfigInFile(wgConfig, network)
	if err != nil {
		wgs.logger.Error("Failed to store WireGuard config in file", "error", err)
		return nil, err
	}

	err = exec.Command("wg-quick", "up", network.Name).Run()
	if err != nil {
		wgs.logger.Error("Failed to create WireGuard network", "error", err)
		return nil, err
	}

	wgs.logger.Info("WireGuard network created successfully", "network", network)
	return network, nil
}

func createWireguradConfig(network *backend.Network) (string, error) {
	key, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		fmt.Println("Failed to generate private key:", err)
		return "", err
	}

	return fmt.Sprintf(`
		[Interface]
		Address = 192.168.1.1/24
		ListenPort = 51821
		PrivateKey = %s
	`, key.String()), nil
}

func storeConfigInFile(config string, network *backend.Network) error {
	filePath := fmt.Sprintf("/etc/wireguard/%s.conf", network.Name)
	err := os.WriteFile(filePath, []byte(config), 0600)
	if err != nil {
		fmt.Println("Failed to write WireGuard config to file:", err)
		return err
	}
	fmt.Println("WireGuard config written to", filePath)
	return nil
}
