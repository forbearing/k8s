package main

import (
	"log"

	"github.com/forbearing/k8s/deployment"
)

func Deployment_Update_Status() {
	handler := deployment.NewOrDie(ctx, "", namespace)
	defer cleanup(handler)

	deploy, err := handler.Apply(filename)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Wait Ready")
	handler.WaitReady(name)

	copiedDeploy := deploy.DeepCopy()
	copiedDeploy.Status.Replicas = 10
	deploy2, err := handler.UpdateStatus(copiedDeploy)
	//deploy.Status.Replicas = 10
	//deploy2, err := handler.UpdateStatus(deploy)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(*deploy2.Spec.Replicas)
	log.Println(deploy2.Status.Replicas)
}
