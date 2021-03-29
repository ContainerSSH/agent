package main

import (
	"fmt"
	"io"
	"os"
	osSignal "os/signal"
)

func waitSignal(stderr io.Writer, args []string, exit exitFunc) {
	var sig []os.Signal
	var message = ""
	var err error

loop:
	for {
		if len(args) == 0 {
			break loop
		}
		switch args[0] {
		case "--message":
			if len(args[0]) < 1 {
				usage("--message requires an argument", true)
			}
			message = args[1]
			args = args[2:]
		case "--signal":
			if len(args[0]) < 1 {
				usage("--signal requires an argument", true)
			}
			signalName := args[1]
			var signalObject os.Signal
			args = args[2:]

			signalObject, err = processSignalName(signalName)
			if err != nil {
				_, _ = stderr.Write([]byte(err.Error()))
				exit(1)
			}
			sig = append(sig, signalObject)
		default:
			usage(fmt.Sprint("unknown option: #{args[0]}"), true)
		}
	}
	if len(message) == 0 {
		usage("missing required --message option", true)
	}
	if sig == nil {
		usage("missing required --signal option", true)
	}

	chanSig := make(chan os.Signal, 1)
	chanDone := make(chan bool, 1)

	for _, signalSubscribe := range sig {
		osSignal.Notify(chanSig, signalSubscribe)
	}

	go func() {
		<-chanSig
		chanDone <- true
	}()
	<-chanDone

	fmt.Printf(message)
	exit(0)
}
