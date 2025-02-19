package kcp

import (
	"context"
	"fmt"

	kcpclientset "github.com/kcp-dev/kcp/sdk/client/clientset/versioned/cluster"
	"github.com/kcp-dev/logicalcluster/v3"

	pluginhelpers "github.com/kcp-dev/kcp/cli/pkg/helpers"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type WorkspaceWalker struct {
	KCPClusterClient kcpclientset.ClusterInterface
	Root             logicalcluster.Path
}

func (walker *WorkspaceWalker) FetchAllWorkspaces(ctx context.Context) ([]logicalcluster.Path, error) {
	contexts := make([]logicalcluster.Path, 0)
	contexts = append(contexts, walker.Root)

	children, err := walker.fetchChildren(ctx, walker.Root)
	if err != nil {
		return []logicalcluster.Path{}, err
	}

	contexts = append(contexts, children...)

	return contexts, nil
}

func (walker *WorkspaceWalker) fetchChildren(ctx context.Context, parent logicalcluster.Path) ([]logicalcluster.Path, error) {
	children := make([]logicalcluster.Path, 0)

	results, err := walker.KCPClusterClient.TenancyV1alpha1().Workspaces().Cluster(parent).List(ctx, metav1.ListOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			return []logicalcluster.Path{}, nil
		}
		return []logicalcluster.Path{}, err
	}

	for _, workspace := range results.Items {
		_, _, err := pluginhelpers.ParseClusterURL(workspace.Spec.URL)
		if err != nil {
			return []logicalcluster.Path{}, fmt.Errorf("current config context URL %q does not point to workspace", workspace.Spec.URL)
		}

		fullName := parent.String() + ":" + workspace.Name
		children = append(children, logicalcluster.NewPath(fullName))
		grandchildren, err := walker.fetchChildren(ctx, logicalcluster.NewPath(fullName))
		if err != nil {
			return []logicalcluster.Path{}, err
		}

		children = append(children, grandchildren...)
	}

	return children, nil
}
