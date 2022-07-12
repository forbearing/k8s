package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"time"

	k8snamespace "github.com/forbearing/k8s/namespace"
)

var (
	ctx, cancel = context.WithTimeout(context.Background(), time.Minute*10)
	kubeconfig  = filepath.Join(os.Getenv("HOME"), ".kube/config")
	filename    = "../../testdata/examples/namespace.yaml"
	name        = "test1"
)

func main() {
	Namespace_Tools()
}

func checkErr(name string, err error) {
	if err != nil {
		log.Printf("%s failed: %v\n", name, err)
	} else {
		log.Printf("%s success.\n", name)
	}
}
func cleanup(handler *k8snamespace.Handler) {
	handler.Delete(name)
}
