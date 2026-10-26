FROM scratch
ARG TARGETPLATFORM
COPY $TARGETPLATFORM/containerssh-agent /usr/bin/containerssh-agent
