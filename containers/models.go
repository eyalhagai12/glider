package containers

type Container struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	DeploymentUUID string `json:"deploymentUUID"`
	NodeID         string `json:"nodeID"`
	IP             string `json:"ip"`
	Port           string `json:"port"`
}
