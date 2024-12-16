package starlark

import (
	"context"
	"fmt"
	kcpclient "github.com/kcp-dev/kcp/pkg/client/clientset/versioned"
	"github.com/kcp-dev/logicalcluster/v3"
	"github.com/pkg/errors"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

const (
	kcpPreviousWorkspaceContextKey string = "workspace.kcp.io/previous"
	kcpCurrentWorkspaceContextKey  string = "workspace.kcp.io/current"
)

type UseWorkspaceCommand struct {

	// Name is the name of the workspace to switch to.
	Name string
	// ShortWorkspaceOutput indicates only the workspace name should be printed.
	ShortWorkspaceOutput bool

	modifyConfig   func(name string, newConfig *clientcmdapi.Config) (string, error)
	adminUcpConfig *clientcmdapi.Config
	kubeconfigPath string
	kcpclient      *kcpclient.Clientset
	workspace      string
}

// NewUseWorkspaceCommand returns a new UseWorkspaceCommand.
func NewUseWorkspaceCommand(workspace string, kubeConfig *clientcmdapi.Config, kubeconfigPath string) *UseWorkspaceCommand {
	kcpConfig, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfigPath},
		&clientcmd.ConfigOverrides{},
	).ClientConfig()
	if err != nil {
		log.Fatalf("Failed to create KCP client config: %v", err)
	}

	kcpClient, err := kcpclient.NewForConfig(kcpConfig)
	if err != nil {
		log.Fatalf("Failed to create KCP client: %v", err)
	}

	return &UseWorkspaceCommand{
		kubeconfigPath: kubeconfigPath,
		adminUcpConfig: kubeConfig,
		kcpclient:      kcpClient,
		Name:           workspace,

		modifyConfig: func(workspace string, newConfig *clientcmdapi.Config) (string, error) {
			// Convert *api.Config back to raw string format
			rawConfig, err := clientcmd.Write(*newConfig)
			if err != nil {
				log.Fatalf("Failed to write kubeconfig: %v", err)
			}

			// Print the raw kubeconfig string
			fmt.Println(string(rawConfig))
			parentDir := os.TempDir()
			ucpConfigs := filepath.Join(parentDir, "ucp-kubeconfigs")
			kubeconfig, err := os.CreateTemp(ucpConfigs, workspace+".*.kubeconfig")
			if err != nil {
				log.Fatalf("Failed to create temporary file: %v", err)
			}

			err = clientcmd.WriteToFile(*newConfig, kubeconfig.Name())
			if err != nil {
				return "", err
			}
			return kubeconfig.Name(), nil
		},
	}
}

// Run executes the "use workspace" logic based on the supplied options.
func (o *UseWorkspaceCommand) Run(ctx context.Context) (kubeconfigPath string, err error) {
	name := o.Name

	cluster := o.adminUcpConfig.Clusters["workspace.kcp.io/current"]
	currentContext, found := o.adminUcpConfig.Contexts[o.adminUcpConfig.CurrentContext]
	if !found {
		return "", fmt.Errorf("current %q context not found", currentContext)
	}

	// make relative paths absolute
	var workspacePath string
	if name[0] == ':' {
		workspacePath = strings.TrimPrefix(name, ":")
	} else {
		workspacePath = strings.TrimSpace(cluster.Server) + ":" + strings.TrimSpace(name)
	}
	//
	//// remove . and ..
	//pth, err := resolveDots(name)
	//if err != nil {
	//	return err
	//}
	//
	//// here we should have a valid absolute path without dots, without : prefix
	//if !pth.IsValid() {
	//	return fmt.Errorf("invalid workspace path: %s", o.Name)
	//}

	//// first check if the workspace exists via discovery
	//groups, err := o.kcpclient.Discovery().ServerGroups()
	//if err != nil && !apierrors.IsForbidden(err) {
	//	return err
	//}
	//denied := false
	//if apierrors.IsForbidden(err) || len(groups.Groups) == 0 {
	//	denied = true
	//}
	//
	//// first try to get Workspace from parent to potentially get a 404. A 403 in the parent though is
	//// not a blocker to enter the workspace. We do discovery as a final check.
	//
	//notFound := false
	//
	//_, workspaceName := logicalcluster.NewPath(name).Split()
	//if workspaceName != "" {
	//	if _, err := o.kcpclient.TenancyV1alpha1().Workspaces().Get(ctx, workspaceName, metav1.GetOptions{}); apierrors.IsNotFound(err) {
	//		notFound = true
	//	}
	//}
	//
	//switch {
	//case denied && notFound:
	//	return fmt.Errorf("workspace %q not found", name)
	//case denied:
	//	return fmt.Errorf("access to workspace %q denied", name)
	//case notFound:
	//	// we are good. Somehow we have access, maybe without having access to the parent or there is no parent.
	//}

	//u.Path = path.Join(u.Path, pth.RequestPath())
	u, err := url.Parse(workspacePath)
	if err != nil {
		return "", err
	}
	return o.commitConfig(name, currentContext, u)
}

func resolveDots(pth string) (logicalcluster.Path, error) {
	var ret logicalcluster.Path
	for _, part := range strings.Split(pth, ":") {
		switch part {
		case ".":
			continue
		case "..":
			if ret.Empty() {
				return logicalcluster.Path{}, errors.New("cannot go up from root")
			}
			ret, _ = ret.Parent()
		default:
			ret = ret.Join(part)
		}
	}
	return ret, nil
}

// swapContexts moves to previous context from the config.
// It will update existing configuration by swapping current & previous configurations.
// This method already commits. Do not use with commitConfig

// commitConfig will take in current config, new host and optional workspaceType and update the kubeconfig.
func (o *UseWorkspaceCommand) commitConfig(workspace string, currentContext *clientcmdapi.Context, u *url.URL) (string, error) {
	// modify kubeconfig, using the "workspace" context and cluster
	newKubeConfig := o.adminUcpConfig.DeepCopy()
	oldCluster, found := o.adminUcpConfig.Clusters[currentContext.Cluster]
	if !found {
		return "", fmt.Errorf("cluster %q not found in kubeconfig", currentContext.Cluster)
	}
	newCluster := *oldCluster
	newCluster.Server = u.String()
	newKubeConfig.Clusters[kcpCurrentWorkspaceContextKey] = &newCluster
	newContext := *currentContext
	newContext.Cluster = kcpCurrentWorkspaceContextKey
	newKubeConfig.Contexts[kcpCurrentWorkspaceContextKey] = &newContext

	// store old context and old cluster
	if currentContext.Cluster == kcpCurrentWorkspaceContextKey {
		currentContext = currentContext.DeepCopy()
		currentContext.Cluster = kcpPreviousWorkspaceContextKey
		newKubeConfig.Clusters[kcpPreviousWorkspaceContextKey] = oldCluster
	}
	newKubeConfig.Contexts[kcpPreviousWorkspaceContextKey] = currentContext

	newKubeConfig.CurrentContext = kcpCurrentWorkspaceContextKey

	kubeConfigPath, err := o.modifyConfig(workspace, newKubeConfig)
	if err != nil {
		return "", err
	}

	return kubeConfigPath, nil
}
