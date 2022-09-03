package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/forbearing/k8s/types"
)

var (
	ctx, cancel = context.WithTimeout(context.Background(), time.Minute*10)
	namespace   = "test"
	kubeconfig  = filepath.Join(os.Getenv("HOME"), ".kube/config")
	deployFile  = "../../testdata/examples/deployment.yaml"
	deployName  = "mydep"
	podFile     = "../../testdata/examples/pod.yaml"
	podName     = "mypod"
)

func main() {
	//Alias()
	Apply()
}

func checkErr(name string, val interface{}, err error) {
	if err != nil {
		log.Printf("%s failed: %v\n", name, err)
	} else {
		log.Printf("%s success: %v\n", name, val)
	}
}

// cleanup will delete or prune created deployments.
func cleanup(handler types.Deleter) {
	fmt.Println(handler.Delete(deployName))
	fmt.Println(handler.Delete(podName))
}
