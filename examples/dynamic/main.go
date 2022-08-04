package main

import (
	"log"

	"github.com/forbearing/k8s/dynamic"
)

func main() {
	Dynamic_Create()
	//Dynamic_Update()
	//Dynamic_Apply()
	//Dynamic_Get()
}

func cleanup(handler *dynamic.Handler) {
	//logrus.Info("===== cleanup =====")
	//logrus.Println(handler.Namespace("test").Delete(deployUnstructName))
	//logrus.Println(handler.Namespace("test").Group("").Resource("pods").Delete(podUnstructData))
	//logrus.Println(handler.Group("").Resource("namespaces").Delete(nsUnstructData))
	//logrus.Println(handler.Group("").Resource("persistentvolumes").Delete(pvUnstructData))
	//logrus.Println(handler.Group("rbac.authorization.k8s.io").Resource("clusterroles").Delete(crUnstructName))
	handler.Namespace("test").Delete(deployUnstructName)
	handler.Namespace("test").Group("").Resource("pods").Delete(podUnstructData)
	handler.Group("").Resource("namespaces").Delete(nsUnstructData)
	handler.Group("").Resource("persistentvolumes").Delete(pvUnstructData)
	handler.Group("rbac.authorization.k8s.io").Resource("clusterroles").Delete(crUnstructName)
}

func checkErr(name string, val interface{}, err error) {
	if err != nil {
		log.Printf("%s failed: %v\n", name, err)
	} else {
		log.Printf("%s success: %v\n", name, val)
	}
}
