package belfry

import (
	"fmt"
	"io"
	"net"
	"strings"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

// RunArgs is the arguments for Run.
type RunArgs struct {
	Stdin   io.Reader
	Stdout  io.Writer
	Stderr  io.Writer
	Bastion string
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

	addr := net.JoinHostPort(args.Host, args.Port)
	if args.Bastion != "" {
		addr = net.JoinHostPort(args.Bastion, "22")
	}

	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return fmt.Errorf("Failed to connect: %w", err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("Failed to create session: %w", err)
	}
	defer session.Close()

	cmd := strings.Join(args.Command, " ")

	if args.Bastion != "" {
		key, err := GetPrivateKey()
		if err != nil {
			return err
		}

		keyring := agent.NewKeyring()
		err = keyring.Add(agent.AddedKey{
			PrivateKey: key,
		})
		if err != nil {
			return fmt.Errorf("Failed to add key: %w", err)
		}

		err = agent.ForwardToAgent(client, keyring)
		if err != nil {
			return fmt.Errorf("Failed to forward to SSH-agent: %w", err)
		}

		err = agent.RequestAgentForwarding(session)
		if err != nil {
			return fmt.Errorf("Failed to request SSH-agent forward: %w", err)
		}

		cmd = fmt.Sprintf("ssh -p %s %s %s", args.Port, args.Host, cmd)
	}

	session.Stdin = args.Stdin
	session.Stdout = args.Stdout
	session.Stderr = args.Stderr
	err = session.Run(cmd)
	if err != nil {
		return fmt.Errorf("Failed to run: %w", err)
	}

	return nil
}
