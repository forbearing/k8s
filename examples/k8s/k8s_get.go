package main

import (
	"log"

	"github.com/forbearing/k8s"
	"github.com/forbearing/k8s/deployment"
	"github.com/forbearing/k8s/pod"
	"github.com/forbearing/k8s/statefulset"
)

func K8S_Get() {
	handler := k8s.NewOrDie(ctx, "", namespace)

	if _, err := handler.Apply(deployFile); err != nil {
		log.Fatal(err)
	}
	if _, err := handler.Apply(stsFile); err != nil {
		log.Fatal(err)
	}
	if _, err := handler.Apply(podFile); err != nil {
		log.Fatal(err)
	}

	deployObj, err := handler.WithGVK(deployment.GVK).Get(deployName)
	checkErr("get deployment", deployObj.GetName(), err)
	stsObj, err := handler.WithGVK(statefulset.GVK).Get(stsName)
	checkErr("get statefulset", stsObj.GetName(), err)
	podObj, err := handler.GetFromFile(podFile)
	checkErr("get pod", podObj.GetName(), err)

	handler.DeleteFromFile(deployFile)
	handler.DeleteFromFile(stsFile)
	handler.DeleteFromFile(podFile)
	handler.WithGVK(deployment.GVK).Delete(deployName)
	handler.WithGVK(statefulset.GVK).Delete(stsName)
	handler.WithGVK(pod.GVK).Delete(podName)

	// Output:
	//2022/10/04 00:34:45 get deployment success: mydep
	//2022/10/04 00:34:45 get statefulset success: mysts
	//2022/10/04 00:34:45 get pod success: mypod
}
