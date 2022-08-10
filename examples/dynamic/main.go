package main

import (
	"log"

	"github.com/forbearing/k8s/clusterrole"
	"github.com/forbearing/k8s/dynamic"
	"github.com/forbearing/k8s/namespace"
	"github.com/forbearing/k8s/persistentvolume"
	"github.com/forbearing/k8s/pod"
)

func main() {
	Dynamic_Create()
	//Dynamic_Update()
	//Dynamic_Apply()
	//Dynamic_Get()
}

func cleanup(handler *dynamic.Handler) {
	handler.WithNamespace("test").Delete(deployUnstructName)
	handler.WithNamespace("test").WithGVR(pod.GVR()).Delete(podUnstructData)
	handler.WithGVR(namespace.GVR()).Delete(nsUnstructData)
	handler.WithGVR(persistentvolume.GVR()).Delete(pvUnstructData)
	handler.WithGVR(clusterrole.GVR()).Delete(crUnstructName)
}

func checkErr(name string, val interface{}, err error) {
	if err != nil {
		log.Printf("%s failed: %v\n", name, err)
	} else {
		log.Printf("%s success: %v\n", name, val)
	}
}
