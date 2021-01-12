package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strings"
)

// execFunc is a function that replaces the current program image with the executed one. If only returns with an error
// if the execution failed
type execFunc func(argv0 string, argv []string, envv []string) (err error)

// exitFunc exits the current program with the specified exit code.
type exitFunc func(code int)

func console(stdin io.Reader, stdout io.Writer, stderr io.Writer, args []string, exec execFunc, exit exitFunc) {
	env, program, wait, reportPID := parseConsoleArgs(args)

	if wait {
		firstByte := make([]byte, 1)
		readBytes, err := stdin.Read(firstByte)
		if err != nil {
			usage(fmt.Sprintf("Failed to read from stdin (%v)", err), false)
			exit(2)
		}
		if readBytes != 1 {
			usage("Nothing read from stdin", false)
			exit(2)
		}
		if firstByte[0] != '\000' && firstByte[0] != '\n' {
			usage("Unexpected first character. Send either a null or a newline character to start program.", false)
			exit(2)
		}
	}

	if reportPID {
		pid := uint32(os.Getpid())
		b := make([]byte, 4)
		binary.LittleEndian.PutUint32(b, pid)
		_, err := stdout.Write(b)
		if err != nil {
			exit(3)
		}
	}

	if err := exec(program[0], program, env); err != nil {
		_, _ = stderr.Write([]byte(err.Error()))
		exit(127)
	}
}

func parseConsoleArgs(args []string) ([]string, []string, bool, bool) {
	var env = os.Environ()
	var program []string
	wait := false
	reportPID := false

loop:
	for {
		if len(args) == 0 {
			break loop
		}
		switch args[0] {
		case "--":
			if len(args) == 1 {
				usage("no program to run passed", true)

			}
			program = args[1:]
			args = []string{}
		case "--env":
			if len(args) == 1 {
				usage("--env requires an argument", true)
			}
			if strings.HasPrefix(args[1], "--") {
				usage("--env requires an argument", true)
			}
			if !strings.Contains(args[1], "=") {
				usage("--env requires an argument in the format of NAME=VALUE", true)
			}
			env = append(env, args[1])
			args = args[2:]
		case "--wait":
			wait = true
			args = args[1:]
		case "--pid":
			reportPID = true
			args = args[1:]
		default:
			usage(fmt.Sprintf("invalid parameter: %s", args[0]), true)
		}
	}
	return env, program, wait, reportPID
}
