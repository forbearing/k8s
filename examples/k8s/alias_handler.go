package main

import (
	"time"

	"github.com/forbearing/k8s"
)

func Alias() {
	deployHandler, err := k8s.NewDeployment(ctx, kubeconfig, namespace)
	if err != nil {
		panic(err)
	}
	podHandler, err := k8s.NewPod(ctx, kubeconfig, namespace)
	if err != nil {
		panic(err)
	}

	deploy, err := deployHandler.Apply(deployFile)
	checkErr("create deployment from file", deploy.Name, err)
	pod, err := podHandler.Apply(podFile)
	checkErr("create pod from file", pod.Name, err)

	time.Sleep(time.Second * 3)
	deployHandler.Delete(deployName)
	podHandler.Delete(podName)
}
