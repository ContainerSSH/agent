package main

import (
	"log"
	"os"
	"strings"
	"syscall"
)

func console(args []string) {
	var env []string
	var program []string
	wait := false

loop:
	for {
		if len(args) == 0 {
			break loop
		}
		nextArg := args[0]
		switch nextArg {
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
		}
	}

	if wait {
		firstByte := make([]byte, 1)
		readBytes, err := os.Stdin.Read(firstByte)
		if err != nil {
			os.Exit(2)
		}
		if readBytes != 1 {
			os.Exit(3)
		}
		if firstByte[0] != '\000' {
			usage("non-null first character", true)
		}
	}

	if err := syscall.Exec(program[0], program[1:], env); err != nil {
		log.Fatal(err)
	}
}
