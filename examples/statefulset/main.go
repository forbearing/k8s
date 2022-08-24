package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/forbearing/k8s"
	"github.com/forbearing/k8s/statefulset"
)

var (
	ctx, cancel = context.WithTimeout(context.Background(), time.Minute*10)
	namespace   = "test"
	kubeconfig  = filepath.Join(os.Getenv("HOME"), ".kube/config")
	filename    = "../../testdata/examples/statefulset.yaml"
	filename2   = "../../testdata/nginx/nginx-sts.yaml"
	name        = "mysts"
	name2       = "nginx-sts"
	label       = "type=statefulset"
)

func main() {
	//StatefulSet_Tools()
	StatefulSet_Scale()
}

func checkErr(name string, val interface{}, err error) {
	if err != nil {
		log.Printf("%s failed: %v\n", name, err)
	} else {
		log.Printf("%s success: %v\n", name, val)
	}
}
func cleanup(handler *statefulset.Handler) {
	handler.Delete(name)
	k8s.DeleteF(ctx, filename2, kubeconfig)
}
