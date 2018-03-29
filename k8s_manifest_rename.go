package k8s_manifest_rename

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"runtime"

	io_util "github.com/bborbe/io/util"
	"github.com/golang/glog"
	"gopkg.in/yaml.v2"
	"path/filepath"
	"strings"
)

const (
	parameterPath  = "path"
	parameterWrite = "write"
)

var (
	pathPtr  = flag.String(parameterPath, "", "path")
	writePtr = flag.Bool(parameterWrite, false, "write")
)

func main() {
	defer glog.Flush()
	glog.CopyStandardLogTo("info")
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())

	writer := os.Stdout

	err := do(writer, *pathPtr, *writePtr)
	if err != nil {
		glog.Exitf("format %v failed: %v", *pathPtr, err)
	}
	glog.V(2).Infof("format %v success", *pathPtr)
}

type Kind string
type Name string

func do(writer io.Writer, path string, write bool) error {
	var err error
	glog.V(2).Infof("format %s and write %v", path, write)
	if len(path) == 0 {
		fmt.Fprintf(writer, "parameter %s missing\n", parameterPath)
		return nil
	}
	if path, err = io_util.NormalizePath(path); err != nil {
		glog.Warningf("normalize path: %s failed: %v", path, err)
		return err
	}
	source, err := ioutil.ReadFile(path)
	if err != nil {
		glog.Warningf("read file %s failed: %v", path, err)
		return err
	}
	var data struct {
		Kind     Kind `yaml:"kind"`
		Metadata struct {
			Name Name `yaml:"name"`
		} `yaml:"metadata"`
	}
	if err = yaml.Unmarshal(source, &data); err != nil {
		glog.Warningf("unmarshal %s failed: %v", path, err)
		return err
	}
	newpath := filepath.Join(filepath.Dir(path), buildName(data.Kind, data.Metadata.Name))
	if path == newpath {
		glog.V(2).Infof("skip rename")
		return nil
	}
	glog.V(1).Infof("rename file from %s to %s", path, newpath)
	return os.Rename(path, newpath)
}

func buildName(kind Kind, name Name) string {
	return fmt.Sprintf("%s-%s.yaml", name, shortenKind(kind))
}

func shortenKind(kind Kind) string {
	switch strings.ToLower(string(kind)) {
	case "ingress":
		return "ing"
	case "deployment":
		return "deploy"
	case "endpoint":
		return "ep"
	case "configmap":
		return "cm"
	case "daemonset":
		return "ds"
	case "namespace":
		return "ns"
	case "persistentvolumeclaim":
		return "pvc"
	case "persistentvolume":
		return "pv"
	case "pod":
		return "po"
	case "replicaset":
		return "rs"
	case "replicationcontroller":
		return "rc"
	case "serviceaccount":
		return "sa"
	case "service":
		return "svc"
	default:
		return strings.ToLower(string(kind))
	}
}
