package main

import (
	"github.com/forbearing/k8s"
	"github.com/forbearing/k8s/deployment"
	"github.com/forbearing/k8s/pod"
	"github.com/forbearing/k8s/statefulset"
)

func K8S_Create() {
	handler := k8s.NewOrDie(ctx, "", namespace)

	deployObj, err := handler.Create(deployFile)
	checkErr("create deployment", deployObj.GetName(), err)

	stsObj, err := handler.Create(stsFile)
	checkErr("create statefulset", stsObj.GetName(), err)

	podObj, err := handler.Create(podFile)
	checkErr("create pod", podObj.GetName(), err)

	handler.DeleteFromFile(deployFile)
	handler.DeleteFromFile(stsFile)
	handler.DeleteFromFile(podFile)
	handler.WithGVK(deployment.GVK).Delete(deployName)
	handler.WithGVK(statefulset.GVK).Delete(stsName)
	handler.WithGVK(pod.GVK).Delete(podName)

	// Output:
	//2022/10/04 00:32:47 create deployment success: mydep
	//2022/10/04 00:32:47 create statefulset success: mysts
	//2022/10/04 00:32:47 create pod success: mypod
}
