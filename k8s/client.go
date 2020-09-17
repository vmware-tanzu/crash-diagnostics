// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package k8s

import (
	"strings"

	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/cli-runtime/pkg/printers"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	LegacyGroupName = "core"
)

// Client prepares and exposes a dynamic, discovery, and Rest clients
type Client struct {
	Client      dynamic.Interface
	Disco       discovery.DiscoveryInterface
	CoreRest    rest.Interface
	JsonPrinter printers.JSONPrinter
}

// New returns a *Client
func New(kubeconfig string) (*Client, error) {
	// creating cfg for each client type because each
	// setup its own cfg default which may not be compatible
	dynCfg, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}
	client, err := dynamic.NewForConfig(dynCfg)
	if err != nil {
		return nil, err
	}

	discoCfg, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}
	disco, err := discovery.NewDiscoveryClientForConfig(discoCfg)
	if err != nil {
		return nil, err
	}

	restCfg, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}
	setCoreDefaultConfig(restCfg)
	restc, err := rest.RESTClientFor(restCfg)
	if err != nil {
		return nil, err
	}

	return &Client{Client: client, Disco: disco, CoreRest: restc}, nil
}

func (k8sc *Client) Search(params SearchParams) ([]SearchResult, error) {
	return k8sc._search(strings.Join(params.Groups, " "),
		strings.Join(params.Kinds, " "),
		strings.Join(params.Namespaces, " "),
		strings.Join(params.Versions, " "),
		strings.Join(params.Names, " "),
		strings.Join(params.Labels, " "),
		strings.Join(params.Containers, " "))
}

// Search does a drill-down search from group, version, resourceList, to resources.  The following rules are applied
// 1) Legacy core group (api/v1) can be specified as "core"
// 2) All specified search params will use AND operator for match (i.e. groups=core AND kinds=pods AND versions=v1 AND ... etc)
// 3) kinds will match resource.Kind or resource.Name
// 4) All search params are passed as comma- or space-separated sets that are matched using OR (i.e. kinds=pods services
//    will match resouces of type pods or services)
func (k8sc *Client) _search(groups, kinds, namespaces, versions, names, labels, containers string) ([]SearchResult, error) {
	// normalize params
	groups = strings.ToLower(groups)
	kinds = strings.ToLower(kinds)
	namespaces = strings.ToLower(namespaces)
	versions = strings.ToLower(versions)
	labels = strings.ToLower(labels)
	containers = strings.ToLower(containers)

	logrus.Debugf(
		"Search filters groups:[%v]; kinds:[%v]; namespaces:[%v]; versions:[%v]; names:[%v]; labels:[%v] containers:[%s]",
		groups, kinds, namespaces, versions, names, labels, containers,
	)

	grpList, err := k8sc.Disco.ServerGroups()
	if err != nil {
		return nil, err
	}

	// if namespace filters not provided, assume all namespaces
	if len(namespaces) == 0 {
		nsNames, err := getNamespaces(k8sc)
		if err != nil {
			return nil, err
		}
		namespaces = strings.Join(nsNames, " ")
	}

	var finalResults []SearchResult
	logrus.Debugf("Searching in %d groups", len(grpList.Groups))
	for _, grp := range grpList.Groups {
		// filter by group
		grpName := strings.TrimSpace(grp.Name)
		grpName = getLegacyGrpName(grpName)
		if len(groups) > 0 && !strings.Contains(groups, strings.ToLower(grpName)) {
			continue
		}

		// filter by group version
		logrus.Debugf("Searching resources in Group %s", grpName)
		for _, discoGV := range grp.Versions {
			if len(versions) > 0 && !strings.Contains(versions, strings.ToLower(discoGV.Version)) {
				continue
			}

			// adjust version for legacy group
			grpVersion := discoGV.GroupVersion
			if grpName == LegacyGroupName {
				grpVersion = discoGV.Version
			}

			logrus.Debugf("Searching resources in GroupVersion %s", discoGV.GroupVersion)
			resources, err := k8sc.Disco.ServerResourcesForGroupVersion(grpVersion)
			if err != nil {
				logrus.Errorf("K8s.Search failed to get resources for %s: %s", discoGV.GroupVersion, err)
				continue
			}

			// filter by resource kind
			for _, res := range resources.APIResources {
				if len(kinds) > 0 && !strings.Contains(kinds, strings.ToLower(res.Kind)) {
					continue
				}

				gvr := schema.GroupVersionResource{
					Group:    toLegacyGrpName(grpName),
					Version:  discoGV.Version,
					Resource: res.Name,
				}
				logrus.Debugf(`Searching for GroupResource %#v`, gvr)

				// retrieve API objects based on GroupVersionResource and
				// filter by namespaces and names
				listOptions := metav1.ListOptions{
					LabelSelector: labels,
				}

				// gather found resources
				var results []SearchResult
				if res.Namespaced {
					for _, ns := range splitParamList(namespaces) {
						logrus.Debugf("Searching for %s in namespace %s [GroupRes: %v]", res.Name, ns, gvr)
						list, err := k8sc.Client.Resource(gvr).Namespace(ns).List(listOptions)
						if err != nil {
							logrus.Debugf(
								"WARN: K8s.Search failed to get %s in %s [GroupRes: %s][labels: %v]: %s",
								res.Name, ns, discoGV.GroupVersion, listOptions.LabelSelector, err,
							)
							continue
						}
						logrus.Debugf("Found %d %s in namespace [%s]", len(list.Items), res.Name, ns)
						result := SearchResult{
							ListKind:             list.GetKind(),
							ResourceName:         res.Name,
							ResourceKind:         res.Kind,
							Namespaced:           res.Namespaced,
							Namespace:            ns,
							GroupVersionResource: gvr,
							List:                 list,
						}
						results = append(results, result)
					}
				} else {
					logrus.Debugf("Searching for resource %s (non-namespaced)", res.Name)
					list, err := k8sc.Client.Resource(gvr).List(listOptions)
					if err != nil {
						logrus.Debugf(
							"WARN: K8s.Search failed to get %s: [GroupRes: %s] [labels: %v]: %s",
							res.Name, discoGV.GroupVersion, listOptions.LabelSelector, err,
						)
						continue
					}
					logrus.Debugf("Found %d %s (non-namespaced)", len(list.Items), res.Name)
					result := SearchResult{
						ListKind:             list.GetKind(),
						ResourceKind:         res.Kind,
						ResourceName:         res.Name,
						Namespaced:           res.Namespaced,
						GroupVersionResource: gvr,
						List:                 list,
					}
					results = append(results, result)
				}

				// apply name filters
				for _, result := range results {
					filteredResult := result
					if len(containers) > 0 && result.ListKind == "PodList" {
						filteredResult = filterPodsByContainers(result, containers)
						logrus.Debugf("Found %d %s with container filter [%s]", len(filteredResult.List.Items), filteredResult.ResourceName, containers)
					}
					if len(names) > 0 {
						filteredResult = filterByNames(result, names)
						logrus.Debugf("Found %d %s with name filter [%s]", len(filteredResult.List.Items), filteredResult.ResourceName, names)
					}
					finalResults = append(finalResults, filteredResult)
				}
			}
		}
	}

	return finalResults, nil
}

