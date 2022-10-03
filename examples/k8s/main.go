package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"time"
)

var (
	ctx, cancel = context.WithTimeout(context.Background(), time.Minute*10)
	namespace   = "test"
	kubeconfig  = filepath.Join(os.Getenv("HOME"), ".kube/config")
	deployFile  = "../../testdata/examples/deployment.yaml"
	stsFile     = "../../testdata/examples/statefulset.yaml"
	podFile     = "../../testdata/examples/pod.yaml"
	deployName  = "mydep"
	stsName     = "mysts"
	podName     = "mypod"
)

func main() {
	//Apply()
	//K8S_Create()
	//K8S_Update()
	//K8S_Apply()
	//K8S_Delete()
	//K8S_Get()
	//K8S_List()
}

func checkErr(name string, val interface{}, err error) {
	if err != nil {
		log.Printf("%s failed: %v\n", name, err)
	} else {
		log.Printf("%s success: %v\n", name, val)
	}
}
