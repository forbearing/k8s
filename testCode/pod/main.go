package main

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"github.com/forbearing/k8s/pod"
)

var (
	ctx, cancel = context.WithTimeout(context.Background(), time.Minute*10)
	namespace   = "test"
	kubeconfig  = filepath.Join(os.Getenv("HOME"), ".kube/config")
	filename    = "../../testData/examples/pod.yaml"
	name        = "mypod"
	label       = "type=pod"
)

func main() {
	defer cancel()

	Pod_Tools()
}

func cleanup(handler *pod.Handler) {
	handler.Delete(name)
}
