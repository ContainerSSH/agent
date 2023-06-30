package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"strconv"

	proto "go.containerssh.io/libcontainerssh/agentprotocol"
	config "go.containerssh.io/libcontainerssh/config"
	log "go.containerssh.io/libcontainerssh/log"
)

const (
	xauth_path = "/usr/bin/xauth"
)

func serveConnection(log log.Logger, from io.ReadWriteCloser, to io.ReadWriteCloser) {
	_, err := io.Copy(from, to)
	if err != nil && errors.Is(err, io.EOF) {
		log.Warning("Connection error", err)
	}
	from.Close()
	to.Close()
}

func checkCreateXAuthority() error {
	xauthority, ok := os.LookupEnv("XAUTHORITY")
	if !ok {
		xauthority = ".Xauthority"
	}
	file, err := os.OpenFile(xauthority, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	return file.Close()
}

func setupX11(log log.Logger, setup proto.SetupPacket) proto.SetupPacket {
	err := checkCreateXAuthority()
	if err != nil {
		log.Error("Failed to create .Xauthority", err)
		panic(err)
	}

	cmd := exec.Command(xauth_path, "add", ":10."+setup.Screen, setup.AuthProtocol, setup.AuthCookie)
	err = cmd.Run()
	if err != nil {
		log.Error("Failed to run xauth", err)
		panic(err)
	}

	setup.BindHost = "127.0.0.1"
	// Magic X11 formula: 6000 + display number
	setup.BindPort = 6010
	return setup
}

func parsePort(proto string, host string, port uint32) string {
	switch proto {
	case "tcp":
		return host + ":" + strconv.Itoa(int(port))
	case "unix":
		return host
	default:
		panic(fmt.Errorf("unknown protocol %s", proto))
	}
}

//nolint:funlen
func localForward(
	log log.Logger,
	forwardCtx *proto.ForwardCtx,
	connChan chan *proto.Connection,
	setup proto.SetupPacket,
) {
	listenAddr := parsePort(setup.Protocol, setup.BindHost, setup.BindPort)

	sock, err := net.Listen(setup.Protocol, listenAddr)
	if err != nil {
		log.Error("Failed to start listening for connections", listenAddr, err)
		os.Exit(1)
	}

	go func() {
		for {
			conn, ok := <-connChan
			if !ok {
				sock.Close()
				break
			}
			_ = conn.Reject()
		}
	}()

	for {
		conn, err := sock.Accept()
		if err != nil {
			log.Warning("Failed to accept connection from os", err)
			break
		}

		var addr string
		var port uint32
		var agentCon io.ReadWriteCloser
		switch setup.Protocol {
		case "tcp":
			addr = conn.RemoteAddr().(*net.TCPAddr).IP.String()
			port = uint32(conn.RemoteAddr().(*net.TCPAddr).Port)
			agentCon, err = forwardCtx.NewConnectionTCP(
				setup.BindHost,
				setup.BindPort,
				addr,
				uint32(port),
				func() error {
					return conn.Close()
				},
			)
		case "unix":
			agentCon, err = forwardCtx.NewConnectionUnix(
				setup.BindHost,
				func() error {
					return conn.Close()
				},
			)
		default:
			panic(fmt.Errorf("unknown protocol %s", setup.Protocol))
		}
		if err != nil {
			log.Warning("Failed to create new connection with backend", err)
		}

		go serveConnection(log, conn, agentCon)
		go serveConnection(log, agentCon, conn)

		if setup.SingleConnection {
			break
		}
	}

	forwardCtx.WaitFinish()
}

func externalDial(log log.Logger, forwardCtx *proto.ForwardCtx, connChan chan *proto.Connection, setup proto.SetupPacket) {
	for {
		agentCon, ok := <-connChan
		if !ok {
			break
		}
		details := agentCon.Details()
		var protocol string
		switch details.Protocol {
		case proto.PROTOCOL_TCP:
			protocol = "tcp"
		case "unix":
			protocol = "unix"
		default:
			panic(fmt.Errorf("unknown protocol %s", details.Protocol))
		}
		log.Warning(fmt.Sprintf("Dialing %s %s:%d", setup.Protocol, details.ConnectedAddress, details.ConnectedPort))

		dialAddr := parsePort(protocol, details.ConnectedAddress, details.ConnectedPort)

		conn, err := net.Dial(protocol, dialAddr)
		if err != nil {
			log.Warning("Failed to dial %s", dialAddr, err)
			_ = agentCon.Reject()
			continue
		}
		err = agentCon.Accept()
		if err != nil {
			log.Warning("Failed to accept connection", err)
			continue
		}
		go serveConnection(log, conn, agentCon)
		go serveConnection(log, agentCon, conn)
	}

	forwardCtx.WaitFinish()
}

func forwardServer(stdin io.Reader, stdout io.Writer, stderr io.Writer, exit exitFunc) {
	logConfig := config.LogConfig{
		Level:       config.LogLevelDebug,
		Destination: config.LogDestinationStdout,
		File:        "/tmp/agent.log",
		Stdout:      stderr,
		Format:      config.LogFormatLJSON,
	}
	log, err := log.NewLogger(logConfig)
	if err != nil {
		panic(err)
	}

	log.Debug("Starting agent")
	forwardCtx := proto.NewForwardCtx(stdin, stdout, log)

	conType, setup, connChan, err := forwardCtx.StartClient()
	if err != nil {
		panic(err)
	}

	switch conType {
	case proto.CONNECTION_TYPE_X11:
		setup = setupX11(log, setup)
		fallthrough
	case proto.CONNECTION_TYPE_SOCKET_FORWARD:
		fallthrough
	case proto.CONNECTION_TYPE_PORT_FORWARD:
		localForward(log, forwardCtx, connChan, setup)
	case proto.CONNECTION_TYPE_SOCKET_DIAL:
		fallthrough
	case proto.CONNECTION_TYPE_PORT_DIAL:
		externalDial(log, forwardCtx, connChan, setup)
	}
}
