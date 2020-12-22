# Using the ContainerSSH agent

This is the agent meant to be installed in the containers run by
ContainerSSH. While images work without this agent, several features
may not be available depending on what the container engine supports.

## Usage

    ./agent MODE [OPTIONS] 

Currently, the agent has one mode:

`console`
: This mode will take care of console-related activities.

---

## Console mode

The console mode takes care of console-related settings that are not
supported by the container backend. The console mode has the following
required argument:

    ./agent MODE [OPTIONS] -- PROGRAM

The program will be launched with `execv`, replacing the agent image
with the launched program as PID 1.

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
the program and if any other byte is received the agent will exit.

