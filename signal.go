package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
)

// Process is an abstraction for process signal handling.
type Process interface {
	// Signal sends the specified sig to the process and returns an error if the signal failed.
	Signal(sig os.Signal) error
}

// signal sends the process indicated with --pid the signal indicated with --signal. The signal names are described in
// RFC 4254 Section 6.9 ( https://tools.ietf.org/html/rfc4254#section-6.9 )
func signal(stderr io.Writer, args []string, exit exitFunc, findProcess func(pid int) (Process, error)) {
	var sig os.Signal = nil
	pid := -1
	var err error

loop:
	for {
		if len(args) == 0 {
			break loop
		}
		switch args[0] {
		case "--pid":
			if len(args[0]) < 1 {
				usage("--pid requires an argument", true)
			}
			if pid > -1 {
				usage("--pid is duplicate", true)
			}
			pid, err = strconv.Atoi(args[1])
			if err != nil || pid < 1 {
				usage("--pid requires a positive integer argument", true)
			}
			args = args[2:]
		case "--signal":
			if len(args[0]) < 1 {
				usage("--signal requires an argument", true)
			}
			signalName := args[1]
			args = args[2:]

			sig, err = processSignalName(signalName)
			if err != nil {
				_, _ = stderr.Write([]byte(err.Error()))
				exit(1)
			}
		default:
			usage(fmt.Sprintf("unknown option: %s", args[0]), true)
		}
	}
	if pid < 0 {
		usage("missing required --pid option", true)
	}
	if sig == nil {
		usage("missing required --signal option", true)
	}
	process, err := findProcess(pid)
	if err != nil {
		_, _ = stderr.Write([]byte(fmt.Sprintf("process not found (%v)", err)))
		exit(2)
	}
	if err := process.Signal(sig); err != nil {
		_, _ = stderr.Write([]byte(fmt.Sprintf("failed to send signal (%v)", err)))
		exit(3)
	}
	exit(0)
}
