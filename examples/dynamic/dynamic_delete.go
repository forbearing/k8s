package main

import (
	"context"
	"log"

	"github.com/forbearing/k8s/deployment"
	"github.com/forbearing/k8s/dynamic"
)

func Dynamic_Delete() {
	handler := dynamic.NewOrDie(context.TODO(), "", namespace)
	defer cleanup(handler)
	if _, err := handler.Apply(deployUnstructData); err != nil {
		log.Fatal(err)
	}
	if _, err := handler.Apply(podUnstructData); err != nil {
		log.Fatal(err)
	}

	if _, err := handler.Apply(nsUnstructData); err != nil {
		log.Fatal(err)
	}

	if _, err := handler.Apply(pvUnstructData); err != nil {
		log.Fatal(err)
	}

	if _, err := handler.Apply(crUnstructData); err != nil {
		log.Fatal(err)
	}

	// if call Delete() to delete a k8s resource and the passed parameter type is string,
	// you should always to explicitly specify the GroupVersionKind by WithGVK() method to delete it.
	err := handler.WithGVK(deployment.GVK()).Delete(deployUnstructName)
	checkErr("delete deployment", "", err)
	handler.Delete(podUnstructData)
	checkErr("delete pod", "", err)
	handler.Delete(nsUnstructData)
	checkErr("delete namespace", "", err)
	handler.Delete(pvUnstructData)
	checkErr("delete persistentvolume", "", err)
	handler.Delete(crUnstructData)
	checkErr("delete clusterrole", "", err)
}
