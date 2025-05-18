package network

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"

	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

func storePrivateKey(interfaceName string, privateKey wgtypes.Key) (string, error) {
	filePath := fmt.Sprintf("./tmp/keys/%s_private.key", interfaceName)
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = file.WriteString(privateKey.String())
	if err != nil {
		return "", err
	}
	return filePath, nil
}

func generateKeyPair(interfaceName string) (wgtypes.Key, string, error) {
	privateKey, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		return wgtypes.Key{}, "", err
	}
	publicKey := privateKey.PublicKey()

	prvKeyFilePath, err := storePrivateKey(interfaceName, privateKey)
	if err != nil {
		return wgtypes.Key{}, "", err
	}

	return publicKey, prvKeyFilePath, nil
}

func InitializeVPN(logger *slog.Logger, net *Network) error {
	interfaceName := net.InterfaceName
	ipAddr := net.IpAddress
	port := net.Port

	logger.Info("Initializing WireGuard VPN interface", "interface", interfaceName)
	if out, err := exec.Command("ip", "link", "add", interfaceName, "type", "wireguard").CombinedOutput(); err != nil {
		logger.Error("Failed to create WireGuard interface", "output", string(out), "error", err)
		return err
	}

	logger.Info("Generating WireGuard key pair")
	publicKey, privateKeyFilePath, err := generateKeyPair(interfaceName)
	if err != nil {
		return err
	}
	logger.Info("Public key generated", "publicKey", publicKey.String())

	logger.Info("Setting up ip address for WireGuard interface")
	if out, err := exec.Command("ip", "addr", "add", ipAddr, "dev", interfaceName).CombinedOutput(); err != nil {
		logger.Error("Failed to set IP address", "output", string(out), "error", err)
		return err
	}

	logger.Info("Setting up WireGuard interface with private key and port")
	if out, err := exec.Command("wg", "set", interfaceName, "private-key", privateKeyFilePath, "listen-port", port).CombinedOutput(); err != nil {
		logger.Error("Failed to set WireGuard interface", "output", string(out), "error", err)
		return err
	}

	logger.Info("Starting WireGuard interface")
	if out, err := exec.Command("ip", "link", "set", interfaceName, "up").CombinedOutput(); err != nil {
		logger.Error("Failed to start WireGuard interface", "output", string(out), "error", err)
		return err
	}

	logger.Info("WireGuard VPN interface initialized", "interface", interfaceName)

	return nil
}
