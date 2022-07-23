package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/forbearing/k8s"
	"github.com/forbearing/k8s/pod"
)

var (
	ctx, cancel = context.WithTimeout(context.Background(), time.Minute*10)
	namespace   = "test"
	kubeconfig  = filepath.Join(os.Getenv("HOME"), ".kube/config")
	filename    = "../../testdata/examples/pod.yaml"
	filename2   = "../../testdata/nginx/nginx-pod.yaml"
	name        = "mypod"
	name2       = "nginx-pod"
	label       = "type=pod"
)

func main() {
	defer cancel()

	Pod_Tools()
	//Pod_Informer()
}

func cleanup(handler *pod.Handler) {
	handler.Delete(name)
	k8s.DeleteF(ctx, kubeconfig, filename2)
}
func checkErr(name string, val interface{}, err error) {
	if err != nil {
		log.Printf("%s failed: %v\n", name, err)
	} else {
		log.Printf("%s success: %v.\n", name, val)
	}
}