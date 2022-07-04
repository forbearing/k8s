package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/forbearing/k8s/deployment"
)

var (
	ctx, cancel = context.WithTimeout(context.Background(), time.Second*10)
	namespace   = "test"
	kubeconfig  = filepath.Join(os.Getenv("HOME"), ".kube/config")
	filename    = "../../testData/examples/deployment.yaml"
	name        = "mydep"
	label       = "type=deployment"
)

func main() {
	Deployment_Create()
	Deployment_Update()
	Deployment_Apply()
	Deployment_Delete()
	Deployment_Get()
	Deployment_List()
	Deployment_Watch()
}

func myerr(name string, err error) {
	if err != nil {
		log.Printf("%s failed: %v\n", name, err)
	} else {
		log.Printf("%s success.\n", name)
	}
}

// create will delete or prune created deployments.
func create(handler *deployment.Handler) {
	handler.Delete(rawName)
	handler.DeleteFromFile(filename)
}
