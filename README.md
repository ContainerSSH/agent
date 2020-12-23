[![ContainerSSH - Launch Containers on Demand](https://containerssh.github.io/images/logo-for-embedding.svg)](https://containerssh.io/)

<!--suppress HtmlDeprecatedAttribute -->
<h1 align="center">ContainerSSH Guest Agent</h1>

[![Go Report Card](https://goreportcard.com/badge/github.com/containerssh/agent?style=for-the-badge)](https://goreportcard.com/report/github.com/containerssh/agent)
[![LGTM Alerts](https://img.shields.io/lgtm/alerts/github/ContainerSSH/agent?style=for-the-badge)](https://lgtm.com/projects/g/ContainerSSH/agent/)

This is the agent meant to be installed in the containers run by ContainerSSH. While images work without this agent, several features may not be available depending on what the container engine supports. 

<p align="center"><strong>Note: This is a developer documentation.</strong><br />The user documentation for ContainerSSH is located at <a href="https://containerssh.github.io">containerssh.io</a>.</p>

## How this application works

This application is intended as a single binary to be embedded into a container image as a wrapper for the actual program to be launched. Currently, this program supports one mode that can be invoked as follows:

```bash
./agent console --env FOO=bar --env TERM=xterm --wait --pid -- /bin/bash
```

The parameters are as follows:

- `console` sets the agent to console mode. (This is the only mode supported.)
- `env` passes an environment variable to the desired program.
- `wait` waits for a `\0` byte on the `stdin` before launching the desired program.
- `pid` writes the process ID of the program to the `stdout` in the first 4 bytes as a little-endian `uint32` before launching the program.

The detailed usage is documented in [USAGE.md](USAGE.md).

## Building this application

This application can be built by running the following two programs:

- `go generate`
- `go build `

