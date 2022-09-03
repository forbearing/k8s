package main

import (
	"log"

	"github.com/forbearing/k8s/clusterrole"
	"github.com/forbearing/k8s/deployment"
	"github.com/forbearing/k8s/dynamic"
	"github.com/forbearing/k8s/namespace"
	"github.com/forbearing/k8s/persistentvolume"
	"github.com/forbearing/k8s/pod"
)

func main() {
	//Dynamic_Create()
	//Dynamic_Update()
	Dynamic_Apply()
	//Dynamic_Get()
}

func cleanup(handler *dynamic.Handler) {
	handler.WithGVK(deployment.GVK()).Delete(deployUnstructName)
	handler.DeleteFromFile("../../testdata/examples/deployment.yaml")
	handler.DeleteFromFile("../../testdata/examples/deployment.json")
	handler.WithGVK(pod.GVK()).Delete(podUnstructData)
	handler.WithGVK(namespace.GVK()).Delete(nsUnstructData)
	handler.WithGVK(persistentvolume.GVK()).Delete(pvUnstructData)
	handler.WithGVK(clusterrole.GVK()).Delete(crUnstructName)
}

func checkErr(name string, val interface{}, err error) {
	if err != nil {
		log.Printf("%s failed: %v\n", name, err)
	} else {
		log.Printf("%s success: %v\n", name, val)
	}
}
