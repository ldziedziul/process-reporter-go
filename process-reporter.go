package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/shirou/gopsutil/v3/process"
)

type ProcessInfo struct {
	PID        int32   `json:"pid"`
	Name       string  `json:"name"`
	Username   string  `json:"user"`
	CPUPercent float64 `json:"cpu_percent"`
	MemBytes   uint64  `json:"mem_bytes"`
}

func getProcessInfo() ([]ProcessInfo, error) {
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
		mem, _ := proc.MemoryInfo()
		memBytes := uint64(0)
		if mem != nil {
			memBytes = mem.RSS
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

func writeJSON(processes []ProcessInfo, file *os.File) error {
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(processes)
}

func writeCSV(processes []ProcessInfo, file *os.File) error {
	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"PID", "Name", "User", "CPU (%)", "Memory (bytes)"})
	for _, proc := range processes {
		writer.Write([]string{
			strconv.Itoa(int(proc.PID)),
			proc.Name,
			proc.Username,
			fmt.Sprintf("%.2f", proc.CPUPercent),
			strconv.FormatUint(proc.MemBytes, 10),
		})
	}
	return nil
}

func main() {
	format := flag.String("format", "json", "Output format: json or csv")
	output := flag.String("output", "", "Optional output file name. If not set, output goes to stdout")
	flag.Parse()

	processes, err := getProcessInfo()
	if err != nil {
		log.Fatalf("Failed to get process info: %v", err)
	}

	var out *os.File
	if *output == "" {
		out = os.Stdout
	} else {
		out, err = os.Create(*output)
		if err != nil {
			log.Fatalf("Failed to create output file: %v", err)
		}
		defer out.Close()
	}

	switch *format {
	case "json":
		err = writeJSON(processes, out)
	case "csv":
		err = writeCSV(processes, out)
	default:
		log.Fatalf("Unsupported format: %s", *format)
	}

	if err != nil {
		log.Fatalf("Failed to write report: %v", err)
	}

	if *output != "" {
		fmt.Printf("Report saved successfully in %s format to %s\n", *format, *output)
	}
}
