package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

// Helper to capture stdout and stderr
func captureOutput(f func()) (stdout string, stderr string) {
	oldOut := os.Stdout
	oldErr := os.Stderr

	rOut, wOut, _ := os.Pipe()
	rErr, wErr, _ := os.Pipe()

	os.Stdout = wOut
	os.Stderr = wErr

	outC := make(chan string)
	errC := make(chan string)

	// Read stdout
	go func() {
		var buf bytes.Buffer
		_, _ = buf.ReadFrom(rOut)
		outC <- buf.String()
	}()

	// Read stderr
	go func() {
		var buf bytes.Buffer
		_, _ = buf.ReadFrom(rErr)
		errC <- buf.String()
	}()

	f()

	_ = wOut.Close()
	_ = wErr.Close()

	os.Stdout = oldOut
	os.Stderr = oldErr

	stdout = <-outC
	stderr = <-errC
	return
}

func TestMain_ValidArgsStdout(t *testing.T) {
	args := []string{"--format", "json"}

	stdout, stderr := captureOutput(func() {
		err := run(args, os.Stdout, os.Stderr)
		if err != nil {
			t.Errorf("run() returned error: %v", err)
		}
	})

	if stderr != "" {
		t.Errorf("unexpected stderr output: %q", stderr)
	}

	if !strings.Contains(stdout, "[") {
		t.Errorf("expected JSON output, got: %s", stdout)
	}
}

func TestMain_InvalidFlag(t *testing.T) {
	args := []string{"--unknownflag"}

	_, stderr := captureOutput(func() {
		err := run(args, os.Stdout, os.Stderr)
		if err == nil {
			t.Error("expected error from invalid flag, got nil")
		}
	})

	if !strings.Contains(stderr, "flag provided but not defined") {
		t.Errorf("unexpected stderr output: %q", stderr)
	}
}
