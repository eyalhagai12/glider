package resources

import (
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

func FetchDisks() ([]Disk, error) {
	partitions, err := disk.Partitions(false)
	if err != nil {
		return nil, err
	}

	disks := make([]Disk, 0, len(partitions))
	for _, partition := range partitions {
		usage, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			return nil, err
		}
		disks = append(disks, Disk{
			Free:  float32(usage.Free) / (1024 * 1024 * 1024),
			Used:  float32(usage.Used) / (1024 * 1024 * 1024),
			Total: float32(usage.Total) / (1024 * 1024 * 1024),
			Name:  partition.Device,
			Mount: partition.Mountpoint,
		})
	}
	return disks, nil
}

func FetchCPUs() ([]CPU, error) {
	cpuUsageList, err := cpu.Percent(100*time.Nanosecond, false)
	if err != nil {
		return nil, err
	}
	if len(cpuUsageList) == 0 {
		return nil, nil
	}

	cpuUsage := make([]CPU, len(cpuUsageList))
	for i, usage := range cpuUsageList {
		cpuUsage[i] = CPU{
			Usage: usage,
			Temp:  0,
		}
	}

	return cpuUsage, nil
}

func FetchMemory() (Memory, error) {
	ram, err := mem.VirtualMemory()
	if err != nil {
		return Memory{}, err
	}

	return Memory{
		Free:  float32(ram.Available) / (1024 * 1024 * 1024),
		Used:  float32(ram.Used) / (1024 * 1024 * 1024),
		Total: float32(ram.Total) / (1024 * 1024 * 1024),
	}, nil
}

func FetchNodeUtilization() (NodeResources, error) {
	disks, err := FetchDisks()
	if err != nil {
		return NodeResources{}, err
	}
	
	cpuData, err := FetchCPUs()
	if err != nil {
		return NodeResources{}, err
	}

	memoryData, err := FetchMemory()
	if err != nil {
		return NodeResources{}, err
	}

	return NodeResources{
		CPUs:   cpuData,
		Memory: memoryData,
		Disks:  disks,
	}, nil
}
