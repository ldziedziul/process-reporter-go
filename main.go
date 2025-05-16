package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	format := flag.String("format", "json", "Output format: json or csv")
	output := flag.String("output", "", "Optional output file name. If not set, output goes to stdout")
	flag.Parse()

	formatter, err := getFormatter(*format)
	if err != nil {
		log.Fatalf("Unsupported format: %v", err)
	}

	provider := &SystemProcessProvider{}
	var out *os.File
	if *output == "" {
		out = os.Stdout
	} else {
		out, err = os.Create(*output)
		if err != nil {
			log.Fatalf("Failed to create output file: %v", err)
		}
		defer func(out *os.File) {
			err := out.Close()
			if err != nil {
				log.Fatalf("Failed to close output file: %v", err)
			}
		}(out)
	}

	app := App{
		Provider:  provider,
		Formatter: formatter,
		Output:    out,
	}

	if err := app.Run(); err != nil {
		log.Fatalf("Failed to run application: %v", err)
	}

	if *output != "" {
		fmt.Printf("Report saved successfully in %s format to %s\n", *format, *output)
	}
}
