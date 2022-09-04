package main

import (
	"bytes"
	"context"
	"fmt"
	"log"

	"github.com/forbearing/k8s/pod"
	"github.com/forbearing/k8s/util/errors"
)

func PortForward_Pod() {

	handler := pod.NewOrDie(context.Background(), "", namespace)
	//defer cleanup(handler)

	if _, err := handler.Create(filename); errors.IgnoreAlreadyExists(err) != nil {
		log.Fatal(err)
	}
	fmt.Println(handler.IsReady(name)) // false
	handler.WaitReady(name)
	fmt.Println(handler.IsReady(name)) // true

	log.Println("PortForward")
	if err := handler.PortForward(name, 8080, 80); err != nil {
		log.Fatal(err)
	}
	if err := handler.PortForward(name, 8080, 80); err != nil {
		log.Fatal(err)
	}
	if err := handler.PortForward(name, 8080, 80); err != nil {
		log.Fatal(err)
	}
	if err := handler.PortForward(name, 8080, 80); err != nil {
		log.Fatal(err)
	}
	// Output
	//false
	//true
	//Forwarding from 127.0.0.1:8080 -> 80
	//Forwarding from [::1]:8080 -> 80
	//Handling connection for 8080
	//Handling connection for 8080
	//Handling connection for 8080
	//Handling connection for 8080
	//if err := handler.PortForward(name, 8080, 80); err != nil {
	//    log.Fatal(err)
	//}

	fmt.Println()
	log.Println("PortForwardWithStream")
	stdout, stderr := new(bytes.Buffer), new(bytes.Buffer)
	if err := handler.PortForwardWithStream(name, 8080, 80, stdout, stderr); err != nil {
		log.Fatal(err)
	}
	fmt.Println(stdout)
	fmt.Println(stderr)

	stdout.Reset()
	stderr.Reset()
	if err := handler.PortForwardWithStream(name, 8080, 80, stdout, stderr); err != nil {
		log.Fatal(err)
	}
	fmt.Println(stdout)
	fmt.Println(stderr)

	stdout.Reset()
	stderr.Reset()
	if err := handler.PortForwardWithStream(name, 8080, 80, stdout, stderr); err != nil {
		log.Fatal(err)
	}
	fmt.Println(stdout)
	fmt.Println(stderr)

	stdout.Reset()
	stderr.Reset()
	if err := handler.PortForwardWithStream(name, 8080, 80, stdout, stderr); err != nil {
		log.Fatal(err)
	}
	fmt.Println(stdout)
	fmt.Println(stderr)
	// Output:
	//true
	//true
	//^CForwarding from 127.0.0.1:8080 -> 80
	//Forwarding from [::1]:8080 -> 80
	//Handling connection for 8080
	//Handling connection for 8080
	//Handling connection for 8080
	//Handling connection for 8080
}
