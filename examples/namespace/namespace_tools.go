package main

import (
	"log"
	"time"

	k8snamespace "github.com/forbearing/k8s/namespace"
)

func Namespace_Tools() {
	handler, err := k8snamespace.New(ctx, kubeconfig)
	if err != nil {
		panic(err)
	}
	//defer cleanup(handler)

	handler.Apply(filename)
	ns, err := handler.Get(name)
	if err != nil {
		panic(err)
	}

	getByName := func() {
		log.Println("===== Get From Namespace Name")
		begin := time.Now()

		age, err := handler.GetAge(name)
		checkErr("GetAge", err)
		log.Println(age)

		end := time.Now()
		log.Println("===== Get From Namespace Name Cost Time:", end.Sub(begin))
		log.Println()
	}

	getByObj := func() {
		log.Println("===== Get From Namespace Object")
		begin := time.Now()

		age, err := handler.GetAge(ns)
		checkErr("GetAge", err)
		log.Println(age)

		end := time.Now()
		log.Println("===== Get From Namespace Object Cost Time:", end.Sub(begin))
		log.Println()
	}

	getByName()
	getByObj()

	// Output:

	//2022/07/08 15:35:04 ===== Get From Namespace Name
	//2022/07/08 15:35:04 GetAge success.
	//2022/07/08 15:35:04 791.47ms
	//2022/07/08 15:35:04 ===== Get From Namespace Name Cost Time: 131.720538ms
	//2022/07/08 15:35:04
	//2022/07/08 15:35:04 ===== Get From Namespace Object
	//2022/07/08 15:35:04 GetAge success.
	//2022/07/08 15:35:04 791.51ms
	//2022/07/08 15:35:04 ===== Get From Namespace Object Cost Time: 6.277µs
	//2022/07/08 15:35:04
}
