package network

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strings"

	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

func storePrivateKey(interfaceName string, privateKey wgtypes.Key, storageDirPath string) (string, error) {
	err := os.MkdirAll(storageDirPath, 0755)
	if err != nil {
		return "", err
	}

	savePath := fmt.Sprintf("%s/%s_private.key", storageDirPath, interfaceName)
	file, err := os.OpenFile(savePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = file.WriteString(privateKey.String())
	if err != nil {
		return "", err
	}
	return savePath, nil
}

func generateKeyPair(interfaceName string, prvKeyDirPath string) (wgtypes.Key, string, error) {
	privateKey, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		return wgtypes.Key{}, "", err
	}
	publicKey := privateKey.PublicKey()

	prvKeyFilePath, err := storePrivateKey(interfaceName, privateKey, prvKeyDirPath)
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
	publicKey, privateKeyFilePath, err := generateKeyPair(interfaceName, "./tmp/keys")
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

func ConnectToVPN(logger *slog.Logger, interfaceName string, ipAddress string, publicKey string, endpoint string, allowedIPs []string) error {
	logger.Info("Connecting to VPN", "interface", interfaceName, "ip", ipAddress, "endpoint", endpoint)
	if out, err := exec.Command("ip", "link", "add", "dev", interfaceName, "type", "wireguard").CombinedOutput(); err != nil {
		logger.Error("Failed to add interface to agent node", "output", string(out), "error", err)
		return err
	}

	if out, err := exec.Command("ip", "addr", "add", ipAddress, "dev", interfaceName).CombinedOutput(); err != nil {
		logger.Error("Failed to set IP address for interface", "output", string(out), "error", err)
		return err
	}

	_, prvKeyFile, err := generateKeyPair(interfaceName, "./.metadata")
	if err != nil {
		logger.Error("Failed to generate key pair", "error", err)
		return err
	}

	if out, err := exec.Command("wg", "set", interfaceName, "private-key", prvKeyFile, "listen-port", "51820").CombinedOutput(); err != nil {
		logger.Error("Failed to set WireGuard interface", "output", string(out), "error", err)
		return err
	}

	persistentKeepalive := 25 // seconds
	if out, err := exec.Command("wg", "set", interfaceName, "persistent-keepalive", fmt.Sprintf("%d", persistentKeepalive)).CombinedOutput(); err != nil {
		logger.Error("Failed to set WireGuard interface", "output", string(out), "error", err)
		return err
	}

	allowedIPs = append(allowedIPs, ipAddress)
	allowedIPsStr := strings.Join(allowedIPs, ",")
	if out, err := exec.Command("wg", "set", interfaceName, "peer", publicKey, "endpoint", endpoint, "allowed-ips", allowedIPsStr).CombinedOutput(); err != nil {
		logger.Error("Failed to set WireGuard peer", "output", string(out), "error", err)
		return err
	}
	if out, err := exec.Command("ip", "link", "set", interfaceName, "up").CombinedOutput(); err != nil {
		logger.Error("Failed to set WireGuard interface up", "output", string(out), "error", err)
		return err
	}

	logger.Info("WireGuard VPN connection established", "interface", interfaceName, "ip", ipAddress, "endpoint", endpoint)

	return nil
}
