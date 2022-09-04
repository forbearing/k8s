package main

import (
	"log"

	"github.com/forbearing/k8s"
)

func K8S_Update() {
	handler := k8s.NewOrDie(ctx, "", namespace)

	if _, err := handler.Create(deployFile); err != nil {
		log.Fatal(err)
	}
	_, err := handler.Update(deployFile)
	checkErr("update deployment", "", err)

	if _, err := handler.Create(stsFile); err != nil {
		log.Fatal(err)
	}
	_, err = handler.Update(stsFile)
	checkErr("update statefulset", "", err)

	handler.DeleteFromFile(deployFile)
	handler.DeleteFromFile(stsFile)
	//handler.WithGVK(deployment.GVK()).Delete(deployName)
	//handler.WithGVK(statefulset.GVK()).Delete(stsName)

	// Output:
	//2022/09/04 16:09:07 update deployment success:
	//2022/09/04 16:09:07 update statefulset success:
}
