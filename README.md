[![ContainerSSH - Launch Containers on Demand](https://containerssh.github.io/images/logo-for-embedding.svg)](https://containerssh.io/)

<!--suppress HtmlDeprecatedAttribute -->
<h1 align="center">ContainerSSH Guest Agent</h1>

[![Go Report Card](https://goreportcard.com/badge/github.com/containerssh/agent?style=for-the-badge)](https://goreportcard.com/report/github.com/containerssh/agent)
[![LGTM Alerts](https://img.shields.io/lgtm/alerts/github/ContainerSSH/agent?style=for-the-badge)](https://lgtm.com/projects/g/ContainerSSH/agent/)

This is the agent meant to be installed in the containers run by ContainerSSH. While images work without this agent, several features may not be available depending on what the container engine supports. 

<p align="center"><strong>Note: This is a developer documentation.</strong><br />The user documentation for ContainerSSH is located at <a href="https://containerssh.github.io">containerssh.io</a>.</p>

## Integrating the agent

This agent is intended to be integrated into container images. The installation process is as follows:

1. Go to the [releases section](https://github.com/ContainerSSH/agent/releases).
2. Download the latest release for your platform, as well as the signature file.
3. Download the GPG key from [https://containerssh.io/gpg.txt](https://containerssh.io/gpg.txt).
4. Import the GPG key into your keychain (`gpg --import gpg.txt`).
5. Edit the key: `gpg --edit-key 3EE5B012FA7B400CD952601E4689F1F0F358FABA`
6. Mark it as trusted: `trust`
7. Quit using the `quit` command.
8. Verify the GPG signature: `gpg --verify downloadedfile.sig downloadedfile`
9. Install the file using your package manager.

You can look at the default [guest image Dockerfile](https://github.com/containerssh/guest-image) for an example on Ubuntu.

## How this application works

This application is intended as a single binary to be embedded into a container image to handle features that the container engine (Docker, Kubernetes) does not support. Currently, the following modes are supported:

```bash
# Run in Console mode
./agent console --env FOO=bar --env TERM=xterm --wait --pid -- /bin/bash

# Run in Signal mode
./agent signal --pid 3 --signal TERM
```

The console mode supports the following parameters:

- `console` sets the agent to console mode.
- `env` passes an environment variable to the desired program.
- `wait` waits for a `\0` byte on the `stdin` before launching the desired program.
- `pid` writes the process ID of the program to the `stdout` in the first 4 bytes as a little-endian `uint32` before launching the program.

The signal mode supports the following parameters:

- `signal` sets the agent to signal mode.
- `pid` passes the process ID to send the signal to.
- `signal` sets the signal to send. These are defined in [RFC 4254 Section 6.9](https://tools.ietf.org/html/rfc4254#section-6.9).

The detailed usage is documented in [USAGE.md](USAGE.md).

## Building this application

This application can be built by running the following two programs:

- `go generate`
- `go build `

