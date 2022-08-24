package main

import (
	"log"
	"time"

	"github.com/forbearing/k8s/replicaset"
)

func ReplicaSet_Scale() {
	handler := replicaset.NewOrDie(ctx, "", namespace)
	defer cleanup(handler)

	if _, err := handler.Apply(filename); err != nil {
		log.Fatal(err)
	}
	log.Println("Wait Ready")
	handler.WaitReady(name)
	handler.Scale(name, 6)
	log.Println("Wait Ready Again")
	handler.WaitReady(name)
	time.Sleep(time.Second)
}
