package main

import (
	"log"
	"time"

	"github.com/forbearing/k8s/clusterrolebinding"
)

func RoleBinding_Tools() {
	handler, err := clusterrolebinding.New(ctx, kubeconfig)
	if err != nil {
		panic(err)
	}
	defer cleanup(handler)

	handler.Apply(filename)
	crb, err := handler.Get(name)
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

		role, err := handler.GetRole(crb)
		checkErr("GetRole", role, err)
		age, err := handler.GetAge(crb)
		checkErr("GetAge", age, err)
		subjects, err := handler.GetSubjects(crb)
		checkErr("GetSubjects", subjects, err)

		end := time.Now()
		log.Println("===== Get RoleBinding By Object Cost Time:", end.Sub(begin))
	}

	getByName()
	getByObj()

	// Output:

	//2022/07/11 12:42:26 ===== Get RoleBinding By Name
	//2022/07/11 12:42:26 GetRole success: &{ClusterRole myclusterrole}
	//2022/07/11 12:42:26 GetAge success: 1h13m56.149103s
	//2022/07/11 12:42:26 GetSubjects success: [{ServiceAccount mysa test}]
	//2022/07/11 12:42:26 ===== Get RoleBinding By Name Cost Time: 46.638167ms
	//2022/07/11 12:42:26
	//2022/07/11 12:42:26 ===== Get RoleBinding By Object
	//2022/07/11 12:42:26 GetRole success: &{ClusterRole myclusterrole}
	//2022/07/11 12:42:26 GetAge success: 1h13m56.154814s
	//2022/07/11 12:42:26 GetSubjects success: [{ServiceAccount mysa test}]
	//2022/07/11 12:42:26 ===== Get RoleBinding By Object Cost Time: 8.164µs
}
