package main

import (
	"bytes"
	"testing"
)

var processes = []ProcessInfo{
	{PID: 123, Name: "init", Username: "root", CPUPercent: 2.0, MemBytes: 1024},
	{PID: 456, Name: "bash", Username: "luk", CPUPercent: 4.5, MemBytes: 125},
}

func TestJSONFormatter_Format(t *testing.T) {
	formatter := &JSONFormatter{}

	var buf bytes.Buffer
	err := formatter.Format(processes, &buf)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	actual := buf.String()
	expected := `[
  {
    "PID": 123,
    "Name": "init",
    "Username": "root",
    "CPUPercent": 2,
    "MemBytes": 1024
  },
  {
    "PID": 456,
    "Name": "bash",
    "Username": "luk",
    "CPUPercent": 4.5,
    "MemBytes": 125
  }
]
`
	// this comparison is naive to be improved
	if actual != expected {
		t.Errorf("Expected %s but got: %s", expected, actual)
	}
}

func TestCSVFormatter_Format(t *testing.T) {
	formatter := &CSVFormatter{}

	var buf bytes.Buffer
	err := formatter.Format(processes, &buf)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	actual := buf.String()
	// this comparison is naive to be improved
	expected := `PID,Name,User,CPU (%),Memory (Bytes)
123,init,root,2.00,1024
456,bash,luk,4.50,125
`
	if actual != expected {
		t.Errorf("Expected %s but got: %s", expected, actual)
	}
}

func TestGetFormatter_ValidFormats(t *testing.T) {
	f, err := getFormatter("json")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if _, ok := f.(*JSONFormatter); !ok {
		t.Error("Expected JSONFormatter")
	}

	f, err = getFormatter("csv")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if _, ok := f.(*CSVFormatter); !ok {
		t.Error("Expected CSVFormatter")
	}
}

func TestGetFormatter_InvalidFormat(t *testing.T) {
	_, err := getFormatter("xml")
	if err == nil {
		t.Error("Expected error for unsupported format")
	}
}
