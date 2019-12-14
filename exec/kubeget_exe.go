package exec

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/vmware-tanzu/crash-diagnostics/k8s"
	"github.com/vmware-tanzu/crash-diagnostics/script"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
)

func exeKubeGet(k8sc *k8s.Client, cmd *script.KubeGetCommand, workdir string) error {
	if k8sc == nil {
		return fmt.Errorf("K8s client not initialized")
	}
	var objects []runtime.Object

	switch cmd.What() {
	case "objects":
		logrus.Debug("KUBEGET what:objects")
		objs, err := search(k8sc, cmd.Groups(), cmd.Kinds(), cmd.Namespaces(), cmd.Versions(), cmd.Names(), cmd.Labels(), cmd.Containers())
		if err != nil {
			return err
		}
		objects = append(objects, objs...)
	case "logs":
		logrus.Debug("KUBEGET what:logs")
		objs, err := search(k8sc, "core", "pods", cmd.Namespaces(), "", cmd.Names(), cmd.Labels(), cmd.Containers())
		if err != nil {
			return err
		}
		objects = append(objects, objs...)
	case "all", "*":
		logrus.Debug("KUBEGET what:all")
		objs, err := search(k8sc, cmd.Groups(), cmd.Kinds(), cmd.Namespaces(), cmd.Versions(), cmd.Names(), cmd.Labels(), cmd.Containers())
		if err != nil {
			return err
		}
		objects = append(objects, objs...)

	default:
		return fmt.Errorf("don't know how to get: %s", cmd.What())
	}

	// print object lists
	for _, obj := range objects {
		objList, ok := obj.(*unstructured.UnstructuredList)
		if !ok {
			logrus.Errorf("KUBEGET: unexpected object type for %T", obj)
			continue
		}
		if err := writeObjectList(k8sc, cmd.What(), objList, workdir); err != nil {
			logrus.Errorf("KUBEGET: %s", err)
			continue
		}
	}

	return nil
}

// search does a drill-down search from group, version, resourceList, to resources.  The following rules are applied
// 1) Legacy core group (api/v1) can be specified as "core"
// 2) All specified search params will use AND operator for match (i.e. groups=core and kinds=pods and versions=v1 and ... etc)
// 3) kinds will match resource.Kind or resource.Name
// 4) All search params are passed as comma- or space-separated sets that are matched using OR (i.e. kinds=pods services
//    will match resouces of type pods or services)
func search(k8sc *k8s.Client, groups, kinds, namespaces, versions, names, labels, containers string) ([]runtime.Object, error) {
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

	var runtimeObjs []runtime.Object
	logrus.Debugf("Searching in %d groups", len(grpList.Groups))
	for _, grp := range grpList.Groups {
		// filter by group
		grpName := strings.TrimSpace(grp.Name)
		grpName = getLegacyGrpName(grpName)
		if len(groups) > 0 && !strings.Contains(groups, strings.ToLower(grpName)) {
			continue
		}
		// filter by group version
		for _, discoGV := range grp.Versions {
			if len(versions) > 0 && !strings.Contains(versions, strings.ToLower(discoGV.Version)) {
				continue
			}

			// adjust version for legacy group
			grpVersion := discoGV.GroupVersion
			if grpName == k8s.LegacyGroupName {
				grpVersion = discoGV.Version
			}

			logrus.Debugf("Searching resources in GroupVersion %s", discoGV.GroupVersion)
			resources, err := k8sc.Disco.ServerResourcesForGroupVersion(grpVersion)
			if err != nil {
				logrus.Errorf("KUBEGET failed to get resources for %s: %s", discoGV.GroupVersion, err)
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

				// retrieve API objects based on GroupVersionResource and
				// filter by namespaces and names
				listOptions := metav1.ListOptions{
					LabelSelector: labels,
				}

				var resList []*unstructured.UnstructuredList
				if len(namespaces) > 0 {
					logrus.Debugf("Searching %s in namespace [%s]", res.Name, namespaces)
					for _, ns := range splitParamList(namespaces) {
						list, err := k8sc.Client.Resource(gvr).Namespace(ns).List(listOptions)
						if err != nil {
							logrus.Debugf(
								"WARN: KUBEGET failed to get resource list for %s/%s [%v]: %s",
								discoGV.GroupVersion, ns, listOptions.LabelSelector, err,
							)
							continue
						}
						resList = append(resList, list)
					}
				} else {
					logrus.Debugf("Searching %s", res.Name)
					list, err := k8sc.Client.Resource(gvr).List(listOptions)
					if err != nil {
						logrus.Debugf(
							"WARN: KUBEGET ObjSearch failed to get resource list for %s [%v]: %s",
							discoGV.GroupVersion, listOptions.LabelSelector, err,
						)
						continue
					}
					resList = append(resList, list)
				}

				for _, list := range resList {
					if list != nil && len(list.Items) > 0 {
						logrus.Debugf("Found %d %s", len(list.Items), res.Name)
						objs := filterByNames(list, names)
						// if podlist, apply container filter
						if list.GetKind() == "PodList" {
							logrus.Debugf("Filtering PodList by containers: %s", containers)
							objs = filterPodsByContainers(objs, containers)
						}
						runtimeObjs = append(runtimeObjs, objs)
					}
				}
			}
		}
	}

	return runtimeObjs, nil
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
	if strings.Contains(nses, ",") {
		return strings.Split(nses, ",")
	}
	return strings.Split(nses, " ")
}

