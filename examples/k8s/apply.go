package main

import (
	"context"
	"fmt"
	"log"

	"github.com/forbearing/k8s"
)

func Apply() {
	filename := "../../testdata/examples/all.yaml"

	//if err := k8s.ApplyF(context.TODO(), kubeconfig, filename, namespace); err != nil {
	//    log.Fatal(err)
	//}
	////time.Sleep(time.Second * 120)
	//if err := k8s.DeleteF(context.TODO(), kubeconfig, filename, namespace); err != nil {
	//    log.Fatal(err)
	//}

	fmt.Println()
	log.Println("first apply")
	if err := k8s.ApplyF(context.TODO(), kubeconfig, filename, namespace); err != nil {
		log.Fatal(err)
	}

	fmt.Println()
	log.Println(`second apply, will output "invalid" error message`)
	// update the pod will return "Invalid" error
	if err := k8s.ApplyF(context.TODO(), kubeconfig, filename, namespace); err != nil {
		log.Fatal(err)
	}
	fmt.Println()
	log.Println("third apply, no error")
	// ignore the "Invalid" error
	if err := k8s.ApplyF(context.TODO(), kubeconfig, filename, namespace, k8s.IgnoreInvalid); err != nil {
		log.Fatal(err)
	}

	//
	//
	fmt.Println()
	log.Println("first delete")
	if err := k8s.DeleteF(context.TODO(), kubeconfig, filename, namespace); err != nil {
		log.Fatal(err)
	}
	fmt.Println()
	log.Println(`second delete, will output "not found" error message`)
	// output the "not found" error message
	if err := k8s.DeleteF(context.TODO(), kubeconfig, filename, namespace); err != nil {
		log.Fatal(err)
	}
	fmt.Println()
	log.Println("third delete, no error")
	// no error
	if err := k8s.DeleteF(context.TODO(), kubeconfig, filename, namespace, k8s.IgnoreNotFound); err != nil {
		log.Fatal(err)
	}

}
