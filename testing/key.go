package testing

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/vladimirvivien/echo"
	"golang.org/x/crypto/ssh"

	"github.com/pkg/errors"
)

// GenerateRSAKeyFiles generates a public/private key pair and stores it in the directory passed as the input.
func GenerateRSAKeyFiles(directory, privFileName string) error {
	priv, err := generatePrivateKey()
	if err != nil {
		return errors.Wrap(err, "could not generate private key")
	}
	rsaFile := filepath.Join(directory, privFileName)
	err = ioutil.WriteFile(rsaFile, encodePrivateKeyToPEM(priv), 0600)
	if err != nil {
		return errors.Wrap(err, "could not write private key to file")
	}
	logrus.Info("Created private key PEM file:", rsaFile)

	pub, err := generatePublicKey(&priv.PublicKey)
	if err != nil {
		return errors.Wrap(err, "could not generate public key")
	}

	pubFileName := fmt.Sprintf("%s.pub", privFileName)
	rsaPubFile := filepath.Join(directory, pubFileName)
	err = ioutil.WriteFile(rsaPubFile, pub, 0600)
	if err != nil {
		return errors.Wrap(err, "could not write public key to file")
	}
	logrus.Info("Created public key file:", rsaPubFile)

	return nil
}

func generatePrivateKey() (*rsa.PrivateKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, err
	}

	// Validate Private Key
	err = privateKey.Validate()
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

// encodePrivateKeyToPEM encodes Private Key from RSA to PEM format
func encodePrivateKeyToPEM(privateKey *rsa.PrivateKey) []byte {
	// Get ASN.1 DER format
	privDER := x509.MarshalPKCS1PrivateKey(privateKey)

	// pem.Block
	privBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   privDER,
	}

	// Private key in PEM format
	privatePEM := pem.EncodeToMemory(&privBlock)

	return privatePEM
}

// generatePublicKey take a rsa.PublicKey and return bytes suitable for writing to .pub file
// returns in the format "ssh-rsa ..."
func generatePublicKey(privatekey *rsa.PublicKey) ([]byte, error) {
	publicRsaKey, err := ssh.NewPublicKey(privatekey)
	if err != nil {
		return nil, err
	}

	pubKeyBytes := ssh.MarshalAuthorizedKey(publicRsaKey)
	return pubKeyBytes, nil
}

func AddKeyToAgent(keyPath string) error {
	e := echo.New()

	logrus.Info("Starting ssh-agent if needed...")
	sshAgentCmd := e.Prog.Avail("ssh-agent")
	if len(sshAgentCmd) == 0 {
		return fmt.Errorf("ssh-agent not found")
	}
	var agentPID string
	if aid := e.Eval("$SSH_AGENT_PID"); len(agentPID) == 0 {
		proc := e.RunProc(fmt.Sprintf(`/bin/sh -c 'eval "$(%s)"'`, sshAgentCmd))
		if proc.Err() != nil {
			return fmt.Errorf("ssh-agent failed: %s: %s", proc.Err(), proc.Result())
		}
		result := proc.Result()
		logrus.Infof("ssh-agent started: %s", result)
		agentPID = strings.Split(result, " ")[2]
	} else {
		agentPID = aid
		logrus.Infof("ssh-agent pid found: %s", aid)
	}

	sshAddCmd := e.Prog.Avail("ssh-add")
	if len(sshAddCmd) == 0 {
		return fmt.Errorf("ssh-add not found")
	}

	logrus.Debugf("adding key to ssh-agent (pid %s): %s", agentPID, keyPath)

	e.SetVar("ssh_agent_pid", agentPID)
	p := e.RunProc(fmt.Sprintf(`/bin/sh -c 'SSH_AGENT_PID=%s %s %s'`, agentPID, sshAddCmd, keyPath))
	if p.Err() != nil {
		return fmt.Errorf("failed to add SSH key to agent: %s: %s", p.Err(), p.Result())
	}
	logrus.Infof("ssh-add result: %s", p.Result())
	return nil
}

func RemoveKeyFromAgent(keyPath string) error {
	e := echo.New()
	sshAddCmd := e.Prog.Avail("ssh-add")
	if len(sshAddCmd) == 0 {
		return fmt.Errorf("ssh-add not found")
	}
	logrus.Debugf("removing key from ssh-agent: %s", keyPath)
	p := e.RunProc(fmt.Sprintf("%s -d %s", sshAddCmd, keyPath))
	if p.Err() != nil {
		return fmt.Errorf("failed to remove SSH key from agent: %s: %s", p.Err(), p.Result())
	}
	logrus.Infof("removal key result: %s", p.Result())
	return nil
}
