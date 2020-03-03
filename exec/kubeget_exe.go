package exec

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/vmware-tanzu/crash-diagnostics/k8s"
	"github.com/vmware-tanzu/crash-diagnostics/script"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
)

func exeKubeGet(k8sc *k8s.Client, cmd *script.KubeGetCommand) ([]k8s.SearchResult, error) {
	if k8sc == nil {
		return nil, fmt.Errorf("K8s client not initialized")
	}
	var searchResults []k8s.SearchResult

	switch cmd.What() {
	case "objects":
		logrus.Debug("KUBEGET what:objects")
		results, err := k8sc.Search(cmd.Groups(), cmd.Kinds(), cmd.Namespaces(), cmd.Versions(), cmd.Names(), cmd.Labels(), cmd.Containers())
		if err != nil {
			return nil, err
		}
		searchResults = append(searchResults, results...)
	case "logs":
		logrus.Debug("KUBEGET what:logs")
		results, err := k8sc.Search("core", "pods", cmd.Namespaces(), "", cmd.Names(), cmd.Labels(), cmd.Containers())
		if err != nil {
			return nil, err
		}
		searchResults = append(searchResults, results...)
	case "all", "*":
		logrus.Debug("KUBEGET what:all")
		results, err := k8sc.Search(cmd.Groups(), cmd.Kinds(), cmd.Namespaces(), cmd.Versions(), cmd.Names(), cmd.Labels(), cmd.Containers())
		if err != nil {
			return nil, err
		}
		searchResults = append(searchResults, results...)
	default:
		return nil, fmt.Errorf("don't know how to get: %s", cmd.What())
	}

	return searchResults, nil
}

func writeSearchResults(k8sc *k8s.Client, what string, searchResults []k8s.SearchResult, workdir string) error {
	if searchResults == nil || len(searchResults) == 0 {
		return fmt.Errorf("cannot write empty (or nil) search result")
	}

	// earch result represents a list of searched item
	// write each list in a namespaced location in working dir
	rootDir := filepath.Join(workdir, "kubeget")
	if err := os.MkdirAll(rootDir, 0744); err != nil && !os.IsExist(err) {
		return err
	}
	for _, result := range searchResults {
		resultDir := rootDir
		if result.Namespaced {
			resultDir = filepath.Join(rootDir, result.Namespace)
		}
		if err := os.MkdirAll(resultDir, 0744); err != nil && !os.IsExist(err) {
			return fmt.Errorf("failed to create search result dir: %s", err)
		}

		if err := saveResultToFile(k8sc, result, resultDir); err != nil {
			return fmt.Errorf("failed to save object: %s", err)
		}

		// print logs
		if (what == "logs" || what == "all") && result.ListKind == "PodList" {
			if len(result.List.Items) == 0 {
				continue
			}
			for _, podItem := range result.List.Items {
				logDir := filepath.Join(resultDir, podItem.GetName())
				if err := os.MkdirAll(logDir, 0744); err != nil && !os.IsExist(err) {
					return fmt.Errorf("failed to create pod log dir: %s", err)
				}

				if err := writePodLogs(k8sc, podItem, logDir); err != nil {
					logrus.Errorf("failed to save logs: pod %s: %s", podItem.GetName(), err)
					continue
				}
			}
		}

	}

	return nil
}

func saveResultToFile(k8sc *k8s.Client, result k8s.SearchResult, resultDir string) error {
	path := filepath.Join(resultDir, fmt.Sprintf("%s.json", result.ResourceName))
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	logrus.Debugf("KUBEGET: saving %s search results to: %s", result.ResourceName, path)

	if err := k8sc.JsonPrinter.PrintObj(result.List, file); err != nil {
		if wErr := writeError(err, file); wErr != nil {
			return fmt.Errorf("failed to write previous err [%s] to file: %s", err, wErr)
		}
		return err
	}
	return nil
}

func writePodLogs(k8sc *k8s.Client, podItem unstructured.Unstructured, logDir string) error {
	logrus.Debugf("KUBEGET: writing logs for pod %s", podItem.GetName())
	containers, err := getPodContainers(podItem)
	if err != nil {
		return fmt.Errorf("failed to retrieve pod containers: %s", err)
	}
	if len(containers) == 0 {
		return nil
	}

	for _, container := range containers {
		if err := writeContainerLogs(k8sc, podItem.GetNamespace(), podItem.GetName(), container, logDir); err != nil {
			return err
		}
	}

	return nil
}

func getPodContainers(podItem unstructured.Unstructured) ([]corev1.Container, error) {
	var containers []corev1.Container

	pod := new(corev1.Pod)
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(podItem.Object, &pod); err != nil {
		return nil, fmt.Errorf("error converting container objects: %s", err)
	}

	for _, c := range pod.Spec.InitContainers {
		containers = append(containers, c)
	}

	for _, c := range pod.Spec.Containers {
		containers = append(containers, c)
	}
	containers = append(containers, getPodEphemeralContainers(pod)...)
	return containers, nil
}

func getPodEphemeralContainers(pod *corev1.Pod) []corev1.Container {
	var containers []corev1.Container
	for _, ec := range pod.Spec.EphemeralContainers {
		containers = append(containers, corev1.Container(ec.EphemeralContainerCommon))
	}
	return containers
}

func writeContainerLogs(k8sc *k8s.Client, namespace string, podName string, container corev1.Container, logDir string) error {
	containerLogDir := filepath.Join(logDir, container.Name)
	if err := os.MkdirAll(containerLogDir, 0744); err != nil && !os.IsExist(err) {
		return fmt.Errorf("error creating container log dir: %s", err)
	}

	path := filepath.Join(containerLogDir, fmt.Sprintf("%s.log", container.Name))
	logrus.Debugf("Writing pod container log %s", path)

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	opts := &corev1.PodLogOptions{Container: container.Name}
	req := k8sc.CoreRest.Get().Namespace(namespace).Name(podName).Resource("pods").SubResource("log").VersionedParams(opts, scheme.ParameterCodec)
	reader, err := req.Stream()
	if err != nil {
		streamErr := fmt.Errorf("failed to create container log stream:\n%s", err)
		if wErr := writeError(streamErr, file); wErr != nil {
			return fmt.Errorf("failed to write previous err [%s] to file: %s", err, wErr)
		}
		return err
	}
	defer reader.Close()

	if _, err := io.Copy(file, reader); err != nil {
		cpErr := fmt.Errorf("failed to copy container log:\n%s", err)
		if wErr := writeError(cpErr, file); wErr != nil {
			return fmt.Errorf("failed to write previous err [%s] to file: %s", err, wErr)
		}
		return err
	}
	return nil
}

func writeError(errStr error, w io.Writer) error {
	_, err := fmt.Fprintln(w, errStr.Error())
	return err
}
