package main

import (
	"flag"
	"fmt"
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
	parameterPath     = "path"
	parameterWrite    = "write"
	parameterValidate = "validate"
)

var (
	pathPtr     = flag.String(parameterPath, "", "path")
	writePtr    = flag.Bool(parameterWrite, false, "write")
	validatePtr = flag.Bool(parameterValidate, false, "validate content is already formated")
)

type Kind string
type Name string

func main() {
	defer glog.Flush()
	glog.CopyStandardLogTo("info")
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())

	if len(*pathPtr) == 0 {
		fmt.Fprintf(os.Stderr, "parameter %s missing\n", parameterPath)
		os.Exit(1)
	}
	path, err := io_util.NormalizePath(*pathPtr)
	if err != nil {
		glog.Exitf("normalize path: %s failed: %v", path, err)
	}
	source, err := ioutil.ReadFile(path)
	if err != nil {
		glog.Exitf("read file %s failed: %v", path, err)
	}
	var data struct {
		Kind Kind `yaml:"kind"`
		Metadata struct {
			Name Name `yaml:"name"`
		} `yaml:"metadata"`
	}
	if err = yaml.Unmarshal(source, &data); err != nil {
		glog.Exitf("unmarshal %s failed: %v", path, err)
	}
	newpath := filepath.Join(filepath.Dir(path), buildName(data.Kind, data.Metadata.Name))
	if *writePtr {
		if path != newpath {
			if err := os.Rename(path, newpath); err != nil {
				glog.Exitf("rename failed: %v", err)
			}
		}
		return
	}
	if *validatePtr {
		if path != newpath {
			fmt.Fprintf(os.Stderr, "path invalid! %s != %s missing\n", path, newpath)
			os.Exit(1)
		}
	}
	fmt.Printf("path should be %s", newpath)
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
