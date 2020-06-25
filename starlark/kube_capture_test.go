package starlark

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go.starlark.net/starlarkstruct"

	"github.com/sirupsen/logrus"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	testcrashd "github.com/vmware-tanzu/crash-diagnostics/testing"
)

var _ = Describe("kube_capture", func() {

	var (
		k8sconfig string
		kind      *testcrashd.KindCluster
		waitTime  = time.Second * 11
		workdir   string

		executor *Executor
		err      error
	)

	BeforeSuite(func() {
		clusterName := "crashd-test-kubecapture"
		tmpFile, err := ioutil.TempFile(os.TempDir(), clusterName)
		Expect(err).NotTo(HaveOccurred())
		k8sconfig = tmpFile.Name()

		// create kind cluster
		kind = testcrashd.NewKindCluster("../testing/kind-cluster-docker.yaml", clusterName)
		err = kind.Create()
		Expect(err).NotTo(HaveOccurred())

		err = kind.MakeKubeConfigFile(k8sconfig)
		Expect(err).NotTo(HaveOccurred())

		logrus.Infof("Sleeping %v ... waiting for pods", waitTime)
		time.Sleep(waitTime)
	})

	AfterSuite(func() {
		kind.Destroy()
		os.RemoveAll(k8sconfig)
	})

	execSetup := func(crashdScript string) {
		executor = New()
		err = executor.Exec("test.kube.capture", strings.NewReader(crashdScript))
		Expect(err).To(BeNil())
	}

	BeforeEach(func() {
		workdir, err = ioutil.TempDir(os.TempDir(), "test")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		os.RemoveAll(workdir)
	})

	It("creates a directory and files for namespaced objects", func() {
		crashdScript := fmt.Sprintf(`
crashd_config(workdir="%s")
kube_config(path="%s")
kube_data = kube_capture(what="objects", groups="core", kinds="services", namespaces=["default", "kube-system"])
		`, workdir, k8sconfig)
		execSetup(crashdScript)
		Expect(executor.result.Has("kube_data")).NotTo(BeNil())

		data := executor.result["kube_data"]
		Expect(data).NotTo(BeNil())

		captureData, _ := data.(*starlarkstruct.Struct)
		Expect(captureData.AttrNames()).To(HaveLen(2))

		errVal, err := captureData.Attr("error")
		Expect(err).NotTo(HaveOccurred())
		Expect(trimQuotes(errVal.String())).To(BeEmpty())

		fileVal, err := captureData.Attr("file")
		Expect(err).NotTo(HaveOccurred())
		Expect(trimQuotes(fileVal.String())).To(BeADirectory())

		kubeCaptureDir := trimQuotes(fileVal.String())
		Expect(filepath.Join(kubeCaptureDir, "default", "services.json")).To(BeARegularFile())
		Expect(filepath.Join(kubeCaptureDir, "kube-system", "services.json")).To(BeARegularFile())
	})

	It("creates a directory and files for non-namespaced objects", func() {
		crashdScript := fmt.Sprintf(`
crashd_config(workdir="%s")
kube_config(path="%s")
kube_data = kube_capture(what="objects", groups="core", kinds="nodes")
		`, workdir, k8sconfig)
		execSetup(crashdScript)
		Expect(executor.result.Has("kube_data")).NotTo(BeNil())

		data := executor.result["kube_data"]
		Expect(data).NotTo(BeNil())

		captureData, _ := data.(*starlarkstruct.Struct)
		Expect(captureData.AttrNames()).To(HaveLen(2))

		errVal, err := captureData.Attr("error")
		Expect(err).NotTo(HaveOccurred())
		Expect(trimQuotes(errVal.String())).To(BeEmpty())

		fileVal, err := captureData.Attr("file")
		Expect(err).NotTo(HaveOccurred())
		Expect(trimQuotes(fileVal.String())).To(BeADirectory())

		kubeCaptureDir := trimQuotes(fileVal.String())
		Expect(filepath.Join(kubeCaptureDir, "nodes.json")).To(BeARegularFile())
	})

	It("creates a directory and log files for all objects in a namespace", func() {
		crashdScript := fmt.Sprintf(`
crashd_config(workdir="%s")
kube_config(path="%s")
kube_data = kube_capture(what="logs", namespaces="kube-system")
		`, workdir, k8sconfig)
		execSetup(crashdScript)
		Expect(executor.result.Has("kube_data")).NotTo(BeNil())

		data := executor.result["kube_data"]
		Expect(data).NotTo(BeNil())

		captureData, _ := data.(*starlarkstruct.Struct)
		Expect(captureData.AttrNames()).To(HaveLen(2))

		errVal, err := captureData.Attr("error")
		Expect(err).NotTo(HaveOccurred())
		Expect(trimQuotes(errVal.String())).To(BeEmpty())

		fileVal, err := captureData.Attr("file")
		Expect(err).NotTo(HaveOccurred())
		Expect(trimQuotes(fileVal.String())).To(BeADirectory())

		kubeCaptureDir := trimQuotes(fileVal.String())
		Expect(filepath.Join(kubeCaptureDir, "kube-system")).To(BeADirectory())

		files, err := ioutil.ReadDir(filepath.Join(kubeCaptureDir, "kube-system"))
		Expect(err).NotTo(HaveOccurred())
		Expect(len(files)).NotTo(BeNumerically("<", 3))
	})

	It("creates a log file for specific container in a namespace", func() {
		crashdScript := fmt.Sprintf(`
crashd_config(workdir="%s")
kube_config(path="%s")
kube_data = kube_capture(what="logs", namespaces="kube-system", containers=["etcd"])
		`, workdir, k8sconfig)
		execSetup(crashdScript)
		Expect(executor.result.Has("kube_data")).NotTo(BeNil())

		data := executor.result["kube_data"]
		Expect(data).NotTo(BeNil())

		captureData, _ := data.(*starlarkstruct.Struct)
		Expect(captureData.AttrNames()).To(HaveLen(2))

		errVal, err := captureData.Attr("error")
		Expect(err).NotTo(HaveOccurred())
		Expect(trimQuotes(errVal.String())).To(BeEmpty())

		fileVal, err := captureData.Attr("file")
		Expect(err).NotTo(HaveOccurred())
		Expect(trimQuotes(fileVal.String())).To(BeADirectory())

		kubeCaptureDir := trimQuotes(fileVal.String())
		Expect(filepath.Join(kubeCaptureDir, "kube-system")).To(BeADirectory())

		files, err := ioutil.ReadDir(filepath.Join(kubeCaptureDir, "kube-system"))
		Expect(err).NotTo(HaveOccurred())
		Expect(files).NotTo(HaveLen(0))
	})

})
