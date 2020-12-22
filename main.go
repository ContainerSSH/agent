package main

import (
	"fmt"
	"os"
)

//go:generate go run scripts/generate-usage.go

func usage(error string, exit bool) {
	if fileInfo, _ := os.Stdout.Stat(); (fileInfo.Mode() & os.ModeCharDevice) != 0 {
		if error != "" {
			_, _ = os.Stderr.Write([]byte("\033[31mError: " + error + "\033[0m\n"))
		}
		_, _ = os.Stdout.Write([]byte("\n"))
		_, _ = os.Stdout.Write([]byte(usageTextColored))
	} else {
		if error != "" {
			_, _ = os.Stderr.Write([]byte("Error: " + error + "\n"))
		}
		_, _ = os.Stdout.Write([]byte("\n"))
		_, _ = os.Stdout.Write([]byte(usageText))
	}
	if exit {
		os.Exit(1)
	}
}

func main() {
	args := os.Args

	if len(args) < 2 {
		usage("not enough parameters", true)
	}

	switch args[1] {
	case "console":
		console(args[2:])
	default:
		usage(fmt.Sprintf("invalid mode: %s", args[1]), false)
	}
}
