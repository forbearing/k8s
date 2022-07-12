package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/forbearing/k8s/clusterrolebinding"
)

var (
	ctx, cancel = context.WithTimeout(context.Background(), time.Minute*10)
	namespace   = "test"
	kubeconfig  = filepath.Join(os.Getenv("HOME"), ".kube/config")
	filename    = "../../testdata/examples/clusterrolebinding.yaml"
	name        = "myclusterrolebinding"
	label       = "type=clusterrolebinding"
)

func main() {
	RoleBinding_Tools()
}

func checkErr(name string, val interface{}, err error) {
	if err != nil {
		log.Printf("%s failed: %v\n", name, err)
	} else {
		log.Printf("%s success: %v\n", name, val)
	}
}
func cleanup(handler *clusterrolebinding.Handler) {
	handler.Delete(name)
}
