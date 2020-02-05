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

func exeKubeGet(k8sc *k8s.Client, cmd *script.KubeGetCommand, workdir string) error {
	if k8sc == nil {
		return fmt.Errorf("K8s client not initialized")
	}
	var objects []runtime.Object

	switch cmd.What() {
	case "objects":
		logrus.Debug("KUBEGET what:objects")
		objs, err := k8sc.Search(cmd.Groups(), cmd.Kinds(), cmd.Namespaces(), cmd.Versions(), cmd.Names(), cmd.Labels(), cmd.Containers())
		if err != nil {
			return err
		}
		objects = append(objects, objs...)
	case "logs":
		logrus.Debug("KUBEGET what:logs")
		objs, err := k8sc.Search("core", "pods", cmd.Namespaces(), "", cmd.Names(), cmd.Labels(), cmd.Containers())
		if err != nil {
			return err
		}
		objects = append(objects, objs...)
	case "all", "*":
		logrus.Debug("KUBEGET what:all")
		objs, err := k8sc.Search(cmd.Groups(), cmd.Kinds(), cmd.Namespaces(), cmd.Versions(), cmd.Names(), cmd.Labels(), cmd.Containers())
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
