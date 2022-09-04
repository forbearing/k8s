package main

import "github.com/forbearing/k8s"

func K8S_Apply() {
	handler := k8s.NewOrDie(ctx, "", namespace)

	_, err := handler.Apply(deployFile)
	checkErr("create deployment", deployName, err)

	_, err = handler.Apply(stsFile)
	checkErr("create statefulset", stsName, err)

	_, err = handler.Apply(podFile)
	checkErr("create pod", podName, err)

	handler.DeleteFromFile(deployFile)
	handler.DeleteFromFile(stsFile)
	handler.DeleteFromFile(podFile)
	//handler.WithGVK(deployment.GVK()).Delete(deployName)
	//handler.WithGVK(statefulset.GVK()).Delete(stsName)
	//handler.WithGVK(pod.GVK()).Delete(podName)

	// Output:
	//2022/09/04 16:08:40 create deployment success: mydep
	//2022/09/04 16:08:40 create statefulset success: mysts
	//2022/09/04 16:08:40 create pod success: mypod
}
