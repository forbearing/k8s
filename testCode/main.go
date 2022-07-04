package main

import (
	"context"
	"os"
	"path/filepath"
	"time"
)

var (
	ctx, cancel = context.WithTimeout(context.Background(), time.Second*10)
	namespace   = "test"
	kubeconfig  = filepath.Join(os.Getenv("HOME"), ".kube/config")
	filename    = "../testData/examples/pod.yaml"
	name        = "mypod"
	label       = "type=pod"
)

func main() {
	defer cancel()

	Pod()
}
