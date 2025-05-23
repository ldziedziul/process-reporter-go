package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

func run(args []string, sysOut io.Writer, sysErrOut io.Writer) error {
	fs := flag.NewFlagSet("process-reporter", flag.ContinueOnError)
	fs.SetOutput(sysErrOut) // Flag errors go to errOut

	format := fs.String("format", "json", "Output format: json or csv")
	output := fs.String("output", "", "Optional output file name. If not set, output goes to stdout")

	if err := fs.Parse(args); err != nil {
		return err
	}

	formatter, err := getFormatter(*format)
	if err != nil {
		return fmt.Errorf("unsupported format: %v", err)
	}

	provider := &SystemProcessProvider{}

	var out = sysOut
	var file *os.File
	if *output != "" {
		file, err = os.Create(*output)
		if err != nil {
			return fmt.Errorf("failed to create output file: %v", err)
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				fmt.Println(fmt.Errorf("error closing file: %w", err))
			}
		}(file)
		out = file
	}

	app := App{
		Provider:  provider,
		Formatter: formatter,
		Output:    out,
	}

	if err := app.Run(); err != nil {
		return fmt.Errorf("failed to run application: %v", err)
	}

	if *output != "" {
		_, err := fmt.Fprintf(sysOut, "Report saved successfully in %s format to %s\n", *format, *output)
		if err != nil {
			return fmt.Errorf("failed to write to output: %w", err)
		}
	}

	return nil
}

func main() {
	if err := run(os.Args[1:], os.Stdout, os.Stderr); err != nil {
		_, err := fmt.Fprintln(os.Stderr, err)
		if err != nil {
			fmt.Println(fmt.Errorf("error writing to stderr: %w", err))
		}
		os.Exit(1)
	}
}
