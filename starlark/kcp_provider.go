package starlark

import (
	"context"
	"fmt"
	"net/url"

	"github.com/sirupsen/logrus"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"

	"github.com/vmware-tanzu/crash-diagnostics/k8s"
	"github.com/vmware-tanzu/crash-diagnostics/kcp"

	"github.com/kcp-dev/logicalcluster/v3"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

const TunnelHost = "localhost"

// KcpProviderFn is a built-in starlark function that creates a kubeconfig with all contexts for all KCP workspaces
// Starlark format: kcp_provider(kcp_admin_secret_namespace="kcp-system", kcp_admin_secret_name="admin-kubeconfig", [kcp_cert_secret_name="admin-cert-data", kube_config=kube_config(), tunnel_config=tunnel_config()]])
func KcpProviderFn(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {

	var (
		kcp_admin_secret_name, kcp_cert_secret_name, kcp_admin_secret_namespace string
		kubeConfig                                                              *starlarkstruct.Struct
		tunnelConfig                                                            *starlarkstruct.Struct
	)
	logrus.Info(kwargs)

	if err := starlark.UnpackArgs("kcp_provider", args, kwargs,
		"kcp_admin_secret_namespace", &kcp_admin_secret_namespace,
		"kcp_admin_secret_name", &kcp_admin_secret_name,
		"kcp_cert_secret_name?", &kcp_cert_secret_name,
		"kube_config?", &kubeConfig,
		"tunnel_config?", &tunnelConfig,
	); err != nil {
		return starlark.None, fmt.Errorf("failed to unpack input arguments: %w", err)
	}

	ctx, ok := thread.Local(identifiers.scriptCtx).(context.Context)
	if !ok || ctx == nil {
		return starlark.None, fmt.Errorf("script context not found")
	}

	if kubeConfig == nil {
		kubeConfig = thread.Local(identifiers.kubeCfg).(*starlarkstruct.Struct)
	}

	path, err := getKubeConfigPathFromStruct(kubeConfig)
	if err != nil {
		return starlark.None, fmt.Errorf("failed to kubeconfig: %w", err)
	}

	clusterCtxName := getKubeConfigContextNameFromStruct(kubeConfig)

	clients, err := k8s.New(path, clusterCtxName)
	if err != nil {
		return starlark.None, fmt.Errorf("could not initialize search client: %w", err)
	}

	kcpAdminKubeConfig, err := fetchKCPAdminKubeConfig(ctx, clients.Typed, kcp_admin_secret_namespace, kcp_admin_secret_name)
	if err != nil {
		return starlark.None, fmt.Errorf("failed to fetch KCP Admin KubeConfig: %w", err)
	}

	caCertBytes, tlsCertBytes, tlsKeyBytes, err := fetchKCPAdminKubeConfigCertData(ctx, clients.Typed, kcp_admin_secret_namespace, kcp_cert_secret_name)
	if err != nil {
		return starlark.None, fmt.Errorf("failed to fetch KCP Admin KubeConfig CertData: %w", err)
	}

	var localPort starlark.Int

	if tunnelConfig != nil {
		logrus.Info("tunnel config is passed. running the port forwarder.")

		var err error

		portforwardStr, err := tunnelConfig.Attr("namespace")
		if err != nil {
			return starlark.None, fmt.Errorf("could not parse the pod namespace: %w", err)
		}

		portForwardNS := portforwardStr.(starlark.String)

		podNameStr, err := tunnelConfig.Attr("pod_name")
		if err != nil {
			return starlark.None, fmt.Errorf("could not parse the pod name: %w", err)
		}

		portForwardPodName := podNameStr.(starlark.String)

		localPort, err = parseLocalPort(tunnelConfig)
		if err != nil {
			return starlark.None, fmt.Errorf("could not parse the local port: %w", err)
		}

		targetPortInt, err := tunnelConfig.Attr("target_port")
		if err != nil {
			return starlark.None, fmt.Errorf("could not get set the roundtripper: %w", err)
		}
		targetPort := targetPortInt.(starlark.Int)

		stopCh, readyCh := make(chan struct{}), make(chan struct{})

		fw, err := k8s.NewPortForwarder(path, string(portForwardNS), string(portForwardPodName), int(localPort.BigInt().Int64()), int(targetPort.BigInt().Int64()), stopCh, readyCh)
		if err != nil {
			return starlark.None, err
		}

		defer fw.Close()

		go func() {
			if err := fw.ForwardPorts(); err != nil {
				logrus.Error("error in forwarding ports ", err)
			}
		}()

		<-readyCh
	}

	for clusterName, clusterConfig := range kcpAdminKubeConfig.Clusters {
		clusterConfig.CertificateAuthority = ""
		clusterConfig.CertificateAuthorityData = caCertBytes

		if tunnelConfig != nil {
			tunneledServerAddress := fmt.Sprintf("%s:%s", TunnelHost, localPort)

			serverURL, err := url.Parse(clusterConfig.Server)
			if err != nil {
				return starlark.None, fmt.Errorf("error parsing server URL: %v", err)
			}

			serverURL.Host = tunneledServerAddress
			clusterConfig.Server = serverURL.String()
		}

		kcpAdminKubeConfig.Clusters[clusterName] = clusterConfig
	}

	for user, userAuthInfo := range kcpAdminKubeConfig.AuthInfos {
		userAuthInfo.ClientCertificate = ""
		userAuthInfo.ClientKey = ""

		userAuthInfo.ClientCertificateData = tlsCertBytes
		userAuthInfo.ClientKeyData = tlsKeyBytes

		kcpAdminKubeConfig.AuthInfos[user] = userAuthInfo
	}

	kcpCoreClient, err := kcp.NewFromConfig(*kcpAdminKubeConfig)
	if err != nil {
		return starlark.None, fmt.Errorf("failed to create KCP client: %w", err)
	}

	walker := kcp.WorkspaceWalker{KCPClusterClient: kcpCoreClient, Root: logicalcluster.NewPath("root")}
	workspaces, err := walker.FetchAllWorkspaces(ctx)
	if err != nil {
		return starlark.None, fmt.Errorf("failed to walk the KCP workspace tree: %w", err)
	}
	wsNames := make([]starlark.Value, 0)
	for _, ws := range workspaces {
		wsNames = append(wsNames, starlark.String(ws.String()))
	}

	if err := populateContexts(kcpAdminKubeConfig, workspaces); err != nil {
		return starlark.None, fmt.Errorf("failed to populate")
	}

	populatedAdminKubeConfigBytes, err := clientcmd.Write(*kcpAdminKubeConfig)
	if err != nil {
		return starlark.None, fmt.Errorf("error converting config to bytes: %w", err)
	}

	kcpProviderDict := starlark.StringDict{
		"kind":           starlark.String(identifiers.kcpProvider),
		"kube_config":    starlark.String(path),
		"kcp_kubeconfig": starlark.String(populatedAdminKubeConfigBytes),
		"contexts":       starlark.NewList(wsNames),
	}

	return starlarkstruct.FromStringDict(starlark.String(identifiers.kcpProvider), kcpProviderDict), nil
}

func fetchKCPAdminKubeConfig(ctx context.Context, k8sClient kubernetes.Interface, namespace string, name string) (*clientcmdapi.Config, error) {
	adminKubeConfigSecretData, err := k8s.GetSecretData(ctx, k8sClient, namespace, name, []string{"kubeconfig"})
	if err != nil {
		return nil, fmt.Errorf("failed to get admin kubeconfig: %w", err)
	}

	adminKubeConfigBytes, ok := adminKubeConfigSecretData["kubeconfig"]
	if !ok {
		return nil, fmt.Errorf("admin kubeconfig secret doesn't have kubeconfig")
	}

	adminKubeConfig, err := clientcmd.Load(adminKubeConfigBytes)
	if err != nil {
		return nil, fmt.Errorf("error converting admin kubeconfig bytes to struct: %v", err)
	}

	return adminKubeConfig, nil
}

func fetchKCPAdminKubeConfigCertData(ctx context.Context, k8sClient kubernetes.Interface, namespace string, name string) ([]byte, []byte, []byte, error) {
	certSecretData, err := k8s.GetSecretData(ctx, k8sClient, namespace, name, []string{"ca.crt", "tls.crt", "tls.key"})
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get admin kubeconfig ")
	}

	caCertBytes, ok := certSecretData["ca.crt"]
	if !ok {
		return nil, nil, nil, fmt.Errorf("cert secret doesn't have ca certificate")
	}

	tlsCertBytes, ok := certSecretData["tls.crt"]
	if !ok {
		return nil, nil, nil, fmt.Errorf("cert secret doesn't have certificate")
	}

	tlsKeyBytes, ok := certSecretData["tls.key"]
	if !ok {
		return nil, nil, nil, fmt.Errorf("cert secret doesn't have certificate key")
	}

	return caCertBytes, tlsCertBytes, tlsKeyBytes, nil
}

func populateContexts(kcpAdminKubeConfig *clientcmdapi.Config, workspaces []logicalcluster.Path) error {
	baseCluster := kcpAdminKubeConfig.Clusters["base"]
	baseContext := kcpAdminKubeConfig.Contexts["base"]

	baseUrl, err := url.Parse(baseCluster.Server)
	if err != nil {
		return err
	}

	for _, wsName := range workspaces {
		cluster := clientcmdapi.NewCluster()
		baseUrl.Path = "/clusters/" + wsName.String()
		cluster.Server = baseUrl.String()
		cluster.CertificateAuthorityData = baseCluster.CertificateAuthorityData

		kcpAdminKubeConfig.Clusters[wsName.String()] = cluster

		context := clientcmdapi.NewContext()
		context.Cluster = wsName.String()
		context.AuthInfo = baseContext.AuthInfo

		kcpAdminKubeConfig.Contexts[wsName.String()] = context
	}

	return nil
}

func parseLocalPort(tunnelConfig *starlarkstruct.Struct) (starlark.Int, error) {
	localPortInt, err := tunnelConfig.Attr("local_port")
	if err != nil {
		return starlark.MakeInt(0), fmt.Errorf("could not get set the roundtripper: %w", err)
	}

	localPort := localPortInt.(starlark.Int)

	return localPort, err
}
