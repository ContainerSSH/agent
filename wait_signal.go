package main

import (
	"fmt"
	"io"
	"os"
	osSignal "os/signal"
)

func waitSignal(stdout io.Writer, args []string, exit exitFunc) {
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
				usage("--message is duplicate", true, 1, exit)
			}
			message = args[1]
			args = args[2:]
		case "--signal":
			if len(args) < 2 {
				usage("--signal requires an argument", true, 1, exit)
			}
			signalName := args[1]
			var signalObject os.Signal
			args = args[2:]

			signalObject, err = processSignalName(signalName)
			if err != nil {
				usage(err.Error(), true, 1, exit)
			}
			sig = append(sig, signalObject)
		default:
			usage(fmt.Sprintf("unknown option: %s", args[0]), true, 1, exit)
		}
	}
	if len(sig) == 0 {
		usage("missing required --signal option", true, 1, exit)
	}

	chanSig := make(chan os.Signal, 1)
	chanDone := make(chan struct{})

	osSignal.Notify(chanSig, sig...)

	go func() {
		<-chanSig
		chanDone <- struct{}{}
	}()
	<-chanDone

	if len(message) > 0 {
		_, _ = stdout.Write([]byte(message))
	}
	exit(0)
}
