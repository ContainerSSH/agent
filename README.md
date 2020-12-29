[![ContainerSSH - Launch Containers on Demand](https://containerssh.github.io/images/logo-for-embedding.svg)](https://containerssh.io/)

<!--suppress HtmlDeprecatedAttribute -->
<h1 align="center">ContainerSSH Guest Agent</h1>

[![Go Report Card](https://goreportcard.com/badge/github.com/containerssh/agent?style=for-the-badge)](https://goreportcard.com/report/github.com/containerssh/agent)
[![LGTM Alerts](https://img.shields.io/lgtm/alerts/github/ContainerSSH/agent?style=for-the-badge)](https://lgtm.com/projects/g/ContainerSSH/agent/)

This is the agent meant to be installed in the containers run by ContainerSSH. While images work without this agent, several features may not be available depending on what the container engine supports. 

<p align="center"><strong>Note: This is a developer documentation.</strong><br />The user documentation for ContainerSSH is located at <a href="https://containerssh.github.io">containerssh.io</a>.</p>

## Integrating the agent

This agent is intended to be integrated into container images. There are two main installation methods:

### Using the base image (recommended)

This method uses the `containerssh/agent` container image as part of a multistage build:

```
FROM containerssh/agent AS agent

FROM your-base-image
COPY --from=agent /usr/bin/containerssh-agent /usr/bin/containerssh-agent
# Your other build commands here
```

### Installing the binaries

To use this method go to the [latest release from the releases section](https://github.com/ContainerSSH/agent/releases) and verify it against our [https://containerssh.io/gpg.txt](https://containerssh.io/gpg.txt) key (`3EE5B012FA7B400CD952601E4689F1F0F358FABA`).

On an Ubuntu image build this would involve the following steps:

```Dockerfile
ARG AGENT_GPG_FINGERPRINT=3EE5B012FA7B400CD952601E4689F1F0F358FABA
ARG AGENT_GPG_SOURCE=https://containerssh.io/gpg.txt

RUN echo "\e[1;32mInstalling ContainerSSH guest agent...\e[0m" && \
    DEBIAN_FRONTEND=noninteractive apt-get -o Dpkg::Options::='--force-confold' update && \
    DEBIAN_FRONTEND=noninteractive apt-get -o Dpkg::Options::='--force-confold' -fuy --allow-downgrades --allow-remove-essential --allow-change-held-packages install gpg && \
    wget -q -O - https://api.github.com/repos/containerssh/agent/releases/latest | grep browser_download_url | grep -e "agent_.*_linux_amd64.deb" | awk ' { print $2 } ' | sed -e 's/"//g' > /tmp/assets.txt && \
    wget -q -O /tmp/agent.deb $(cat /tmp/assets.txt |grep -v .sig) && \
    wget -q -O /tmp/agent.deb.sig $(cat /tmp/assets.txt |grep .sig) && \
    wget -q -O - $AGENT_GPG_SOURCE | gpg --import && \
    echo -e "5\ny\n" | gpg --command-fd 0 --batch --expert --edit-key $AGENT_GPG_FINGERPRINT trust && \
    test $(gpg --status-fd=1 --verify /tmp/agent.deb.sig /tmp/agent.deb | grep VALIDSIG | grep $AGENT_GPG_FINGERPRINT | wc -l) -eq 1 && \
    dpkg -i /tmp/agent.deb && \
    rm -rf /tmp/* && \
    rm -rf ~/.gnupg && \
    DEBIAN_FRONTEND=noninteractive apt-get -o Dpkg::Options::='--force-confold' -fuy --allow-downgrades --allow-remove-essential --allow-change-held-packages remove gpg && \
    DEBIAN_FRONTEND=noninteractive apt-get -o Dpkg::Options::='--force-confold' -y clean && \
    /usr/bin/containerssh-agent -h
```

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

