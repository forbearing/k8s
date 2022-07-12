package main

import (
	"log"
	"time"

	"github.com/forbearing/k8s/ingress"
)

func Ingress_Tools() {
	handler, err := ingress.New(ctx, namespace, kubeconfig)
	if err != nil {
		panic(err)
	}
	defer cleanup(handler)

	handler.Apply(filename)
	handler.Apply(filename2)

	ing1, err := handler.Get(name)
	if err != nil {
		panic(err)
	}
	ing2, err := handler.Get(name2)
	if err != nil {
		panic(err)
	}

	getByName := func() {
		begin := time.Now()

		class1, err := handler.GetClass(name)
		checkErr("GetClass by ingress name", err)
		class2, err := handler.GetClass(name2)
		checkErr("GetClass by ingress name", err)
		log.Println(class1)
		log.Println(class2)

		host1, err := handler.GetHosts(name)
		checkErr("GetHosts by ingress name", err)
		host2, err := handler.GetHosts(name2)
		checkErr("GetHosts by ingress name", err)
		log.Println(host1)
		log.Println(host2)

		addr1, err := handler.GetAddress(name)
		checkErr("GetAddress by ingress name", err)
		addr2, err := handler.GetAddress(name2)
		checkErr("GetAddress by ingress name", err)
		log.Println(addr1)
		log.Println(addr2)

		age1, err := handler.GetAge(name)
		checkErr("GetAge by ingress name", err)
		age2, err := handler.GetAge(name2)
		checkErr("GetAge by ingress name", err)
		log.Println(age1)
		log.Println(age2)

		end := time.Now()
		log.Println("===== Get From Ingress Name Cost Time:", end.Sub(begin))
		log.Println()
	}

	getByObj := func() {
		begin := time.Now()

		class1, err := handler.GetClass(ing1)
		checkErr("GetClass by ingress object", err)
		class2, err := handler.GetClass(ing2)
		checkErr("GetClass by ingress object", err)
		log.Println(class1)
		log.Println(class2)

		host1, err := handler.GetHosts(ing1)
		checkErr("GetHosts by ingress object", err)
		host2, err := handler.GetHosts(ing2)
		checkErr("GetHosts by ingress object", err)
		log.Println(host1)
		log.Println(host2)
		addr1, err := handler.GetAddress(ing1)
		checkErr("GetAddress by ingress object", err)
		addr2, err := handler.GetAddress(ing2)
		checkErr("GetAddress by ingress object", err)
		log.Println(addr1)
		log.Println(addr2)

		age1, err := handler.GetAge(ing1)
		checkErr("GetAge by ingress object", err)
		age2, err := handler.GetAge(ing2)
		checkErr("GetAge by ingress object", err)
		log.Println(age1)
		log.Println(age2)

		end := time.Now()
		log.Println("===== Get From Ingress Object Cost Time:", end.Sub(begin))
		log.Println()
	}
	getByName()
	getByObj()
}
