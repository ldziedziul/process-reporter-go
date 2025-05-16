package main

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

type MockProcessProvider struct {
	Processes []ProcessInfo
	Err       error
}

var mockProcessProvider = MockProcessProvider{
	Processes: []ProcessInfo{
		{PID: 1, Name: "init", Username: "root", CPUPercent: 2.0, MemBytes: 1024},
		{PID: 1, Name: "bash", Username: "luk", CPUPercent: 4.5, MemBytes: 125},
	},
}

func (m *MockProcessProvider) ListProcesses() ([]ProcessInfo, error) {
	return m.Processes, m.Err
}

func TestApp_Run_JSON(t *testing.T) {
	var buf bytes.Buffer
	app := App{Provider: &mockProcessProvider, Formatter: &JSONFormatter{}, Output: &buf}

	err := app.Run()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	output := buf.String()
	if !strings.Contains(output, "init") {
		t.Errorf("Expected output to contain 'init', got %s", output)
	}
}

func TestApp_Run_CSV(t *testing.T) {
	var buf bytes.Buffer
	app := App{Provider: &mockProcessProvider, Formatter: &CSVFormatter{}, Output: &buf}

	err := app.Run()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	output := buf.String()
	if !strings.Contains(output, "bash") {
		t.Errorf("Expected output to contain 'bash', got %s", output)
	}
}

func TestApp_Run_Error(t *testing.T) {
	provider := &MockProcessProvider{Err: errors.New("mock failure")}
	var buf bytes.Buffer
	app := App{Provider: provider, Formatter: &JSONFormatter{}, Output: &buf}

	err := app.Run()
	if err == nil || !strings.Contains(err.Error(), "mock failure") {
		t.Errorf("Expected mock failure error, got %v", err)
	}
}
