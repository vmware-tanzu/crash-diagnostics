package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	spaceSep = regexp.MustCompile(`\s`)
	quoteSet = regexp.MustCompile(`[\"\']`)
)

// CommandProcessor processes cli commands, redirect output to files,
type CommandProcessor struct {
	outputPath string
	workDir    string
	cmds       []string
}

func NewCommandProcessor(cmds []string) *CommandProcessor {
	return &CommandProcessor{workDir: "/tmp/flare", outputPath: "./flareout.tar.gz", cmds: cmds}
}

func (p *CommandProcessor) Process() error {
	if len(p.cmds) == 0 {
		return nil
	}

	var filesWritten []string
	defer func() {
		for _, f := range filesWritten {
			os.RemoveAll(f)
		}
	}()

	// setup file output dir
	if err := os.MkdirAll(p.workDir, 0744); err != nil && !os.IsExist(err) {
		return err
	}

	for _, cmd := range p.cmds {
		parts := p.splitCommand(cmd)
		exec := NewCommandExecutor(parts[0], parts[1:]...)
		reader, err := exec.Execute()
		if err != nil {
			return err
		}

		fileName := fmt.Sprintf("%s.txt", p.genFileName(cmd))
		filePath := filepath.Join(p.workDir, fileName)
		if err := p.redirectToFile(reader, filePath); err != nil {
			return err
		}
		log.Printf("wrote file %s\n", filePath)
		filesWritten = append(filesWritten, filePath)
	}

	// tar directory
	if err := p.tar(p.outputPath, p.workDir); err != nil {
		return err
	}

	return nil
}

func (p *CommandProcessor) splitCommand(cmd string) []string {
	return spaceSep.Split(cmd, -1)
}

func (p *CommandProcessor) redirectToFile(source io.Reader, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err := io.Copy(file, source); err != nil {
		return err
	}
	return nil
}

func (p *CommandProcessor) genFileName(cmd string) string {
	str := quoteSet.ReplaceAllString(cmd, "")
	return spaceSep.ReplaceAllString(str, "_")
}

func (p *CommandProcessor) tar(tarPath, sourcePath string) error {
	log.Printf("tarring source %s to %s\n", sourcePath, tarPath)

	if !filepath.IsAbs(tarPath) {
		return fmt.Errorf("tar path must be absolute: %s", tarPath)
	}

	if !filepath.IsAbs(sourcePath) {
		return fmt.Errorf("tar source path must be absolute:%s", sourcePath)
	}

	_, err := os.Stat(sourcePath)
	if err != nil {
		return err
	}

	tarFile, err := os.Create(tarPath)
	if err != nil {
		return err
	}
	defer tarFile.Close()
	zipper := gzip.NewWriter(tarFile)
	defer zipper.Close()
	tarrer := tar.NewWriter(zipper)
	defer tarrer.Close()

	prefix := "/"
	if strings.HasPrefix(sourcePath, "/tmp") {
		prefix = "/tmp"
	}

	return filepath.Walk(sourcePath, func(file string, finfo os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relFilePath, err := filepath.Rel(prefix, file)
		if err != nil {
			return err
		}

		header, err := tar.FileInfoHeader(finfo, finfo.Name())
		if err != nil {
			return err
		}
		header.Name = relFilePath
		if err := tarrer.WriteHeader(header); err != nil {
			return err
		}
		if finfo.IsDir() {
			return nil
		}

		// add file to tar
		srcFile, err := os.Open(file)
		if err != nil {
			return err
		}
		defer srcFile.Close()
		_, err = io.Copy(tarrer, srcFile)
		if err != nil {
			return err
		}
		return nil
	})
}
