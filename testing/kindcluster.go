// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package testing

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/vladimirvivien/gexe"
)

var (
	kindVersion = "v0.14.0"
)

type KindCluster struct {
	name       string
	config     string
	tmpRootDir string
	e          *gexe.Echo
}

func NewKindCluster(config, name, tmpRootDir string) *KindCluster {
	return &KindCluster{name: name, config: config, tmpRootDir: tmpRootDir, e: gexe.New()}
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

func (k *KindCluster) MakeKubeConfigFile(path string) error {
	logrus.Infof("Creating kind kubeconfig file: %s", path)
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("failed to initialize kind kubeconfig file: %s", err)
	}
	defer f.Close()

	logrus.Infof("Retrieving kind kubeconfig for cluster: kind get kubeconfig --name %s", k.name)
	p := k.e.RunProc(fmt.Sprintf(`kind get kubeconfig --name %s`, k.name))
	if p.Err() != nil {
		return fmt.Errorf("failed to generate kind kubeconfig: %s: %s", p.Result(), p.Err())
	}

	if _, err := io.Copy(f, p.Out()); err != nil {
		return fmt.Errorf("failed to write kind kubeconfig file: %s", err)
	}

	logrus.Infof("kind kubeconfig file created: %s", f.Name())
	return nil
}

func (k *KindCluster) SimulateTerminatingPod() error {
	logrus.Infof("Simulating terminating pod in kind cluster %s", k.name)
	podConfig := `
apiVersion: v1
kind: Pod
metadata:
  name: stuck-pod
  namespace: default
  labels:
    app: test
  finalizers:
    - example.com/finalizer
spec:
  containers:
    - name: busybox
      image: busybox
      command:
        - sh
        - -c
        - while true; do echo "Simulating a stuck pod"; sleep 5; done
---
apiVersion: v1
kind: Pod
metadata:
  name: non-stuck-pod
  namespace: default
  labels:
    app: test
spec:
  containers:
    - name: busybox
      image: busybox
      command:
        - sh
        - -c
        - while true; do echo "Simulating a non-stuck pod"; sleep 5; done
`
	// Write pod configuration to a temporary file in the directory k.tmpRootDir
	filePath := fmt.Sprintf("%s/stuck-pod.yaml", k.tmpRootDir)
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to create file for pod configuration: %s", err)
	}
	defer file.Close()

	if _, err := file.WriteString(podConfig); err != nil {
		return fmt.Errorf("failed to write pod configuration to file: %s", err)
	}
	p := k.e.RunProc(fmt.Sprintf(`kubectl --context kind-%s apply -f %s`, k.name, filePath))
	if p.Err() != nil {
		return fmt.Errorf("failed to apply pod configuration: %s: %s", p.Err(), p.Result())
	}

	p = k.e.RunProc(fmt.Sprintf("kubectl --context kind-%s wait --for=condition=Ready pod -l app=test --timeout=60s", k.name))
	if p.Err() != nil {
		return fmt.Errorf("failed to simulate terminating pod: %s: %s", p.Err(), p.Result())
	}

	p = k.e.RunProc(fmt.Sprintf(`kubectl --context kind-%s delete pod stuck-pod --wait=false --grace-period=0 --force`, k.name))
	if p.Err() != nil {
		return fmt.Errorf("failed to simulate terminating pod: %s: %s", p.Err(), p.Result())
	}

	// Wait until pod is in Error state or max 10 seconds
	timeout := time.After(10 * time.Second)
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop() // Ensure the ticker is stopped to prevent resource leaks

	for {
		select {
		case <-timeout:
			return fmt.Errorf("timed out waiting for pod to be in Terminating state")
		case <-ticker.C:
			p = k.e.RunProc(fmt.Sprintf(`kubectl --context kind-%s get pod stuck-pod -o jsonpath='{.status.phase}'`, k.name))
			if p.Err() != nil {
				return fmt.Errorf("failed to check pod status: %s: %s", p.Err(), p.Result())
			}
			if strings.Contains(p.Result(), "Failed") {
				logrus.Infof("Pod is in Error state: %s", p.Result())
				return nil
			}
		}
	}
}

func (k *KindCluster) StartNginxPod() error {
	logrus.Infof("Starting pod in kind cluster %s", k.name)
	podConfig := `
apiVersion: v1
kind: Pod
metadata:
  name: nginx
  labels:
    app: nginx
spec:
  containers:
  - name: nginx
    image: nginx
    ports:
    - containerPort: 80

`
	
	filePath := fmt.Sprintf("%s/nginx-pod.yaml", k.tmpRootDir)
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to create file for pod configuration: %s", err)
	}
	defer file.Close()

	if _, err := file.WriteString(podConfig); err != nil {
		return fmt.Errorf("failed to write pod configuration to file: %s", err)
	}
	p := k.e.RunProc(fmt.Sprintf("timeout 60s bash -c 'while ! kubectl --context kind-%s get sa default -n default &>/dev/null; do sleep 1; done'", k.name))
	if p.Err() != nil {
		return fmt.Errorf("default service account has not been created: %s: %s", p.Err(), p.Result())
	}

	p = k.e.RunProc(fmt.Sprintf(`kubectl --context kind-%s apply -f %s`, k.name, filePath))
	if p.Err() != nil {
		return fmt.Errorf("failed to apply pod configuration: %s: %s", p.Err(), p.Result())
	}

	p = k.e.RunProc(fmt.Sprintf("kubectl --context kind-%s wait --for=condition=Ready pod -l app=nginx --timeout=60s", k.name))
	if p.Err() != nil {
		return fmt.Errorf("faild to schedule a pod: %s: %s", p.Err(), p.Result())
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

func findOrInstallKind(e *gexe.Echo) error {
	if len(e.Prog().Avail("kind")) == 0 {
		logrus.Info(`kind not found, installing with GO111MODULE="on" go get sigs.k8s.io/kind@v0.7.0`)
		if err := installKind(e); err != nil {
			return err
		}
	}
	return nil
}
func installKind(e *gexe.Echo) error {
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
