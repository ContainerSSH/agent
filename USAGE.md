# Using the ContainerSSH agent

This is the agent meant to be installed in the containers run by
ContainerSSH. While images work without this agent, several features
may not be available depending on what the container engine supports.

## Usage

    ./agent MODE [OPTIONS] 

Currently, the agent has the following modes:

`console`
: This mode will take care of console-related activities.

`signal`
: This mode will send a signal to a process.

---

## Console mode

The console mode takes care of console-related settings that are not
supported by the container backend. The console mode has the following
required argument:

    ./agent MODE [OPTIONS] -- PROGRAM

The program will be launched with `execve` and replace the agent
completely. (The agent will not be running once the program is launched.)
This makes the agent safe to run even when you need the program to be run
as PID 1.

### Setting environment variables

The console mode accepts environment variables because some container
backends, such as Kubernetes on `exec`, do not support setting environment
variables.

Environment variables can be set using the following argument form:

    --env NAME=VALUE

## Waiting

Kubernetes does not wait for a console to be attached, therefore a first
line of a shell or equivalent will be lost. Therefore, the agent has a
`--wait` flag.

The `--wait` flag will not launch the desired program until a byte with
the value of `\0` is received on `stdin`. This byte will not be passed to
the program and if any other byte is received the agent will exit to
prevent data stream corruption.

## Reporting PID

Since most container implementations don't report the PID of the process
inside the container in the `exec` method the `--pid` flag will report
the PID of the current program over stdout in the first 4 bytes as a 
little-endian uint32.

## Return codes

This program will exit with one of the following exit codes:

`1`
: General configuration error. See `stderr` for details.

`2`
: Could not read from `stdin` with `--wait`.

`3`
: Could not write PID to `stdout` with `--pid`.

`127`
: Could not execute desired program.

---

## Signal mode

The signal mode is useful for forwarding signals to processes that are
already running.

The signal mode has the following syntax:

    ./agent signal --pid 1234 --signal ABRT

The signal names are, as per RFC 4254 Section 6.9:

- `ABRT`
- `ALRM`
- `FPE`
- `HUP`
- `ILL`
- `INT`
- `KILL`
- `PIPE`
- `QUIT`
- `SEGV`
- `TERM`
- `USR1`
- `USR2`

**Note:** `USR1` and `USR2` are not supported on Windows.

This mode has the following exit codes:

`0`
: Signal successfully sent.

`1`
: General configuration error. See `stderr` for details.

`2`
: The specified process was not found. See `stderr` for details.

`3`
: Failed to send signal to process. See `stderr` for details.
