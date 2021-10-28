package main

import (
	"fmt"
	"os"
	"syscall"
)

//go:generate go run scripts/generate-usage.go
//go:generate go run scripts/generate-license.go

func usage(error string, exit bool, exitCode int, exitFunc func(code int)) {
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
		exitFunc(exitCode)
	}
}

func license(argv []string) {
	if len(argv) > 0 {
		usage("license takes no arguments", true, 1, os.Exit)
	}
	_, _ = os.Stdout.Write([]byte(licenseText))
}

func main() {
	args := os.Args

	if len(args) < 2 {
		usage("not enough parameters", true, 1, os.Exit)
	}

	switch args[1] {
	case "-h":
		fallthrough
	case "--help":
		usage("", false, 0, os.Exit)
		os.Exit(0)
	case "console":
		console(os.Stdin, os.Stdout, os.Stderr, args[2:], syscall.Exec, os.Exit)
	case "signal":
		signal(os.Stderr, args[2:], syscall.Exit, func(pid int) (Process, error) {
			return os.FindProcess(pid)
		})
	case "wait-signal":
		waitSignal(os.Stdout, args[2:], os.Exit)
	case "write-file":
		writeFile(args[2:], os.Stdin, os.Stderr, os.Exit)
	case "license":
		license(args[2:])
	default:
		usage(fmt.Sprintf("invalid mode: %s", args[1]), true, 1, os.Exit)
	}
}
