package main

import (
	"log"
	"time"

	"github.com/forbearing/k8s/secret"
)

func Secret_Tools() {
	handler, err := secret.New(ctx, namespace, kubeconfig)
	if err != nil {
		panic(err)
	}
	sec, err := handler.Get(name)
	if err != nil {
		panic(err)
	}

	getByName := func() {
		log.Println("===== Get Secret By Name")
		begin := time.Now()

		typ, err := handler.GetType(name)
		checkErr("GetType", typ, err)
		numData, err := handler.GetNumData(name)
		checkErr("GetNumData", numData, err)
		age, err := handler.GetAge(name)
		checkErr("GetAge", age, err)

		end := time.Now()
		log.Println("===== Get Secret By Name Cost Time:", end.Sub(begin))
		log.Println()
	}

	getByObj := func() {
		log.Println("===== Get Secret By Object")
		begin := time.Now()

		typ, err := handler.GetType(sec)
		checkErr("GetType", typ, err)
		numData, err := handler.GetNumData(sec)
		checkErr("GetNumData", numData, err)
		age, err := handler.GetAge(sec)
		checkErr("GetAge", age, err)

		end := time.Now()
		log.Println("===== Get Secret By Object Cost Time:", end.Sub(begin))
	}

	getByName()
	getByObj()

	// Output:

	//2022/07/10 22:17:12 ===== Get Secret By Name
	//2022/07/10 22:17:12 GetType success: Opaque
	//2022/07/10 22:17:12 GetNumData success: 2
	//2022/07/10 22:17:12 GetAge success: 1h3m22.90115s
	//2022/07/10 22:17:12 ===== Get Secret By Name Cost Time: 8.955858ms
	//2022/07/10 22:17:12
	//2022/07/10 22:17:12 ===== Get Secret By Object
	//2022/07/10 22:17:12 GetType success: Opaque
	//2022/07/10 22:17:12 GetNumData success: 2
	//2022/07/10 22:17:12 GetAge success: 1h3m22.901185s
	//2022/07/10 22:17:12 ===== Get Secret By Object Cost Time: 6.526µs
}