func setCoreDefaultConfig(config *rest.Config) {
	config.GroupVersion = &corev1.SchemeGroupVersion
	config.APIPath = "/api"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()

	if config.UserAgent == "" {
		config.UserAgent = rest.DefaultKubernetesUserAgent()
	}
}

func toLegacyGrpName(str string) string {
	if str == "core" {
		return ""
	}
	return str
}

func getLegacyGrpName(str string) string {
	if str == "" {
		return "core"
	}
	return str
}

func splitParamList(nses string) []string {
	if len(nses) == 0 {
		return []string{}
	}
	if strings.Contains(nses, ",") {
		return strings.Split(nses, ",")
	}
	return strings.Split(nses, " ")
}

func filterByNames(result SearchResult, names string) SearchResult {
	if len(names) == 0 {
		return result
	}
	var filteredItems []unstructured.Unstructured
	for _, item := range result.List.Items {
		if len(names) > 0 && !strings.Contains(names, item.GetName()) {
			continue
		}
		filteredItems = append(filteredItems, item)
	}
	result.List.Items = filteredItems
	return result
}

func filterPodsByContainers(result SearchResult, containers string) SearchResult {
	if result.ListKind != "PodList" {
		return result
	}
	var filteredItems []unstructured.Unstructured
	for _, podItem := range result.List.Items {
		containerItems := getPodContainers(podItem)
		for _, containerItem := range containerItems {
			containerObj, ok := containerItem.(map[string]interface{})
			if !ok {
				logrus.Errorf("Failed to assert ustructured item (type %T) as Container", containerItem)
				continue
			}
			name, ok, err := unstructured.NestedString(containerObj, "name")
			if !ok || err != nil {
				logrus.Errorf("Failed to get container object name: %s", err)
				continue
			}
			if len(containers) > 0 && !strings.Contains(containers, name) {
				logrus.Debugf("Container %s not found in filter list %s", name, containers)
				continue
			}
			filteredItems = append(filteredItems, podItem)
		}

	}

	result.List.Items = filteredItems
	return result
}

// getPodContainers collect and return init-containers, containers, ephemeral containers from pod item
func getPodContainers(podItem unstructured.Unstructured) []interface{} {
	var result []interface{}

	initContainers, ok, err := unstructured.NestedSlice(podItem.Object, "spec", "initContainers")
	if err != nil {
		logrus.Errorf("Failed to get init-containers for pod %s: %s", podItem.GetName(), err)
	}
	if ok {
		result = append(result, initContainers...)
	}

	containers, ok, err := unstructured.NestedSlice(podItem.Object, "spec", "containers")
	if err != nil {
		logrus.Errorf("Failed to get containers for pod %s: %s", podItem.GetName(), err)
	}
	if ok {
		result = append(result, containers...)
	}

	return result
}

// getNamespaces collect all available namespaces in cluster
func getNamespaces(k8sc *Client) ([]string, error) {
	gvr := schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "namespaces",
	}
	objList, err := k8sc.Client.Resource(gvr).List(metav1.ListOptions{})

	if err != nil {
		return nil, err
	}

	names := make([]string, len(objList.Items))
	for i, obj := range objList.Items {
		names[i] = obj.GetName()
	}

	return names, nil
}
