package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"time"
)

var (
	ctx, cancel = context.WithTimeout(context.Background(), time.Minute*10)
	namespace   = "test"
	kubeconfig  = filepath.Join(os.Getenv("HOME"), ".kube/config")
	masterName  = "d11-k8s-master1"
	workerName  = "d11-k8s-worker1"
)

func main() {
	Node_Tools()
}
func checkErr(name string, val interface{}, err error) {
	if err != nil {
		log.Printf("%s failed: %v\n", name, err)
	} else {
		log.Printf("%s success: %v\n", name, val)
	}
}
