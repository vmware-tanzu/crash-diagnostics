package k8s

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"k8s.io/cli-runtime/pkg/printers"
	"k8s.io/client-go/rest"
)

type ResultWriter struct {
	workdir   string
	writeLogs bool
	restApi   rest.Interface
	printer   printers.ResourcePrinter
}

func NewResultWriter(workdir, what, outputFormat string, restApi rest.Interface) (*ResultWriter, error) {
	var err error
	workdir = filepath.Join(workdir, BaseDirname)
	if err := os.MkdirAll(workdir, 0744); err != nil && !os.IsExist(err) {
		return nil, err
	}

	writeLogs := what == "logs" || what == "all"
	var printer printers.ResourcePrinter
	if outputFormat == "" || outputFormat == "json" {
		printer = &printers.JSONPrinter{}
	} else if outputFormat == "yaml" {
		printer = &printers.YAMLPrinter{}
	} else {
		return nil, errors.Errorf("unsupported output format: %s", outputFormat)
	}
	return &ResultWriter{
		workdir:   workdir,
		printer:   printer,
		writeLogs: writeLogs,
		restApi:   restApi,
	}, err
}

func (w *ResultWriter) GetResultDir() string {
	return w.workdir
}

func (w *ResultWriter) Write(ctx context.Context, searchResults []SearchResult) error {
	if len(searchResults) == 0 {
		return fmt.Errorf("cannot write empty (or nil) search result")
	}

	// each result represents a list of searched item
	// write each list in a namespaced location in working dir
	for _, result := range searchResults {
		objWriter := ObjectWriter{
			writeDir: w.workdir,
			printer:  w.printer,
		}
		writeDir, err := objWriter.Write(result)
		if err != nil {
			return err
		}

		if w.writeLogs && result.ListKind == "PodList" {
			if len(result.List.Items) == 0 {
				continue
			}
			for _, podItem := range result.List.Items {
				logDir := filepath.Join(writeDir, podItem.GetName())
				if err := os.MkdirAll(logDir, 0744); err != nil && !os.IsExist(err) {
					return fmt.Errorf("failed to create pod log dir: %s", err)
				}

				containers, err := GetContainers(podItem)
				if err != nil {
					return err
				}
				for _, containerLogger := range containers {
					reader, err := containerLogger.Fetch(ctx, w.restApi)
					if err != nil {
						return err
					}
					err = containerLogger.Write(reader, logDir)
					if err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}
