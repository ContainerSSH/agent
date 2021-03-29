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
			if len(message) > 0 {
				usage("--message is duplicate", true)
			}
			message = args[1]
			args = args[2:]
		case "--signal":
			if len(args) < 2 {
				usage("--signal requires an argument", true)
			}
			signalName := args[1]
			var signalObject os.Signal
			args = args[2:]

			signalObject, err = processSignalName(signalName)
			if err != nil {
				usage(err.Error(), true)
			}
			sig = append(sig, signalObject)
		default:
			usage(fmt.Sprintf("unknown option: %s", args[0]), true)
		}
	}
	if len(sig) == 0 {
		usage("missing required --signal option", true)
	}

	chanSig := make(chan os.Signal, 1)
	chanDone := make(chan struct{})

	for _, signalSubscribe := range sig {
		osSignal.Notify(chanSig, signalSubscribe)
	}

	go func() {
		<-chanSig
		chanDone <- struct{}{}
	}()
	<-chanDone

	if len(message) > 0 {
		fmt.Print(message)
	}
	exit(0)
}
