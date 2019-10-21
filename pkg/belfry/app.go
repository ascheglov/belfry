package belfry

import (
	"fmt"
	"io"
	"net"
	"strings"

	"golang.org/x/crypto/ssh"
)

// RunArgs is the arguments for Run.
type RunArgs struct {
	Stdin   io.Reader
	Stdout  io.Writer
	Stderr  io.Writer
	Host    string
	Port    string
	Command []string
}

// Run is the main entry point.
func Run(args *RunArgs) error {
	config, err := DefaultSSHConfig()
	if err != nil {
		return err
	}

	client, err := ssh.Dial("tcp", net.JoinHostPort(args.Host, args.Port), config)
	if err != nil {
		return fmt.Errorf("Failed to connect: %w", err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("Failed to create session: %w", err)
	}
	defer session.Close()

	session.Stdin = args.Stdin
	session.Stdout = args.Stdout
	session.Stderr = args.Stderr
	err = session.Run(strings.Join(args.Command, " "))
	if err != nil {
		return fmt.Errorf("Failed to run: %w", err)
	}

	return nil
}
