package main

import (
	"github.com/forbearing/k8s"
	"github.com/forbearing/k8s/deployment"
	"github.com/forbearing/k8s/pod"
	"github.com/forbearing/k8s/statefulset"
)

func K8S_Apply() {
	handler := k8s.NewOrDie(ctx, "", namespace)

	deployObj, err := handler.Apply(deployFile)
	checkErr("create deployment", deployObj.GetName(), err)

	stsObj, err := handler.Apply(stsFile)
	checkErr("create statefulset", stsObj.GetName(), err)

	podObj, err := handler.Apply(podFile)
	checkErr("create pod", podObj.GetName(), err)

	handler.DeleteFromFile(deployFile)
	handler.DeleteFromFile(stsFile)
	handler.DeleteFromFile(podFile)
	handler.WithGVK(deployment.GVK).Delete(deployName)
	handler.WithGVK(statefulset.GVK).Delete(stsName)
	handler.WithGVK(pod.GVK).Delete(podName)

	// Output:
	//2022/10/04 00:33:40 create deployment success: mydep
	//2022/10/04 00:33:40 create statefulset success: mysts
	//2022/10/04 00:33:40 create pod success: mypod
}
