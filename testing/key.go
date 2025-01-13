package testing

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

// WriteKeys copies the static private key in variable privateKey
// to a local file.
func WriteKeys(rootPath string) error {
	// Use 0444 on github and 0600 on local dev because for some reason the docker volume mount in the action runner requires 0444 to work with the file
	isGitHubActions := os.Getenv("GITHUB_ACTIONS") == "true"
	var keyPerm os.FileMode
	if isGitHubActions {
		keyPerm = 0444
	} else {
		keyPerm = 0600
	}
	pkPath := filepath.Join(rootPath, "id_rsa")
	pkFile, err := os.OpenFile(pkPath, os.O_RDWR|os.O_CREATE, keyPerm)
	if err != nil {
		return err
	}
	defer pkFile.Close()
	logrus.Infof("Writing private key file: %s", pkPath)
	if _, err := io.Copy(pkFile, strings.NewReader(privateKey)); err != nil {
		return fmt.Errorf("")
	}

	pubPath := filepath.Join(rootPath, "id_rsa.pub")
	pubFile, err := os.OpenFile(pubPath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer pubFile.Close()
	logrus.Infof("Writing private key file: %s", pubPath)
	if _, err := io.Copy(pubFile, strings.NewReader(publicKey)); err != nil {
		return fmt.Errorf("")
	}
	return nil
}

var privateKey = `-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----`

var publicKey = `ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDG5ob1agkTUJVc1sAGG9Ey+YEJiYzai5K0QFy7fRTp590pUwDOdatd9zTA7tE3fhhHqZOx+5YKNAheBco3jUu9SY5MaFpw/raXZblKgqULmq6JDQ1Xs4sT8FcPpD+TGNfEKQZEBGtZgZK3WMilUw9JmOMn25w858lLonz9vH5OPvrPh/MS8Gk7//pb/6SgHAKBlfhiSMiCqVxWhsK6IQua4kT8hsGV1gIFYNXEtKPPRUom07mp/CzJSxKbJbGJqpN4sSuSD8FA1LefNnDarL7JcartoR47WHm2ZuXaARZlbYVDFVIbsBR4u0gCzoOMlsPpwa2+RUSesBPh2+xtBZiGEKcmkWjn/gq48LwxRChtJB9aOgQ4a6msea+nVRWSnVl4c9kk6k/0QoWwQ5V3UmCxal8bg8b+wKWnA2bWWIJlNPxhvhSfmYGpOIyEIiPlFGZPEsZUzQaULrXQe3zuKHGLgoE9eKhKh1yG/sHDxzwy5iKEDKmluiPuNi0kL344ElxtL52C8ovP0HtcTg/5wvM+hNjqe7WdhYrEe8XScvxWojhrBvFDbBlWS14mfpejrcQGLcIqOZ6UDUUnasFEfGoacCSe8mSoj+l4mDx9ZX7a9IHbn8xZn6uEPIV7feY8Q7luVldGGeZfHOW1omwl+Ibc7Msd6Pvgp6ler30VKojlNQ==`
