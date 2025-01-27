package k8s

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/cli-runtime/pkg/printers"
)

type ObjectWriter struct {
	writeDir   string
	printer    printers.ResourcePrinter
	singleFile bool
}

func (w *ObjectWriter) Write(result SearchResult) (string, error) {
	// namespaced on group and version to avoid overwrites
	grp := func() string {
		if result.GroupVersionResource.Group == "" {
			return "core"
		}
		return result.GroupVersionResource.Group
	}()

	w.writeDir = filepath.Join(w.writeDir, fmt.Sprintf("%s_%s", grp, result.GroupVersionResource.Version))

	// add resource namespace if needed
	if result.Namespaced {
		w.writeDir = filepath.Join(w.writeDir, result.Namespace)
	}

	now := time.Now().Format("2006-01-02T15-04-05Z.0000")
	var extension string
	if _, ok := w.printer.(*printers.JSONPrinter); ok {
		extension = "json"
	} else {
		extension = "yaml"
	}

	if w.singleFile {
		if err := os.MkdirAll(w.writeDir, 0744); err != nil && !os.IsExist(err) {
			return "", fmt.Errorf("failed to create search result dir: %s", err)
		}
		path := filepath.Join(w.writeDir, fmt.Sprintf("%s-%s.%s", result.ResourceName, now, extension))
		return w.writeDir, w.writeFile(result.List, path)
	} else {
		w.writeDir = filepath.Join(w.writeDir, result.ResourceName)
		if err := os.MkdirAll(w.writeDir, 0744); err != nil && !os.IsExist(err) {
			return "", fmt.Errorf("failed to create search result dir: %s", err)
		}

		for i := range result.List.Items {
			u := &result.List.Items[i]
			path := filepath.Join(w.writeDir, fmt.Sprintf("%s-%s.%s", u.GetName(), now, extension))
			if err := w.writeFile(u, path); err != nil {
				return "", err
			}
		}
	}

	return w.writeDir, nil
}

func (w *ObjectWriter) writeFile(o runtime.Object, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	logrus.Debugf("objectWriter: saving %s search results to: %s", o.GetObjectKind(), path)

	if err := w.printer.PrintObj(o, file); err != nil {
		if wErr := writeError(err, file); wErr != nil {
			return fmt.Errorf("objectWriter: failed to write previous err [%s] to file: %s", err, wErr)
		}
		return err
	}

	return nil
}
