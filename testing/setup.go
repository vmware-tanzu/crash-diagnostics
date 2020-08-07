// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package testing

import (
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/vladimirvivien/echo"
)

const charset = "abcdefghijklmnopqrstuvwxyz"

var (
	InfraSetupWait   = time.Second * 11
	rnd              = rand.New(rand.NewSource(time.Now().Unix()))
	sshContainerName = "test-sshd"
	sshPort          = NextPortValue()
)

type TestSupport struct {
	username       string
	portValue      string
	resourceName   string
	testingRoot    string
	workdirRoot    string
	tmpDirRoot     string
	sshPKFileName  string
	sshPKFilePath  string
	maxConnRetries int
	sshServer      *SSHServer
	kindKubeCfg    string
	kindCluster    *KindCluster
}

// Init initializes and returns TestSupport instance
func Init() (*TestSupport, error) {
	debug := false
	flag.BoolVar(&debug, "debug", debug, "Enables debug level")
	flag.Parse()
	e := echo.New()

	logLevel := logrus.InfoLevel
	if debug {
		logLevel = logrus.DebugLevel
	}
	logrus.SetLevel(logLevel)

	// get username
	username, err := Username()
	if err != nil {
		return nil, err
	}

	resource := NextResourceName()

	// setup workdir
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	testingRoot := filepath.Join(homeDir, ".crashd-testing", resource)
	if err := os.MkdirAll(testingRoot, 0765); err != nil && !os.IsExist(err) {
		return nil, err
	}
	logrus.Infof("Created testing root dir: %s", testingRoot)

	workDir := filepath.Join(testingRoot, "work")
	if err := os.MkdirAll(workDir, 0765); err != nil && !os.IsExist(err) {
		return nil, err
	}
	logrus.Infof("Created testing work dir: %s", workDir)

	sshKeyPath, err := filepath.Abs(filepath.Join("..", "testing"))
	if err != nil {
		return nil, err
	}
	cpCmd := fmt.Sprintf(`/bin/sh -c "cp %s/id_rsa* %s"`, sshKeyPath, workDir)
	logrus.Infof("Copying SSH key files: %s", cpCmd)
	proc := e.RunProc(cpCmd)
	if proc.Err() != nil {
		logrus.Errorf("Error copying key files: %s %s", proc.Err(), proc.Result())
		return nil, proc.Err()
	}

	// setup tempDir
	tmpDirRoot := filepath.Join(testingRoot, "tmp")
	if err := os.MkdirAll(tmpDirRoot, 0765); err != nil && !os.IsExist(err) {
		return nil, err
	}
	logrus.Infof("Created testing temp root dir: %s", tmpDirRoot)

	pkName := "id_rsa"
	pkPath := filepath.Join(workDir, pkName)
	return &TestSupport{
		username:       username,
		portValue:      NextPortValue(),
		resourceName:   resource,
		testingRoot:    testingRoot,
		workdirRoot:    workDir,
		tmpDirRoot:     tmpDirRoot,
		sshPKFileName:  pkName,
		sshPKFilePath:  pkPath,
		maxConnRetries: 100,
	}, nil
}

// PortValue returns a string with a random value that can be used as port
func (t *TestSupport) PortValue() string {
	return t.portValue
}

// ResourceName resturns string that can be used to name resource
func (t *TestSupport) ResourceName() string {
	return t.resourceName
}

// CurrentUsername returns the current username or error
func (t *TestSupport) CurrentUsername() string {
	return t.username
}

func (t *TestSupport) WorkDirRoot() string {
	return t.workdirRoot
}

func (t *TestSupport) TmpDirRoot() string {
	return t.tmpDirRoot
}

func (t *TestSupport) PrivateKeyPath() string {
	return t.sshPKFilePath
}

func (t *TestSupport) MaxConnectionRetries() int {
	return t.maxConnRetries
}

