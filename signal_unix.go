// +build linux freebsd openbsd darwin

package main

import (
	"fmt"
	"os"
	"syscall"
)

// processSignalName creates an os.Signal from the specified name from RFC 4254 Section 6.9
// ( https://tools.ietf.org/html/rfc4254#section-6.9 ) or returns an error if the name is not supported on the
// platform.
func processSignalName(name string) (os.Signal, error) {
	switch name {
	case "ABRT":
		return syscall.SIGABRT, nil
	case "ALRM":
		return syscall.SIGALRM, nil
	case "FPE":
		return syscall.SIGFPE, nil
	case "HUP":
		return syscall.SIGHUP, nil
	case "ILL":
		return syscall.SIGILL, nil
	case "INT":
		return syscall.SIGINT, nil
	case "KILL":
		return syscall.SIGKILL, nil
	case "PIPE":
		return syscall.SIGPIPE, nil
	case "QUIT":
		return syscall.SIGQUIT, nil
	case "SEGV":
		return syscall.SIGSEGV, nil
	case "TERM":
		return syscall.SIGTERM, nil
	case "USR1":
		return syscall.SIGUSR1, nil
	case "USR2":
		return syscall.SIGUSR2, nil
	default:
		return nil, fmt.Errorf("unsupported signal: %s", name)
	}
}
