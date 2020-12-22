package main

import (
	"encoding/binary"
	"io"
	"os"
	"strings"
	"testing"
	"time"
)

type exitData struct {
	code int
}

func exit(code int) {
	panic(exitData{
		code: code,
	})
}

type execData struct {
	argv0 string
	argv  []string
	envv  []string
}

func exec(argv0 string, argv []string, envv []string) (err error) {
	panic(execData{
		argv0: argv0,
		argv:  argv,
		envv:  envv,
	})
}

func TestExec(t *testing.T) {
	stdin, _ := io.Pipe()
	_, stdout := io.Pipe()
	_, stderr := io.Pipe()

	defer func() {
		execData := recover().(execData)
		if "test" != execData.argv0 {
			t.FailNow()
		}
		if len(execData.argv) != 0 {
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
		if len(execData.argv) != 2 {
			t.FailNow()
		}
		if execData.argv[0] != "test1" {
			t.FailNow()
		}
		if execData.argv[1] != "test2" {
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
		if len(execData.argv) != 0 {
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
		if len(execData.argv) != 0 {
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
			t.FailNow()
		case <-time.After(500 * time.Millisecond):
		}
		if _, err := stdinWriter.Write([]byte{'\000'}); err != nil {
			t.FailNow()
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
		if len(execData.argv) != 0 {
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
			t.FailNow()
		}
		if bytesRead != 4 {
			t.FailNow()
		}
		pid := binary.LittleEndian.Uint32(data)
		if pid != uint32(os.Getpid()) {
			t.FailNow()
		}
	}()

	console(stdin, stdout, stderr, []string{"--pid", "--", "test"}, exec, exit)
	t.FailNow()
}
