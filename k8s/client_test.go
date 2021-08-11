// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package k8s

import (
	"context"
	"testing"
)

func TestClient_Search(t *testing.T) {
	tests := []struct {
		name       string
		params     SearchParams
		shouldFail bool
		eval       func(t *testing.T, results []SearchResult)
	}{
		{
			name:       "empty params",
			params:     SearchParams{},
			shouldFail: true,
		},
		{
			name:   "groups only",
			params: SearchParams{Groups: []string{"apps"}},
			eval: func(t *testing.T, results []SearchResult) {
				if len(results) == 0 {
					t.Errorf("no objects found")
				}
				count := 0
				for _, result := range results {
					count = len(result.List.Items) + count
				}
				t.Logf("found %d objects", count)
			},
		},
		{
			name:   "kinds (resources) only",
			params: SearchParams{Kinds: []string{"pods"}},
			eval: func(t *testing.T, results []SearchResult) {
				if len(results) == 0 {
					t.Errorf("no objects found")
				}
				count := 0
				for _, result := range results {
					count = len(result.List.Items) + count
				}
				t.Logf("found %d objects", count)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client, err := New(support.KindKubeConfigFile())
			if err != nil {
				t.Fatal(err)
			}
			results, err := client.Search(context.TODO(), test.params)
			if err != nil && !test.shouldFail {
				t.Fatal(err)
			}
			if test.eval != nil {
				test.eval(t, results)
			}
		})
	}
}
