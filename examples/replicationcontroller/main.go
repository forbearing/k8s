package main

import (
	"context"
	"log"
	"time"

	"github.com/forbearing/k8s"
	"github.com/forbearing/k8s/replicationcontroller"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	ctx, cancel = context.WithTimeout(context.Background(), time.Minute*10)
	kubeconfig  = clientcmd.RecommendedHomeFile
	namespace   = "test"
	filename    = "../../testdata/examples/replicationcontroller.yaml"
	name        = "myrc"
	filename2   = "../../testdata/nginx/nginx-rc.yaml"
	name2       = "nginx-rc"
)

func main() {
	defer cancel()

	ReplicationController_Tools()
	//ReplicationController_Scale()
}

func checkErr(name string, val interface{}, err error) {
	if err != nil {
		log.Printf("%s failed: %v\n", name, err)
	} else {
		log.Printf("%s success: %v\n", name, val)
	}
}

// cleanup will delete or prune created deployments.
func cleanup(handler *replicationcontroller.Handler) {
	log.Println(handler.Delete(name))
	log.Println(k8s.DeleteF(ctx, kubeconfig, filename2))
}
