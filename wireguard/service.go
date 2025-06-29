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

	out, err := exec.Command("wg-quick", "up", network.Name).CombinedOutput()
	if err != nil {
		wgs.logger.Error("Failed to create WireGuard network", "error", err, "output", string(out))
		return nil, err
	}

	wgs.logger.Info("WireGuard network created successfully", "output", string(out))

	wgs.logger.Info("WireGuard network created successfully", "network", network)
	return network, nil
}

func (wgs *WireGuardService) GenerateKeys() (string, string, error) {
	privateKey, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		wgs.logger.Error("Failed to generate WireGuard private key", "error", err)
		return "", "", err
	}

	publicKey := privateKey.PublicKey()
	wgs.logger.Info("Generated WireGuard keys", "privateKey", privateKey.String(), "publicKey", publicKey.String())
	return privateKey.String(), publicKey.String(), nil
}

func createWireguradConfig(network *backend.Network) (string, error) {
	return fmt.Sprintf(`
	[Interface]
	Address = %s/24
	ListenPort = %d
	PrivateKey = %s
	`, network.Address, network.ListenPort, network.PrivateKey), nil
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
