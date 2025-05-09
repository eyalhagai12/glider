package resources

type CPU struct {
	Usage float64 `json:"usage"`
	Temp  float32 `json:"temp"`
}

type Disk struct {
	Free  float32 `json:"free"`
	Used  float32 `json:"used"`
	Total float32 `json:"total"`
	Name  string  `json:"name"`
	Mount string  `json:"mount"`
}

type Memory struct {
	Free  float32 `json:"free"`
	Used  float32 `json:"used"`
	Total float32 `json:"total"`
}

type NodeResources struct {
	CPUs   []CPU  `json:"cpus"`
	Memory Memory `json:"memory"`
	Disks  []Disk `json:"disks"`
}

type ContainerUtilization struct {
	ContainerID   string  `json:"container_id"`
	ContainerName string  `json:"container_name"`
	CPUUsage      float64 `json:"cpu_usage"`
	MemoryUsage   float64 `json:"memory_usage"`
	NetworkIO     float64 `json:"network_io"`
	DiskIO        float64 `json:"disk_io"`
}

type ResourceResponse struct {
	NodeResources NodeResources          `json:"nodeResources"`
	Containers    []ContainerUtilization `json:"containersUtilization"`
}
