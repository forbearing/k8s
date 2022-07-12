package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/forbearing/k8s/persistentvolumeclaim"
)

var (
	ctx, cancel = context.WithTimeout(context.Background(), time.Minute*10)
	namespace   = "test"
	kubeconfig  = filepath.Join(os.Getenv("HOME"), ".kube/config")
	filename    = "../../testdata/examples/persistentvolumeclaim.yaml"
	name        = "mypvc"
	label       = "type=persistentvolumeclaim"
)

func main() {
	PersistentVolumeClaim_Tools()
}
func checkErr(name string, val interface{}, err error) {
	if err != nil {
		log.Printf("%s failed: %v\n", name, err)
	} else {
		log.Printf("%s success: %v\n", name, val)
	}
}
func cleanup(handler *persistentvolumeclaim.Handler) {
	handler.Delete(name)
}
