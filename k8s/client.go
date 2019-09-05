package k8s

import (
	"fmt"
	"io"
	"os"

	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/cli-runtime/pkg/printers"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	namespace           = "kube-system"
	namespacesResource  = schema.GroupVersionResource{Group: "", Version: "v1", Resource: "namespaces"}
	nodesResource       = schema.GroupVersionResource{Group: "", Version: "v1", Resource: "nodes"}
	eventsResource      = schema.GroupVersionResource{Group: "", Version: "v1", Resource: "events"}
	rcsResource         = schema.GroupVersionResource{Group: "", Version: "v1", Resource: "replicationcontrollers"}
	servicesResource    = schema.GroupVersionResource{Group: "", Version: "v1", Resource: "services"}
	dssResource         = schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "daemonsets"}
	deploymentsResource = schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}
	rssResource         = schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "replicasets"}
	podsResource        = schema.GroupVersionResource{Group: "", Version: "v1", Resource: "pods"}
)

func GetClient(kubeconfig string) (dynamic.Interface, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}
	return dynamic.NewForConfig(config)
}

// dumpClusterInfo attempts to retrieve cluster information from API
// When error is encountered, it simply logs it.
func DumpClusterInfo(client dynamic.Interface, path string) error {
	namespaces, err := getNamespaces(client)
	if err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	logrus.Debugf("Dumping cluster info in %s", path)

	pr := new(printers.JSONPrinter)
	var objs []runtime.Object

	// dump node info
	if nodes, err := client.Resource(nodesResource).List(metav1.ListOptions{}); err != nil {
		dumpError(file, nodesResource.String(), err)
	} else {
		objs = append(objs, nodes)
	}

	// get namespaced objs
	for _, ns := range namespaces {
		if events, err := client.Resource(eventsResource).Namespace(ns).List(metav1.ListOptions{}); err != nil {
			dumpError(file, eventsResource.String(), err)
		} else {
			objs = append(objs, events)
		}

		if rcs, err := client.Resource(rcsResource).Namespace(ns).List(metav1.ListOptions{}); err != nil {
			dumpError(file, rcsResource.String(), err)
		} else {
			objs = append(objs, rcs)
		}

		if services, err := client.Resource(servicesResource).Namespace(ns).List(metav1.ListOptions{}); err != nil {
			dumpError(file, servicesResource.String(), err)
		} else {
			objs = append(objs, services)
		}

		if daemonsets, err := client.Resource(dssResource).Namespace(ns).List(metav1.ListOptions{}); err != nil {
			dumpError(file, dssResource.String(), err)
		} else {
			objs = append(objs, daemonsets)
		}

		if deployments, err := client.Resource(deploymentsResource).Namespace(ns).List(metav1.ListOptions{}); err != nil {
			dumpError(file, deploymentsResource.String(), err)
		} else {
			objs = append(objs, deployments)
		}

		if replicas, err := client.Resource(rssResource).Namespace(ns).List(metav1.ListOptions{}); err != nil {
			dumpError(file, rssResource.String(), err)
		} else {
			objs = append(objs, replicas)
		}

		if pods, err := client.Resource(podsResource).Namespace(ns).List(metav1.ListOptions{}); err != nil {
			dumpError(file, podsResource.String(), err)
		} else {
			objs = append(objs, pods)

			// TODO capture container logs for each pod.
		}
	}

	return printToFile(pr, file, objs...)
}

func getNamespaces(client dynamic.Interface) ([]string, error) {
	var result []string
	nss, err := client.Resource(namespacesResource).List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, nsObj := range nss.Items {
		result = append(result, nsObj.GetName())
	}
	return result, nil
}

func printToFile(pr *printers.JSONPrinter, f *os.File, objs ...runtime.Object) error {
	for _, obj := range objs {
		if obj != nil {
			if err := pr.PrintObj(obj, f); err != nil {
				dumpError(f, obj.(*unstructured.Unstructured).GetName(), err)
			}
		}
	}
	return nil
}

func dumpError(writer io.Writer, resource string, err error) {
	errMsg := fmt.Sprintf("Error during cluster-info collection for %s: %s\n", resource, err)
	logrus.Errorf(errMsg)
	fmt.Fprint(writer, errMsg)
}
