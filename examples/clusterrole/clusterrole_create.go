package main

import (
	"io/ioutil"

	"github.com/forbearing/k8s/clusterrole"
	"k8s.io/apimachinery/pkg/runtime"
)

func ClusterRole_Create() {
	// New returns a handler used to multiples deployment.
	handler, err := clusterrole.New(ctx, kubeconfig)
	if err != nil {
		panic(err)
	}
	//defer cleanup(handler)

	// 1. create clusterrole from file.
	deploy, err := handler.Create(filename)
	checkErr("create clusterrole from file", "", err)
	handler.Delete(name)

	// 2. create clusterrole from bytes.
	var data []byte
	if data, err = ioutil.ReadFile(filename); err != nil {
		panic(err)
	}
	_, err = handler.Create(data)
	checkErr("create clusterrole from bytes", "", err)
	handler.Delete(name)

	// 3. create clusterrole from *rbacv1.ClusterRole.
	_, err = handler.Create(deploy)
	checkErr("create clusterrole from *rbacv1.ClusterRole", "", err)
	handler.Delete(name)

	// 4. create clusterrole from rbacv1.ClusterRole.
	_, err = handler.Create(*deploy)
	checkErr("create clusterrole from rbacv1.ClusterRole", "", err)
	handler.Delete(name)

	// 5. create clusterrole from object.Runtime.
	object := runtime.Object(deploy)
	_, err = handler.Create(object)
	checkErr("create clusterrole from runtime.Object", "", err)
	handler.Delete(name)
}
