package main

import (
	"log"

	"github.com/forbearing/k8s"
	"github.com/forbearing/k8s/deployment"
	"github.com/forbearing/k8s/statefulset"
)

func K8S_Delete() {
	handler := k8s.NewOrDie(ctx, "", namespace)

	if _, err := handler.Apply(deployFile); err != nil {
		log.Fatal(err)
	}
	if _, err := handler.Apply(stsFile); err != nil {
		log.Fatal(err)
	}
	err := handler.DeleteFromFile(deployFile)
	checkErr("delete deployment", "", err)
	err = handler.DeleteFromFile(stsFile)
	checkErr("delete statefulset", "", err)

	if _, err := handler.Apply(deployFile); err != nil {
		log.Fatal(err)
	}
	if _, err := handler.Apply(stsFile); err != nil {
		log.Fatal(err)
	}
	err = handler.WithGVK(deployment.GVK).Delete(deployName)
	checkErr("delete deployment", "", err)
	err = handler.WithGVK(statefulset.GVK).Delete(stsName)
	checkErr("delete statefulset", "", err)

	// Output
	//2022/10/04 00:34:08 delete deployment success:
	//2022/10/04 00:34:08 delete statefulset success:
	//2022/10/04 00:34:08 delete deployment success:
	//2022/10/04 00:34:08 delete statefulset success:
}
