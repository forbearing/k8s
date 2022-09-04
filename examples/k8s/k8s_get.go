package main

import (
	"log"

	"github.com/forbearing/k8s"
	"github.com/forbearing/k8s/deployment"
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

	deploy, err := handler.WithGVK(deployment.GVK()).Get(deployName)
	checkErr("get deployment", deploy.GetName(), err)
	sts, err := handler.WithGVK(statefulset.GVK()).Get(stsName)
	checkErr("get statefulset", sts.GetName(), err)
	po, err := handler.GetFromFile(podFile)
	checkErr("get pod", po.GetName(), err)

	handler.DeleteFromFile(deployFile)
	handler.DeleteFromFile(stsFile)
	handler.DeleteFromFile(podFile)
	//handler.WithGVK(deployment.GVK()).Delete(deployName)
	//handler.WithGVK(statefulset.GVK()).Delete(stsName)
	//handler.WithGVK(pod.GVK()).Delete(podName)

	// Output:
	//2022/09/04 16:25:00 get deployment success: mydep
	//2022/09/04 16:25:00 get statefulset success: mysts
	//2022/09/04 16:25:00 get pod success: mypod
}
