package main

import (
	"log"
	"time"

	"github.com/forbearing/k8s/configmap"
)

func ConfigMap_Tools() {
	handler, err := configmap.New(ctx, namespace, kubeconfig)
	if err != nil {
		panic(err)
	}
	defer cleanup(handler)

	handler.Apply(filename)
	cm, err := handler.Get(name)
	if err != nil {
		panic(err)
	}

	getByName := func() {
		log.Println("===== Get ConfigMap By Name")
		begin := time.Now()

		data, err := handler.GetData(name)
		checkErr("GetData", data, err)
		numData, err := handler.NumData(name)
		checkErr("NumData", numData, err)
		age, err := handler.GetAge(name)
		checkErr("GetAge", age, err)

		end := time.Now()
		log.Println("===== Get ConfigMap By Name Cost Time:", end.Sub(begin))
		log.Println()
	}

	getByObj := func() {
		log.Println("===== Get ConfigMap By Object")
		begin := time.Now()

		data, err := handler.GetData(cm)
		checkErr("GetData", data, err)
		numData, err := handler.NumData(cm)
		checkErr("NumData", numData, err)
		age, err := handler.GetAge(cm)
		checkErr("GetAge", age, err)

		end := time.Now()
		log.Println("===== Get ConfigMap By Object Cost Time:", end.Sub(begin))
	}

	getByName()
	getByObj()

	// Output:

	//2022/07/11 14:03:49 ===== Get ConfigMap By Name
	//2022/07/11 14:03:49 GetData success: map[name:kubernetes version:v1.24 welcome:hello]
	//2022/07/11 14:03:49 NumData success: 3
	//2022/07/11 14:03:49 GetAge success: 662.83ms
	//2022/07/11 14:03:49 ===== Get ConfigMap By Name Cost Time: 35.621035ms
	//2022/07/11 14:03:49
	//2022/07/11 14:03:49 ===== Get ConfigMap By Object
	//2022/07/11 14:03:49 GetData success: map[name:kubernetes version:v1.24 welcome:hello]
	//2022/07/11 14:03:49 NumData success: 3
	//2022/07/11 14:03:49 GetAge success: 662.863ms
}
