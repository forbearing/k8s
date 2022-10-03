package main

import (
	"log"

	"github.com/forbearing/k8s"
	"github.com/forbearing/k8s/deployment"
	"github.com/forbearing/k8s/statefulset"
)

func K8S_Update() {
	handler := k8s.NewOrDie(ctx, "", namespace)

	if _, err := handler.Create(deployFile); err != nil {
		log.Fatal(err)
	}
	deployObj, err := handler.Update(deployFile)
	checkErr("update deployment", deployObj.GetName(), err)

	if _, err := handler.Create(stsFile); err != nil {
		log.Fatal(err)
	}
	stsObj, err := handler.Update(stsFile)
	checkErr("update statefulset", stsObj.GetName(), err)

	handler.DeleteFromFile(deployFile)
	handler.DeleteFromFile(stsFile)
	handler.WithGVK(deployment.GVK).Delete(deployName)
	handler.WithGVK(statefulset.GVK).Delete(stsName)

	// Output:
	//2022/10/04 00:33:21 update deployment success: mydep
	//2022/10/04 00:33:21 update statefulset success: mysts
}
