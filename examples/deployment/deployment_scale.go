package main

import (
	"log"

	"github.com/forbearing/k8s/deployment"
)

func Deployment_Scale() {
	handler := deployment.NewOrDie(ctx, "", namespace)
	defer cleanup(handler)

	if _, err := handler.Apply(filename); err != nil {
		log.Fatal(err)
	}
	log.Println("Wait Ready")
	handler.WaitReady(name)
	deploy, err := handler.Scale(name, 10)
	checkErr("Scale Deployment", deploy.Name, err)
	log.Println("Wait Ready Again.")
	handler.WaitReady(name)
}
