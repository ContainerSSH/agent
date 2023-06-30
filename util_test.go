package main

// exitData contains the code from the exit call.
type exitData struct {
	code int
}

// exit panics with the exitData structure to emulate exiting the process.
func exit(code int) {
	panic(exitData{
		code: code,
	})
}

// execData contains the called parameters on exec.
type execData struct {
	argv0 string
	argv  []string
	envv  []string
}

// exec panics with the execData structure to emulate executing a program and replacing the current process image.
func fakeExec(argv0 string, argv []string, envv []string) (err error) {
	panic(execData{
		argv0: argv0,
		argv:  argv,
		envv:  envv,
	})
}
