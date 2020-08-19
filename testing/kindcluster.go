// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package testing

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/vladimirvivien/echo"
)

var (
	kindVersion = "v0.7.0"
	clusterName = "test-crashd-cluster"
)

type KindCluster struct {
	name   string
	config string
	e      *echo.Echo
}

func NewKindCluster(config, name string) *KindCluster {
	return &KindCluster{name: name, config: config, e: echo.New()}
}

func (k *KindCluster) Create() error {
	logrus.Infof("Creating kind cluster %s", k.name)
	// is kind program available
	if err := findOrInstallKind(k.e); err != nil {
		return err
	}

	if strings.Contains(k.e.Run("kind get clusters"), k.name) {
		logrus.Infof("Skipping KindCluster.Create: cluster already created: %s", k.name)
		return nil
	}

	// create kind cluster using kind-cluster-docker.yaml config file
	logrus.Infof("launching: kind create cluster --config %s --name %s", k.config, k.name)
	p := k.e.RunProc(fmt.Sprintf(`kind create cluster --config %s --name %s`, k.config, k.name))
	if p.Err() != nil {
		return fmt.Errorf("failed to install kind: %s: %s", p.Err(), p.Result())
	}

	clusters := k.e.Run("kind get clusters")
	logrus.Infof("kind clusters available: %s", clusters)

	return nil
}

func (k *KindCluster) GetKubeConfig() (io.Reader, error) {
	logrus.Infof("Retrieving kind kubeconfig for cluster: %s", k.name)
	p := k.e.RunProc(fmt.Sprintf(`kind get kubeconfig --name %s`, k.name))
	if p.Err() != nil {
		return nil, p.Err()
	}
	return p.Out(), nil
}

func (k *KindCluster) MakeKubeConfigFile(path string) error {
	logrus.Infof("Creating kind kubeconfig file: %s", path)
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("failed to initialize kind kubeconfig file: %s", err)
	}
	defer f.Close()

	reader, err := k.GetKubeConfig()
	if err != nil {
		return fmt.Errorf("failed to generate kind kubeconfig: %s", err)
	}
	if _, err := io.Copy(f, reader); err != nil {
		return fmt.Errorf("failed to write kind kubeconfig file: %s", err)
	}
	return nil
}

func (k *KindCluster) GetKubeCtlContext() string {
	return fmt.Sprintf("kind-%s", k.name)
}

func (k *KindCluster) Destroy() error {
	logrus.Infof("Destroying kind cluster %s", k.name)
	if err := findOrInstallKind(k.e); err != nil {
		return err
	}
	// deleteting kind cluster
	p := k.e.RunProc(fmt.Sprintf(`kind delete cluster --name %s`, k.name))
	if p.Err() != nil {
		return fmt.Errorf("failed to install kind: %s: %s", p.Err(), p.Result())
	}

	logrus.Info("Kind cluster destroyed")

	clusters := k.e.Run("kind get clusters")
	logrus.Infof("Available kind clusters: %s", clusters)

	return nil
}

func findOrInstallKind(e *echo.Echo) error {
	if len(e.Prog.Avail("kind")) == 0 {
		logrus.Info(`kind not found, installing with GO111MODULE="on" go get sigs.k8s.io/kind@v0.7.0`)
		if err := installKind(e); err != nil {
			return err
		}
	}
	return nil
}
func installKind(e *echo.Echo) error {
	logrus.Infof("installing: go get sigs.k8s.io/kind@%s", kindVersion)
	p := e.SetEnv("GO111MODULE", "on").RunProc(fmt.Sprintf("go get sigs.k8s.io/kind@%s", kindVersion))
	if p.Err() != nil {
		return fmt.Errorf("failed to install kind: %s: %s", p.Err(), p.Result())
	}
	if !p.IsSuccess() || p.ExitCode() != 0 {
		return fmt.Errorf("failed to install kind: %s", p.Result())
	}
	return nil
}
