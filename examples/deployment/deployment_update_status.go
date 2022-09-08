package main

import (
	"log"

	"github.com/forbearing/k8s/deployment"
)

func Deployment_Update_Status() {
	handler := deployment.NewOrDie(ctx, "", namespace)
	//defer cleanup(handler)

	deploy, err := handler.Apply(filename)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Wait Ready")
	handler.WaitReady(name)

	deploy.Status.AvailableReplicas = 1
	deploy.Status.ObservedGeneration = 100
	deploy2, err := handler.UpdateStatus(deploy)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := handler.Update(deploy2); err != nil {
		log.Fatal(err)
	}
}
