package main

import (
	"io"
	"os"
	"log"
)

func writeFile(args []string, stdin io.Reader, stderr io.Writer, exit exitFunc) {
	l := log.New(stderr, "", 0)
	if len(args) < 1 {
		exit(1)
	}

	file, err := os.OpenFile(args[0], os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		l.Printf("Error opening file: %+v", err)
		exit(1)
	}

	err = file.Chmod(0600)
	if err != nil {
		l.Printf("Error chmod'ing file: %+v", err)
		exit(1)
	}

	buf := make([]byte, 8192)
	for {
		nBytes, err := stdin.Read(buf)
		if err != nil {
			exit(0)
		}

		wBytes, err := file.Write(buf[0:nBytes])
		if err != nil {
			exit(1)
		}
		if wBytes != nBytes {
			// Write should never fail without error
			exit(1)
		}
	}
}
