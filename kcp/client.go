package kcp

import (
	"fmt"
	"net/url"

	"k8s.io/client-go/tools/clientcmd"

	kcpclientset "github.com/kcp-dev/kcp/sdk/client/clientset/versioned/cluster"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

func NewFromConfig(config clientcmdapi.Config) (kcpclientset.ClusterInterface, error) {
	// Convert api.Config to rest.Config
	clientConfig := clientcmd.NewDefaultClientConfig(config, &clientcmd.ConfigOverrides{})
	kcpRestConfig, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to create kcp REST config: %w", err)
	}

	u, err := url.Parse(kcpRestConfig.Host)
	if err != nil {
		return nil, err
	}
	u.Path = ""
	kcpRestConfig.Host = u.String()

	kcpCoreClient, err := kcpclientset.NewForConfig(kcpRestConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create kcpCoreClient: %w", err)
	}

	return kcpCoreClient, nil
}
