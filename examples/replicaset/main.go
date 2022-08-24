package main

import (
	"context"
	"log"
	"time"

	"github.com/forbearing/k8s"
	"github.com/forbearing/k8s/replicaset"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	ctx, cancel = context.WithTimeout(context.Background(), time.Minute*10)
	kubeconfig  = clientcmd.RecommendedHomeFile
	namespace   = "test"
	filename    = "../../testdata/examples/replicaset.yaml"
	filename2   = "../../testdata/nginx/nginx-rs.yaml"
	name        = "myrs"
	name2       = "nginx-rs"
)

func main() {
	defer cancel()

	ReplicaSet_Tools()
	//ReplicaSet_Scale()

}

func checkErr(name string, val interface{}, err error) {
	if err != nil {
		log.Printf("%s failed: %v\n", name, err)
	} else {
		log.Printf("%s success: %v\n", name, val)
	}
}

// cleanup will delete or prune created replicaset.
func cleanup(handler *replicaset.Handler) {
	handler.Delete(name)
	k8s.DeleteF(ctx, "", filename2)
}
