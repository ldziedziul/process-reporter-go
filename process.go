package main

import (
	"github.com/shirou/gopsutil/v3/process"
	"math"
)

type ProcessInfo struct {
	PID        int32
	Name       string
	Username   string
	CPUPercent float64
	MemBytes   uint64
}

type ProcessProvider interface {
	ListProcesses() ([]ProcessInfo, error)
}

func round(val float64, precision int) float64 {
	pow := math.Pow(10, float64(precision))
	return math.Round(val*pow) / pow
}

type SystemProcessProvider struct{}

func (p *SystemProcessProvider) ListProcesses() ([]ProcessInfo, error) {
	pids, err := process.Pids()
	if err != nil {
		return nil, err
	}

	var processes []ProcessInfo
	for _, pid := range pids {
		proc, err := process.NewProcess(pid)
		if err != nil {
			continue
		}
		name, _ := proc.Name()
		user, _ := proc.Username()
		cpu, _ := proc.CPUPercent()
		cpu = round(cpu, 2)
		memInfo, _ := proc.MemoryInfo()
		memBytes := uint64(0)
		if memInfo != nil {
			memBytes = memInfo.RSS
		}
		processes = append(processes, ProcessInfo{
			PID:        pid,
			Name:       name,
			Username:   user,
			CPUPercent: cpu,
			MemBytes:   memBytes,
		})
	}
	return processes, nil
}
