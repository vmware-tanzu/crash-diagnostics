package testing

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"path/filepath"

	"golang.org/x/crypto/ssh"

	"github.com/pkg/errors"
)

// GenerateKeyPair generates a public/private key pair and stores it in the directory passed as the input.
func GenerateKeyPair(directory string) error {
	priv, err := generatePrivateKey()
	if err != nil {
		return errors.Wrap(err, "could not generate private key")
	}
	err = ioutil.WriteFile(filepath.Join(directory, "id_rsa"), encodePrivateKeyToPEM(priv), 0600)
	if err != nil {
		return errors.Wrap(err, "could not write public key to file")
	}

	pub, err := generatePublicKey(&priv.PublicKey)
	if err != nil {
		return errors.Wrap(err, "could not generate public key")
	}
	err = ioutil.WriteFile(filepath.Join(directory, "id_rsa.pub"), pub, 0600)
	if err != nil {
		return errors.Wrap(err, "could not write public key to file")
	}
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
