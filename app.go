package main

import (
	"fmt"
	"io"
)

type App struct {
	Provider  ProcessProvider
	Formatter Formatter
	Output    io.Writer
}

func (a *App) Run() error {
	processes, err := a.Provider.ListProcesses()
	if err != nil {
		return fmt.Errorf("could not get process info: %w", err)
	}
	return a.Formatter.Format(processes, a.Output)
}
