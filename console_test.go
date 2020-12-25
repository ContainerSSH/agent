package main

import (
	"encoding/binary"
	"io"
	"os"
	"strings"
	"testing"
	"time"
)

func TestExec(t *testing.T) {
	stdin, _ := io.Pipe()
	_, stdout := io.Pipe()
	_, stderr := io.Pipe()

	defer func() {
		execData := recover().(execData)
		if "test" != execData.argv0 {
			t.FailNow()
		}
		if "test" != execData.argv[0] {
			t.FailNow()
		}
		if len(execData.argv) != 1 {
			t.FailNow()
		}
		if len(execData.envv) != len(os.Environ()) {
			t.FailNow()
		}
	}()

	console(stdin, stdout, stderr, []string{"--", "test"}, exec, exit)
	t.FailNow()
}

func TestExecWithParams(t *testing.T) {
	stdin, _ := io.Pipe()
	_, stdout := io.Pipe()
	_, stderr := io.Pipe()

	defer func() {
		execData := recover().(execData)
		if "test" != execData.argv0 {
			t.FailNow()
		}
		if len(execData.argv) != 3 {
			t.FailNow()
		}
		if execData.argv[0] != "test" {
			t.FailNow()
		}
		if execData.argv[1] != "test1" {
			t.FailNow()
		}
		if execData.argv[2] != "test2" {
			t.FailNow()
		}
		if len(execData.envv) != len(os.Environ()) {
			t.FailNow()
		}
	}()

	console(stdin, stdout, stderr, []string{"--", "test", "test1", "test2"}, exec, exit)
	t.FailNow()
}

func TestExecWithEnv(t *testing.T) {
	stdin, _ := io.Pipe()
	_, stdout := io.Pipe()
	_, stderr := io.Pipe()

	defer func() {
		execData := recover().(execData)
		if "test" != execData.argv0 {
			t.FailNow()
		}
		if len(execData.argv) != 1 {
			t.FailNow()
		}
		for _, env := range execData.envv {
			parts := strings.SplitN(env, "=", 2)
			if parts[0] == "FOO" && parts[1] == "bar" {
				return
			}
		}
		t.FailNow()
	}()

	console(stdin, stdout, stderr, []string{"--env", "FOO=bar", "--", "test"}, exec, exit)
	t.FailNow()
}

func TestExecWait(t *testing.T) {
	stdin, stdinWriter := io.Pipe()
	_, stdout := io.Pipe()
	_, stderr := io.Pipe()

	exitSignal := make(chan struct{}, 1)

	defer func() {
		execData := recover().(execData)
		if "test" != execData.argv0 {
			t.FailNow()
		}
		if len(execData.argv) != 1 {
			t.FailNow()
		}
		if len(execData.envv) != len(os.Environ()) {
			t.FailNow()
		}
		exitSignal <- struct{}{}
	}()

	go func() {
		select {
		case <-exitSignal:
			t.Fail()
		case <-time.After(500 * time.Millisecond):
		}
		if _, err := stdinWriter.Write([]byte{'\000'}); err != nil {
			t.Fail()
		}
		<-exitSignal
	}()

	console(stdin, stdout, stderr, []string{"--wait", "--", "test"}, exec, exit)
	t.FailNow()
}

func TestExecPid(t *testing.T) {
	stdin, _ := io.Pipe()
	stdoutReader, stdout := io.Pipe()
	_, stderr := io.Pipe()

	defer func() {
		execData := recover().(execData)
		if "test" != execData.argv0 {
			t.FailNow()
		}
		if len(execData.argv) != 1 {
			t.FailNow()
		}
		if len(execData.envv) != len(os.Environ()) {
			t.FailNow()
		}
	}()

	go func() {
		data := make([]byte, 4)
		bytesRead, err := stdoutReader.Read(data)
		if err != nil {
			t.Fail()
		}
		if bytesRead != 4 {
			t.Fail()
		}
		pid := binary.LittleEndian.Uint32(data)
		if pid != uint32(os.Getpid()) {
			t.Fail()
		}
	}()

	console(stdin, stdout, stderr, []string{"--pid", "--", "test"}, exec, exit)
	t.FailNow()
}
