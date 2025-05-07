package nodeapi

type NodeDeployRequest struct {
	DeploymentUUID string `json:"deploymentUUID"`
	DeploymentName string `json:"deploymentName"`
	NodeUUID       string `json:"nodeUID"`
	Replicas       int    `json:"replicas"`
	Image          string `json:"image"`
}
