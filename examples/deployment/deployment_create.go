package main

import (
	"io/ioutil"

	"github.com/forbearing/k8s/deployment"
	"k8s.io/apimachinery/pkg/runtime"
)

func Deployment_Create() {
	// New returns a handler used to multiples deployment.
	handler, err := deployment.New(ctx, kubeconfig, namespace)
	if err != nil {
		panic(err)
	}
	defer cleanup(handler)

	// 1. create deployment from file.
	deploy, err := handler.Create(filename)
	checkErr("create deployment from file", "", err)
	handler.Delete(name)

	// 2. create deployment from bytes.
	var data []byte
	if data, err = ioutil.ReadFile(filename); err != nil {
		panic(err)
	}
	_, err = handler.Create(data)
	checkErr("create deployment from bytes", "", err)
	handler.Delete(name)

	// 3. create deployment from *appsv1.Deployment.
	_, err = handler.Create(deploy)
	checkErr("create deployment from *appsv1.Deployment", "", err)
	handler.Delete(name)

	// 4. create deployment from appsv1.Deployment.
	_, err = handler.Create(*deploy)
	checkErr("create deployment from appsv1.Deployment", "", err)
	handler.Delete(name)

	// 5. create deployment from object.Runtime.
	object := runtime.Object(deploy)
	_, err = handler.Create(object)
	checkErr("create deployment from runtime.Object", "", err)
	handler.Delete(name)

	// 6. create deployment from unstructured data, aka map[string]interface{}.
	_, err = handler.Create(unstructData)
	checkErr("create deployment from unstructured data", "", err)
	handler.Delete(name)

	// Output:

	//2022/07/22 20:53:36 create deployment from file success:
	//2022/07/22 20:53:36 create deployment from bytes success:
	//2022/07/22 20:53:36 create deployment from *appsv1.Deployment success:
	//2022/07/22 20:53:36 create deployment from appsv1.Deployment success:
	//2022/07/22 20:53:36 create deployment from runtime.Object success:
	//2022/07/22 20:53:36 create deployment from unstructured data success:
}
