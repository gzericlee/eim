package metric

import (
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
		return nil, err
	}
	totalPercent, err := cpu.Percent(time.Second, false)
	if err != nil {
		return nil, err
	}
	return &Metric{
		MemUsed: vm.UsedPercent,
		CpuUsed: totalPercent[0],
	}, nil
}
