package k8s

import (
	"fmt"
	"strings"

	"github.com/vladimirvivien/echo"
)

func FetchBastionIpAddress(clusterName, namespace, kubeConfigPath string) (string, error) {
	if namespace == "" {
		namespace = "default"
	}
	p := echo.New().RunProc(fmt.Sprintf(
		`kubectl get awscluster/%s -o jsonpath='{.status.bastion.publicIp}' --namespace %s --kubeconfig %s`,
		clusterName,
		namespace,
		kubeConfigPath,
	))

	if p.Err() != nil {
		return "", fmt.Errorf("kubectl get awscluster failed: %s: %s", p.Err(), p.Result())
	}

	result := strings.TrimSpace(p.Result())
	return strings.ReplaceAll(result, "'", ""), nil
}