func strSliceContains(values []string, val string) bool {
	if len(values) == 0 {
		return false
	}

	for _, v := range values {
		if strings.EqualFold(strings.TrimSpace(v), strings.TrimSpace(val)) {
			return true
		}
	}
	return false
}

func filterByNames(list *unstructured.UnstructuredList, names string) *unstructured.UnstructuredList {
	if len(names) == 0 {
		return list
	}
	var filteredItems []unstructured.Unstructured
	for _, item := range list.Items {
		if len(names) > 0 && !strings.Contains(names, item.GetName()) {
			continue
		}
		filteredItems = append(filteredItems, item)
	}
	list.Items = filteredItems
	return list
}

func filterPodsByContainers(list *unstructured.UnstructuredList, containers string) *unstructured.UnstructuredList {
	if list.GetKind() != "PodList" {
		return list
	}
	var filteredItems []unstructured.Unstructured
	for _, podItem := range list.Items {
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

	list.Items = filteredItems
	return list
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

func writeObjectList(k8sc *k8s.Client, what string, objList *unstructured.UnstructuredList, workdir string) error {
	if objList == nil {
		return fmt.Errorf("cannot write nil object list")
	}

	writer := os.Stdout
	if workdir != "stdout" {
		kind := objList.GetKind()
		path := filepath.Join(workdir, fmt.Sprintf("kubeget-%s.json", kind))
		file, err := os.Create(path)
		if err != nil {
			return err
		}
		defer file.Close()
		writer = file
		logrus.Debugf("KUBEGET: writing objects to %s", path)
	}

	if err := k8sc.JsonPrinter.PrintObj(objList, writer); err != nil {
		return err
	}

	// if obj is PodList, write logs for items in list
	if what == "logs" || what == "all" {
		if objList.GetKind() != "PodList" {
			return nil
		}
		for _, podItem := range objList.Items {
			if err := writePodLogs(k8sc, podItem, workdir); err != nil {
				logrus.Errorf("Failed to write logs for pod %s: %s", podItem.GetName(), err)
				continue
			}
		}
	}
	return nil
}

func writePodLogs(k8sc *k8s.Client, podItem unstructured.Unstructured, workdir string) error {
	ns := podItem.GetNamespace()
	name := podItem.GetName()

	writer := os.Stdout
	if workdir != "stdout" {
		path := filepath.Join(workdir, fmt.Sprintf("kubeget-podlog-%s-%s.txt", ns, name))
		logFile, err := os.Create(path)
		if err != nil {
			return err
		}
		defer logFile.Close()
		writer = logFile
		logrus.Debugf("KUBEGET: writing pod logs to %s", path)
	}

	req := k8sc.CoreRest.Get().Namespace(ns).Name(name).Resource("pods").SubResource("log").VersionedParams(&corev1.PodLogOptions{}, scheme.ParameterCodec)
	reader, err := req.Stream()
	if err != nil {
		return err
	}
	defer reader.Close()

	writeLogHeader(podItem, writer)
	if _, err := io.Copy(writer, reader); err != nil {
		return err
	}

	return nil
}

func writeLogHeader(podItem unstructured.Unstructured, w io.Writer) {
	fmt.Fprintln(w, "----------------------------------------------------------------")
	fmt.Fprintf(w, "Log pod %s/%s\n", podItem.GetNamespace(), podItem.GetName())
	fmt.Fprintln(w, "----------------------------------------------------------------")
}
