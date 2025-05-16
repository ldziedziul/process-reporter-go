package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
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

	writer.Write([]string{"PID", "Name", "User", "CPU (%)", "Memory (Bytes)"})
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

func getFormatter(format string) (Formatter, error) {
	switch format {
	case "json":
		return &JSONFormatter{}, nil
	case "csv":
		return &CSVFormatter{}, nil
	default:
		return nil, errors.New("unsupported format")
	}
}
