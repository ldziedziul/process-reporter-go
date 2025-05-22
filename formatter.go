package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
)

type Formatter interface {
	Format([]ProcessInfo, io.Writer) error
}

type JSONFormatter struct{}

func (j *JSONFormatter) Format(processes []ProcessInfo, w io.Writer) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(processes)
}

type CSVFormatter struct{}

func (c *CSVFormatter) Format(processes []ProcessInfo, w io.Writer) error {
	writer := csv.NewWriter(w)
	defer writer.Flush()

	err := writer.Write([]string{"PID", "Name", "User", "CPU (%)", "Memory (Bytes)"})
	if err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}
	for _, proc := range processes {
		err := writer.Write([]string{
			strconv.Itoa(int(proc.PID)),
			proc.Name,
			proc.Username,
			fmt.Sprintf("%.2f", proc.CPUPercent),
			strconv.FormatUint(proc.MemBytes, 10),
		})
		if err != nil {
			return fmt.Errorf("failed to write CSV record for PID %d: %w", proc.PID, err)
		}
	}
	return nil
}

func getFormatter(format string) (Formatter, error) {
	switch format {
	case "json":
		return &JSONFormatter{}, nil
	case "csv":
		return &CSVFormatter{}, nil
	default:
		return nil, fmt.Errorf("unsupported format: %v", format)
	}
}