func (t *TestSupport) SetupSSHServer() error {
	if t.sshServer == nil {
		//privKeyPath := filepath.Join(t.workdirRoot, t.sshPKFileName)
		//if err := GenerateRSAKeyFiles(t.workdirRoot, t.sshPKFileName); err != nil {
		//	return err
		//}
		//
		//if err := AddKeyToAgent(privKeyPath); err != nil {
		//	logrus.Errorf("Failed to add private key to SSH agent: %s", err)
		//} else {
		//	logrus.Infof("Added private key to ssh-agent: %s ", privKeyPath)
		//}

		server, err := NewSSHServer(t.resourceName, t.username, t.portValue, t.workdirRoot)
		if err != nil {
			return err
		}

		if err := server.Start(); err != nil {
			return err
		}

		t.sshServer = server
	}
	return nil
}

func (t *TestSupport) SetupKindCluster() error {
	if t.kindCluster == nil {
		yamlPath, err := filepath.Abs(filepath.Join("..", "./testing", "/kind-cluster-docker.yaml"))
		if err != nil {
			return err
		}

		kind := NewKindCluster(yamlPath, t.resourceName)
		if err := kind.Create(); err != nil {
			return err
		}
		logrus.Infof("kind cluster created")

		// stall to wait for kind pods initialization
		waitTime := time.Second * 10
		logrus.Debugf("waiting %s for kind pods to initialize...", waitTime)
		time.Sleep(waitTime)

		t.kindCluster = kind
	}
	return nil
}

func (t *TestSupport) SetupKindKubeConfig() (string, error) {
	if t.kindCluster == nil {
		return "", fmt.Errorf("kind not set: call SetupKindCluster() first")
	}

	if len(t.kindKubeCfg) > 0 {
		return t.kindKubeCfg, nil
	}

	kubeCfgFile := filepath.Join(t.tmpDirRoot, "kubeconfig")
	if err := t.kindCluster.MakeKubeConfigFile(kubeCfgFile); err != nil {
		return "", err
	}
	t.kindKubeCfg = kubeCfgFile
	return kubeCfgFile, nil
}

func (t *TestSupport) KindKubeConfigFile() string {
	return t.kindKubeCfg
}

func (t *TestSupport) TearDown() error {
	var errs []error

	if t.kindCluster != nil {
		logrus.Infof("Destroying kind cluster...")
		if err := t.kindCluster.Destroy(); err != nil {
			logrus.Error(err)
			errs = append(errs, err)
		}
	}

	//privKeyPath := filepath.Join(t.workdirRoot, t.sshPKFileName)
	//logrus.Infof("Removing private key from agent: %s", privKeyPath)
	//if err := RemoveKeyFromAgent(privKeyPath); err != nil {
	//	logrus.Errorf("Unable to remove private key from SSH agent: %s", err)
	//}

	if t.sshServer != nil {
		logrus.Infof("Stopping SSH server container....")
		if err := t.sshServer.Stop(); err != nil {
			logrus.Error(err)
			errs = append(errs, err)
		}
		time.Sleep(time.Millisecond * 500)
	}

	logrus.Infof("Removing dir: %s", t.testingRoot)
	if err := os.RemoveAll(t.testingRoot); err != nil {
		// do return err:
		// ssh-server container does not cleanly release mounted dir
		// workaround to GitHub Actions permission issue during tests
		logrus.Errorf("Unable to remove testing root dir: %s", err)
	}

	if errs != nil {
		return errors.New(fmt.Sprintf("%v", errs))
	}

	return nil
}

//NextPortValue returns a pseudo-rando test [2200 .. 2290]
func NextPortValue() string {
	port := 2200 + rnd.Intn(90)
	return fmt.Sprintf("%d", port)
}

// NextResourceName returns crashd-test-XXXX name
func NextResourceName() string {
	return fmt.Sprintf("crashd-test-%x", rnd.Uint64())
}

// Username returns current username
func Username() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return usr.Username, nil
}
