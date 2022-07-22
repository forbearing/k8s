package main

import (
	"io/ioutil"

	"github.com/forbearing/k8s/deployment"
	"k8s.io/apimachinery/pkg/runtime"
)

func Deployment_Update() {
	// New returns a handler used to multiples deployment.
	handler, err := deployment.New(ctx, namespace, kubeconfig)
	if err != nil {
		panic(err)
	}
	defer cleanup(handler)

	deploy, err := handler.Apply(filename)
	if err != nil {
		panic(err)
	}
	if _, err := handler.Apply(unstructData); err != nil {
		panic(err)
	}

	// 1. update deployment from file.
	_, err = handler.Update(updateFile)
	checkErr("update deployment from file", "", err)

	// 2. update deployment from bytes.
	var data []byte
	if data, err = ioutil.ReadFile(updateFile); err != nil {
		panic(err)
	}
	_, err = handler.Update(data)
	checkErr("update deployment from bytes", "", err)

	// 3. update deployment from *appsv1.Deployment
	_, err = handler.Update(deploy)
	checkErr("update deployment from *appsv1.Deployment", "", err)

	// 4. update deployment from appsv1.Deployment
	_, err = handler.Update(*deploy)
	checkErr("update deployment from appsv1.Deployment", "", err)

	// 5. update deployment from runtime.Object.
	object := runtime.Object(deploy)
	_, err = handler.Update(object)
	checkErr("update deployment from runtime.Object", "", err)

	// 6. update deployment from unstructured data, aka map[string]interface{}.
	_, err = handler.Update(unstructData)
	checkErr("update deployment from unstructured data", "", err)

	// Output:

	//2022/07/22 22:02:25 apply deployment from file (deployment not exists) success:
	//2022/07/22 22:02:25 apply deployment from file (deployment exists) success:
	//2022/07/22 22:02:25 apply deployment from bytes (deployment not exists) success:
	//2022/07/22 22:02:25 apply deployment from bytes (deployment exists) success:
	//2022/07/22 22:02:25 apply deployment from *appsv1.Deployment (deployment not exists) success:
	//2022/07/22 22:02:25 apply deployment from *appsv1.Deployment (deployment exists) success:
	//2022/07/22 22:02:26 apply deployment from appsv1.Deployment (deployment not exists) success:
	//2022/07/22 22:02:26 apply deployment from appsv1.Deployment (deployment exists) success:
	//2022/07/22 22:02:26 apply deployment from runtime.Object (deployment not exists) success:
	//2022/07/22 22:02:27 apply deployment from runtime.Object (deployment exists) success:
	//2022/07/22 22:02:27 apply deployment from unstructured data (deployment not exists) success:
	//2022/07/22 22:02:28 apply deployment from unstructured data (deployment exists) success:
}
