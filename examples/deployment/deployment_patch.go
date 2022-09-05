package main

import (
	"log"
	"time"

	"github.com/forbearing/k8s/deployment"
)

func Deployment_Patch() {
	handler := deployment.NewOrDie(ctx, "", namespace)
	defer cleanup(handler)

	deploy, err := handler.Apply(unstructData)
	if err != nil {
		log.Fatal(err)
	}
	handler.WaitReady(unstructName)

	modifiedDeploy := deploy.DeepCopy()
	replicas := *deploy.Spec.Replicas

	replicas += 1
	modifiedDeploy.Spec.Replicas = &replicas
	deploy2, err := handler.StrategicMergePatch(deploy, modifiedDeploy)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(*deploy.Spec.Replicas, *deploy2.Spec.Replicas) //1  2
	handler.WaitReady(unstructName)

	replicas += 1
	modifiedDeploy.Spec.Replicas = &replicas
	deploy3, err := handler.MergePatch(deploy, modifiedDeploy)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(*deploy.Spec.Replicas, *deploy3.Spec.Replicas) //1  3
	handler.WaitReady(unstructName)
	time.Sleep(time.Second * 3)

	// Output:

	//2022/09/05 16:54:11 1 2
	//2022/09/05 16:54:16 1 3
}
