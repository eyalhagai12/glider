package nodeapi

type NodeDeployRequest struct {
	DeploymentUUID string `json:"deploymentUUID"`
	DeploymentName string `json:"deploymentName"`
	NodeUUID       string `json:"nodeUID"`
	Replicas       int    `json:"replicas"`
	Image          string `json:"image"`
}

type NodeConnectRequest struct {
	Interface           string   `json:"interface"`
	IP                  string   `json:"ip"`
	PublicKey           string   `json:"publicKey"`
	Endpoint            string   `json:"endpoint"`
	AllowedIPs          []string `json:"allowedIps"`
	PresharedKey        string   `json:"presharedKey,omitempty"`
	PersistentKeepalive int      `json:"keepalive,omitempty"`
}
