package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/forbearing/k8s/ingress"
)

var (
	ctx, cancel = context.WithTimeout(context.Background(), time.Minute*10)
	namespace   = "test"
	kubeconfig  = filepath.Join(os.Getenv("HOME"), ".kube/config")
	filename    = "../../testdata/examples/ingress.yaml"
	filename2   = "../../testdata/examples/ingress-multi.yaml"
	update1File = "../../testdata/examples/ingress-update1.yaml"
	update2File = "../../testdata/examples/ingress-update2.yaml"
	update3File = "../../testdata/examples/ingress-update3.yaml"
	name        = "mying"
	name2       = "mying-multi"
	label       = "type=ingress"
)

func main() {
	Ingress_Tools()
}

func checkErr(name string, err error) {
	if err != nil {
		log.Printf("%s failed: %v\n", name, err)
	} else {
		log.Printf("%s success.\n", name)
	}
}

// cleanup will delete or prune created deployments.
func cleanup(handler *ingress.Handler) {
	handler.Delete(name)
	handler.Delete(name2)
	handler.DeleteFromFile(update1File)
	handler.DeleteFromFile(update2File)
	handler.DeleteFromFile(update3File)
}
