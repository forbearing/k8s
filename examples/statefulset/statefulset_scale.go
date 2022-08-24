package main

import (
	"log"

	"github.com/forbearing/k8s/statefulset"
)

func StatefulSet_Scale() {
	handler := statefulset.NewOrDie(ctx, "", namespace)
	defer cleanup(handler)
	if _, err := handler.Apply(filename); err != nil {
		log.Fatal(err)
	}

	log.Println("Wait Ready")
	handler.WaitReady(name)
	sts, err := handler.Scale(name, 6)
	checkErr("Scale StatefulSet", sts.Name, err)
	log.Println("Wait Ready Again")
	handler.WaitReady(name)
}
