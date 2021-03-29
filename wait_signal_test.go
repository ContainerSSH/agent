// +build linux freebsd openbsd darwin

package main

import (
	"bytes"
	"os"
	"syscall"
	"testing"
)

func TestWaitSignal(t *testing.T) {
	stdout := &bytes.Buffer{}
	exit := func(code int) {
		panic(code)
	}
	done := make(chan struct{})
	exitCode := -1
	exitMessage := "Exited with INT."
	go func() {
		defer func() {
			exitCode = recover().(int)
			done <- struct{}{}
		}()
		waitSignal(stdout, []string{
			"--signal", "INT",
			"--message", exitMessage,
		}, exit)
	}()
	p, err := os.FindProcess(os.Getpid())
	if err != nil {
		t.Fatal(err)
	}
	if err := p.Signal(syscall.SIGINT); err != nil {
		t.Fatal(err)
	}
	<-done
	if exitCode != 0 {
		t.Fatalf("Invalid exit code: %d", exitCode)
	}
	output := stdout.String()
	if output != "Exited with INT." {
		t.Fatalf("Invalid exit message: %s", output)
	}
}
