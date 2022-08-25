package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/forbearing/k8s/clusterrole"
)

var (
	ctx, cancel = context.WithTimeout(context.Background(), time.Minute*10)
	namespace   = "test"
	kubeconfig  = filepath.Join(os.Getenv("HOME"), ".kube/config")
	filename    = "../../testdata/examples/clusterrole.yaml"
	name        = "myclusterrole"
	label       = "type=clusterrole"
)

func main() {
	ClusterRole_Create()
	//ClusterRole_Update()
	//ClusterRole_Apply()
	//ClusterRole_Delete()
	//ClusterRole_Get()
	//ClusterRole_List()
	//ClusterRole_Watch()
	//ClusterRole_Informer()
	//ClusterRole_Tools()

}

func checkErr(name string, val interface{}, err error) {
	if err != nil {
		log.Printf("%s failed: %v\n", name, err)
	} else {
		log.Printf("%s success: %v\n", name, val)
	}
}

// cleanup will delete or prune created deployments.
func cleanup(handler *clusterrole.Handler) {
	handler.Delete(name)
}
