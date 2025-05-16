package main

import (
	"testing"
)

func TestSystemProcessProvider_ListProcesses_NoError(t *testing.T) {
	provider := &SystemProcessProvider{}
	_, err := provider.ListProcesses()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

}

func TestSystemProcessProvider_ListProcesses_GoIsRunning(t *testing.T) {
	provider := &SystemProcessProvider{}
	processes, err := provider.ListProcesses()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	found := false
	for _, process := range processes {
		if process.Name == "go" || process.Name == "go.exe" {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Expected to find a process named 'go', but it was not found")
	}
}

func TestSystemProcessProvider_ListProcesses_NonExistingProcess(t *testing.T) {
	provider := &SystemProcessProvider{}
	processes, err := provider.ListProcesses()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	for _, process := range processes {
		dummyProcessName := "non-existing-process"
		if process.Name == dummyProcessName {
			t.Errorf("Did not expect to find a process named '%s', but it was found", dummyProcessName)
		}
	}
}
