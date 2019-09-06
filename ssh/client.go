package ssh

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"gitlab.eng.vmware.com/vivienv/flare/script"
	"golang.org/x/crypto/ssh"
)

type SSHClient struct {
	cfg     *ssh.ClientConfig
	hostKey ssh.PublicKey
	client  *ssh.Client
}

func NewSSHClient(sshCmd *script.SSHConfigCommand) (*SSHClient, error) {
	if _, err := os.Stat(sshCmd.GetPrivateKeyPath()); err != nil {
		return nil, err
	}

	privateKey, err := ioutil.ReadFile(sshCmd.GetPrivateKeyPath())
	if err != nil {
		return nil, err
	}

	signer, err := ssh.ParsePrivateKey(privateKey)
	if err != nil {
		return nil, nil
	}

	if err != nil {
		log.Fatalf("unable to parse private key: %v", err)
	}

	client := &SSHClient{
		cfg: &ssh.ClientConfig{
			User: sshCmd.GetUserId(),
			Auth: []ssh.AuthMethod{
				ssh.PublicKeys(signer),
			},
		},
	}
	client.cfg.HostKeyCallback = ssh.FixedHostKey(client.hostKey)

	return client, nil
}

func (c *SSHClient) GetClientConfig() *ssh.ClientConfig {
	return c.cfg
}

func (c *SSHClient) Dial(addr string) error {
	client, err := ssh.Dial("tcp", addr, c.cfg)
	if err != nil {
		return err
	}
	c.client = client
	return nil
}

func (c *SSHClient) SSHRun(cmd string, args ...string) (io.Reader, error) {
	session, err := c.client.NewSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()

	output := new(bytes.Buffer)
	session.Stdout = output
	session.Stderr = output
	cmdStr := strings.TrimSpace(fmt.Sprintf("%s %s", cmd, strings.Join(args, " ")))
	if err := session.Run(cmdStr); err != nil {
		return nil, err
	}
	return output, nil
}

func (c *SSHClient) Hangup() error {
	return c.client.Close()
}
