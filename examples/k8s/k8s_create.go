package main

import (
	"github.com/forbearing/k8s"
)

func K8S_Create() {
	handler := k8s.NewOrDie(ctx, "", namespace)

	_, err := handler.Create(deployFile)
	checkErr("create deployment", deployName, err)

	_, err = handler.Create(stsFile)
	checkErr("create statefulset", stsName, err)

	_, err = handler.Create(podFile)
	checkErr("create pod", podName, err)

	handler.DeleteFromFile(deployFile)
	handler.DeleteFromFile(stsFile)
	handler.DeleteFromFile(podFile)
	//handler.WithGVK(deployment.GVK()).Delete(deployName)
	//handler.WithGVK(statefulset.GVK()).Delete(stsName)
	//handler.WithGVK(pod.GVK()).Delete(podName)

	// Output:
	//2022/09/04 16:10:21 create deployment success: mydep
	//2022/09/04 16:10:21 create statefulset success: mysts
	//2022/09/04 16:10:21 create pod success: mypod
}
