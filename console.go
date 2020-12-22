package main

import (
	"encoding/binary"
	"fmt"
	"os"
	"strings"
	"syscall"
)

func console(args []string) {
	env, program, wait, reportPID := parseConsoleArgs(args)

	if wait {
		firstByte := make([]byte, 1)
		readBytes, err := os.Stdin.Read(firstByte)
		if err != nil {
			os.Exit(2)
		}
		if readBytes != 1 {
			os.Exit(2)
		}
		if firstByte[0] != '\000' {
			usage("non-null first character", true)
		}
	}

	if reportPID {
		pid := uint32(os.Getpid())
		b := make([]byte, 4)
		binary.LittleEndian.PutUint32(b, pid)
		_, err := os.Stdout.Write(b)
		if err != nil {
			os.Exit(3)
		}
	}

	if err := syscall.Exec(program[0], program[1:], env); err != nil {
		_, _ = os.Stderr.Write([]byte(err.Error()))
		os.Exit(127)
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
		case "--reportPID":
			reportPID = true
			args = args[1:]
		default:
			usage(fmt.Sprintf("invalid parameter: %s", args[0]), true)
		}
	}
	return env, program, wait, reportPID
}
