package main

import (
	"log"

	"github.com/forbearing/k8s/deployment"
	"github.com/forbearing/k8s/dynamic"
)

var (
	namespace = "test"
)

func main() {
	//Dynamic_Create()
	//Dynamic_Update()
	//Dynamic_Apply()
	//Dynamic_Delete()
	//Dynamic_Get()
	//Dynamic_List()
	//Dynamic_Watch_Single()
}

func cleanup(handler *dynamic.Handler) {
	// if call Delete() to delete a k8s resource and the passed parameter type is string,
	// you should always to explicitly specify the GroupVersionKind by WithGVK() method to delete it.
	handler.WithGVK(deployment.GVK()).Delete(deployUnstructName)
	handler.DeleteFromFile("../../testdata/examples/deployment.yaml")
	handler.DeleteFromFile("../../testdata/examples/deployment.json")
	handler.Delete(podUnstructData)
	handler.Delete(nsUnstructData)
	handler.Delete(pvUnstructData)
	handler.Delete(crUnstructData)
}

func checkErr(name string, val interface{}, err error) {
	if err != nil {
		log.Printf("%s failed: %v\n", name, err)
	} else {
		log.Printf("%s success: %v\n", name, val)
	}
}
