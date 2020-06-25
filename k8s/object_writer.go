package k8s

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"k8s.io/cli-runtime/pkg/printers"
)

type ObjectWriter struct {
	writeDir string
}

func (w ObjectWriter) Write(result SearchResult) (string, error) {
	resultDir := w.writeDir
	if result.Namespaced {
		resultDir = filepath.Join(w.writeDir, result.Namespace)
	}
	if err := os.MkdirAll(resultDir, 0744); err != nil && !os.IsExist(err) {
		return "", fmt.Errorf("failed to create search result dir: %s", err)
	}

	path := filepath.Join(resultDir, fmt.Sprintf("%s.json", result.ResourceName))
	file, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	logrus.Debugf("kube_capture(): saving %s search results to: %s", result.ResourceName, path)

	printer := new(printers.JSONPrinter)
	if err := printer.PrintObj(result.List, file); err != nil {
		if wErr := writeError(err, file); wErr != nil {
			return "", fmt.Errorf("failed to write previous err [%s] to file: %s", err, wErr)
		}
		return "", err
	}
	return resultDir, nil
}
