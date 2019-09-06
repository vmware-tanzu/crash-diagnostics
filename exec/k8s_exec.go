package exec

import (
	"os"

	"github.com/sirupsen/logrus"
	"gitlab.eng.vmware.com/vivienv/flare/k8s"
	"gitlab.eng.vmware.com/vivienv/flare/script"
)

func exeClusterInfo(src *script.Script, path string) {
	cfgs, ok := src.Preambles[script.CmdKubeConfig]
	if !ok {
		logrus.Warn("Skipping cluster-info, KUBECONFIG not provided")
		return
	}
	cfgCmd := cfgs[0].(*script.KubeConfigCommand)
	if _, err := os.Stat(cfgCmd.Config()); err != nil {
		logrus.Warnf("Skipping cluster-info, unable to load KUBECONFIG %s: %s", cfgCmd.Config(), err)
		return
	}

	logrus.Debugf("Using KUBECONFIG %s", cfgCmd.Config())
	k8sClient, err := k8s.GetClient(cfgCmd.Config())
	if err != nil {
		logrus.Errorf("Skipping cluster-info, failed to create Kubernetes API server client: %s", err)
		return
	}

	if err := k8s.DumpClusterInfo(k8sClient, path); err != nil {
		logrus.Errorf("Failed to retrieve cluster information: %s", err)
		return
	}
}
