package main

import (
	"context"
	"fmt"
	"log"

	"github.com/forbearing/k8s"
)

func Handler() {
	filename := "../../testdata/examples/deployment.yaml"

	handler := k8s.NewOrDie(context.TODO(), "", "test")
	deploy, err := handler.Apply(filename)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(deploy.GetName())
	if err := handler.DeleteFromFile(filename); err != nil {
		log.Fatal(err)
	}
}
