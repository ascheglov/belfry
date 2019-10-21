package belfry

import (
	"fmt"
	"io/ioutil"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

// GetPrivateKey reads default private key.
func GetPrivateKey() (interface{}, error) {
	user, err := user.Current()
	if err != nil {
		return nil, fmt.Errorf("Failed to get user info: %w", err)
	}

	pemBytes, err := ioutil.ReadFile(filepath.Join(user.HomeDir, ".ssh/id_rsa"))
	if err != nil {
		return nil, fmt.Errorf("Failed to read private key: %w", err)
	}

	key, err := ssh.ParseRawPrivateKey(pemBytes)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse private key: %w", err)
	}

	return key, nil
}

// GetSSHKey gets ssh key.
func GetSSHKey(keyPath string) (ssh.AuthMethod, error) {
	privateKey, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("Failed to read private key: %w", err)
	}

	signer, err := ssh.ParsePrivateKey(privateKey)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse private key: %w", err)
	}

	return ssh.PublicKeys(signer), nil
}

// DefaultSSHConfig makes default SSH config.
func DefaultSSHConfig() (*ssh.ClientConfig, error) {
	user, err := user.Current()
	if err != nil {
		return nil, fmt.Errorf("Failed to get user info: %w", err)
	}

	keyAuth, err := GetSSHKey(filepath.Join(user.HomeDir, ".ssh/id_rsa"))
	if err != nil {
		return nil, err
	}

	userName := user.Username
	if runtime.GOOS == "windows" {
		// On Windows Username is "domain\user".
		userName = strings.SplitN(userName, `\`, 2)[1]
	}

	config := &ssh.ClientConfig{
		User:            userName,
		Auth:            []ssh.AuthMethod{keyAuth},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		// BannerCallback:  ssh.BannerDisplayStderr(),
		Timeout: 5 * time.Second,
	}

	return config, nil
}
