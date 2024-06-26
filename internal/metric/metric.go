package metric

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

type Metric struct {
	MemUsed float64 `json:"memUsed"`
	CpuUsed float64 `json:"cpuUsed"`
}

func GetMachineMetric() (*Metric, error) {
	vm, err := mem.VirtualMemory()
	if err != nil {
		return nil, fmt.Errorf("get virtual memory -> %w", err)
	}
	totalPercent, err := cpu.Percent(time.Second, false)
	if err != nil {
		return nil, fmt.Errorf("get cpu percent -> %w", err)
	}
	return &Metric{
		MemUsed: vm.UsedPercent,
		CpuUsed: totalPercent[0],
	}, nil
}
