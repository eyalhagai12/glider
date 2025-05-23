package nodes

import (
	"bytes"
	"encoding/json"
	"glider/network"
	"net/http"
	"strings"
)

func (n Node) AddNetwork(net *network.Network, publicKey string, allowedIPs []string) error {
	allowedIPs = append(allowedIPs, net.IpAddress+"/32")
	allowedIPsStr := strings.Join(allowedIPs, ",")

	requestData := network.GetVPNConnectionRequest(net.InterfaceName, net.IpAddress, net.Port, publicKey, allowedIPsStr)
	data, err := json.Marshal(requestData)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, n.ConnectionURL, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	return nil
}
