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

type ResourceResponse struct {
	CPUs   []CPU  `json:"cpus"`
	Memory Memory `json:"memory"`
	Disks  []Disk `json:"disks"`
}
