package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/forbearing/k8s"
	"github.com/forbearing/k8s/daemonset"
)

var (
	ctx, cancel = context.WithTimeout(context.Background(), time.Minute*10)
	namespace   = "test"
	kubeconfig  = filepath.Join(os.Getenv("HOME"), ".kube/config")
	filename    = "../../testdata/examples/daemonset.yaml"
	filename2   = "../../testdata/nginx/nginx-ds.yaml"
	name        = "myds"
	name2       = "nginx-ds"
	label       = "type=daemonset"
)

func main() {
	DaemonSet_Tools()
}

func checkErr(name string, val interface{}, err error) {
	if err != nil {
		log.Printf("%s failed: %v\n", name, err)
	} else {
		log.Printf("%s success: %v\n", name, val)
	}
}
func cleanup(handler *daemonset.Handler) {
	handler.Delete(name)
	k8s.DeleteF(ctx, kubeconfig, filename2)
}
