package main_test

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var pathToBinary string
var session *gexec.Session

var _ = BeforeSuite(func() {
	var err error
	pathToBinary, err = gexec.Build("github.com/bborbe/k8s-manifest-rename")
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterEach(func() {
	session.Interrupt()
	Eventually(session).Should(gexec.Exit())
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})

func TestCheck(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "K8s Manifest Rename Suite")
}

var _ = Describe("K8s Rename Tool", func() {
	var tmpdir string
	var removefiles []string
	var err error
	BeforeEach(func() {
		tmpdir, err = ioutil.TempDir("", "")
		Expect(err).To(BeNil())
	})
	AfterEach(func() {
		for _, file := range removefiles {
			os.Remove(file)
		}
	})
	It("exit with exitcode 1", func() {
		session, err = gexec.Start(exec.Command(pathToBinary), GinkgoWriter, GinkgoWriter)
		Expect(err).To(BeNil())
		session.Wait(100 * time.Millisecond)
		Expect(session.ExitCode()).To(Equal(1))
	})
	Context("path is valid", func() {
		var file string
		var args []string
		BeforeEach(func() {
			file = path.Join(tmpdir, expectName)
			removefiles = append(removefiles, file)
			args = append(args, "-path", file)
			args = append(args, "-write")
			args = append(args, "-v", "4", "-logtostderr")
			err = ioutil.WriteFile(file, []byte(validContent), 0755)
			Expect(err).To(BeNil())
		})
		It("exit with exitcode 0", func() {
			session, err = gexec.Start(exec.Command(pathToBinary, args...), GinkgoWriter, GinkgoWriter)
			Expect(err).To(BeNil())
			session.Wait(100 * time.Millisecond)
			Expect(session.ExitCode()).To(Equal(0))
		})
		It("file still exists", func() {
			session, err = gexec.Start(exec.Command(pathToBinary, args...), GinkgoWriter, GinkgoWriter)
			Expect(err).To(BeNil())
			session.Wait(100 * time.Millisecond)
			_, err := os.Stat(file)
			Expect(err).To(BeNil())
		})
	})
	Context("path is invalid", func() {
		var file string
		var args []string
		BeforeEach(func() {
			file = path.Join(tmpdir, "hello-ingress.yaml")
			removefiles = append(removefiles, file, path.Join(tmpdir, expectName))
			args = append(args, "-path", file)
			args = append(args, "-write")
			args = append(args, "-v", "4", "-logtostderr")
			err = ioutil.WriteFile(file, []byte(validContent), 0755)
			Expect(err).To(BeNil())
		})
		It("exit with exitcode 0", func() {
			session, err = gexec.Start(exec.Command(pathToBinary, args...), GinkgoWriter, GinkgoWriter)
			Expect(err).To(BeNil())
			session.Wait(100 * time.Millisecond)
			Expect(session.ExitCode()).To(Equal(0))
		})
		It("orginal file does not exists anymore", func() {
			session, err = gexec.Start(exec.Command(pathToBinary, args...), GinkgoWriter, GinkgoWriter)
			Expect(err).To(BeNil())
			session.Wait(100 * time.Millisecond)
			_, err = os.Stat(file)
			Expect(err).NotTo(BeNil())
		})
		It("file exists under expected name", func() {
			session, err = gexec.Start(exec.Command(pathToBinary, args...), GinkgoWriter, GinkgoWriter)
			Expect(err).To(BeNil())
			session.Wait(100 * time.Millisecond)
			_, err = os.Stat(path.Join(tmpdir, expectName))
			Expect(err).To(BeNil())
		})
	})
})

var expectName = "hello-ing.yaml"
var validContent = `apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: hello
  namespace: world`
