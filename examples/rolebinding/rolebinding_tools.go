package main

import (
	"log"
	"time"

	"github.com/forbearing/k8s/rolebinding"
)

func RoleBinding_Tools() {
	handler, err := rolebinding.New(ctx, kubeconfig, namespace)
	if err != nil {
		panic(err)
	}
	defer cleanup(handler)

	handler.Apply(filename)
	rb, err := handler.Get(name)
	if err != nil {
		panic(err)
	}

	getByName := func() {
		log.Println("===== Get RoleBinding By Name")
		begin := time.Now()

		role, err := handler.GetRole(name)
		checkErr("GetRole", role, err)
		age, err := handler.GetAge(name)
		checkErr("GetAge", age, err)
		subjects, err := handler.GetSubjects(name)
		checkErr("GetSubjects", subjects, err)

		end := time.Now()
		log.Println("===== Get RoleBinding By Name Cost Time:", end.Sub(begin))
		log.Println()
	}

	getByObj := func() {
		log.Println("===== Get RoleBinding By Object")
		begin := time.Now()

		role, err := handler.GetRole(rb)
		checkErr("GetRole", role, err)
		age, err := handler.GetAge(rb)
		checkErr("GetAge", age, err)
		subjects, err := handler.GetSubjects(rb)
		checkErr("GetSubjects", subjects, err)

		end := time.Now()
		log.Println("===== Get RoleBinding By Object Cost Time:", end.Sub(begin))
	}

	getByName()
	getByObj()

	// Output:

	//2022/07/11 12:39:08 ===== Get RoleBinding By Name
	//2022/07/11 12:39:08 GetRole success: &{Role myrole}
	//2022/07/11 12:39:08 GetAge success: 243h37m22.31434s
	//2022/07/11 12:39:08 GetSubjects success: [{ServiceAccount mysa test}]
	//2022/07/11 12:39:08 ===== Get RoleBinding By Name Cost Time: 22.249564ms
	//2022/07/11 12:39:08
	//2022/07/11 12:39:08 ===== Get RoleBinding By Object
	//2022/07/11 12:39:08 GetRole success: &{Role myrole}
	//2022/07/11 12:39:08 GetAge success: 243h37m22.322935s
	//2022/07/11 12:39:08 GetSubjects success: [{ServiceAccount mysa test}]
	//2022/07/11 12:39:08 ===== Get RoleBinding By Object Cost Time: 12.289µs
}
