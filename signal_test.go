package main

import (
	"fmt"
	"io"
	"os"
	"testing"
)

type testProcess struct {
	receivedSignals []os.Signal
	reject          bool
}

func (t *testProcess) Signal(sig os.Signal) error {
	if t.reject {
		return fmt.Errorf("rejected")
	}
	t.receivedSignals = append(t.receivedSignals, sig)
	return nil
}

func TestSendSignal(t *testing.T) {
	defer func() {
		data := recover().(exitData)
		if data.code != 0 {
			t.Fail()
		}
	}()
	_, stderr := io.Pipe()
	args := []string{"--pid", "1", "--signal", "TERM"}
	testProcess := &testProcess{}
	signal(stderr, args, exit, func(pid int) (Process, error) {
		if pid != 1 {
			return nil, fmt.Errorf("invalid PID")
		}
		return testProcess, nil
	})
	if len(testProcess.receivedSignals) != 1 {
		t.FailNow()
	}
	if testProcess.receivedSignals[0].String() != "TERM" {
		t.FailNow()
	}
}
