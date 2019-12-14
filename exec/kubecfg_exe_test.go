// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package exec

import (
	"testing"

	"github.com/vmware-tanzu/crash-diagnostics/script"
)

// TestExecKUBECONFIG tests KUBECONFIG command against a running API server.
// If setup properly, comment out t.Skip()
// TODO setup end-2-end tests
func TestExecKUBECONFIG(t *testing.T) {
	t.Skip(`Skipping KUBECONFIG exec: it requires a running API server`)
	tests := []execTest{
		{
			name: "KUBECFG",
			source: func() string {
				return `
				FROM local
				KUBECONFIG $HOME/.kube/kind-config-kind
				`
			},
			exec: func(s *script.Script) error {
				e := New(s)
				if err := e.Execute(); err != nil {
					return err
				}

				return nil
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			runExecutorTest(t, test)
		})
	}
}
