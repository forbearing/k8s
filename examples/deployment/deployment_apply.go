package main

import (
	"io/ioutil"

	"github.com/forbearing/k8s/deployment"
	"k8s.io/apimachinery/pkg/runtime"
)

// ApplyXXX method will creates deployment if specified deployment not exist and
// it will updates deployment if already exist.
func Deployment_Apply() {
	// New returns a handler used to multiples deployment.
	handler, err := deployment.New(ctx, namespace, kubeconfig)
	if err != nil {
		panic(err)
	}
	defer cleanup(handler)

	// 1. apply deployment from file, it will updates the deployment, if already exist, or creates it.
	handler.Delete(name)
	_, err = handler.Apply(filename)
	checkErr("apply deployment from file (deployment not exists)", "", err)
	_, err = handler.Apply(updateFile)
	checkErr("apply deployment from file (deployment exists)", "", err)

	// 2. apply deployment from bytes, it will updates the deployment, if already exist, or creates it.
	handler.Delete(name)
	var data, data2 []byte
	if data, err = ioutil.ReadFile(filename); err != nil {
		panic(err)
	}
	if data2, err = ioutil.ReadFile(updateFile); err != nil {
		panic(err)
	}
	deploy, err := handler.Apply(data)
	checkErr("apply deployment from bytes (deployment not exists)", "", err)
	_, err = handler.Apply(data2)
	checkErr("apply deployment from bytes (deployment exists)", "", err)

	// 3. apply deployment from *appsv1.Deployment, it will updates the deployment, if already exist, or creates it.
	replicas := *deploy.Spec.Replicas
	handler.Delete(name)
	_, err = handler.Apply(deploy)
	checkErr("apply deployment from *appsv1.Deployment (deployment not exists)", "", err)
	replicas += 1
	deploy.Spec.Replicas = &replicas
	_, err = handler.Apply(deploy)
	checkErr("apply deployment from *appsv1.Deployment (deployment exists)", "", err)

	// 4. apply deployment from appsv1.Deployment, it will updates the deployment, if already exist, or creates it.
	handler.Delete(name)
	_, err = handler.Apply(*deploy)
	checkErr("apply deployment from appsv1.Deployment (deployment not exists)", "", err)
	replicas += 1
	deploy.Spec.Replicas = &replicas
	_, err = handler.Apply(*deploy)
	checkErr("apply deployment from appsv1.Deployment (deployment exists)", "", err)

	// 5. apply deployment from runtime.Object, it will updates the deployment, if already exist, or creates it.
	handler.Delete(name)
	object := runtime.Object(deploy)
	replicas += 1
	deploy.Spec.Replicas = &replicas
	object2 := runtime.Object(deploy)
	_, err = handler.Apply(object)
	checkErr("apply deployment from runtime.Object (deployment not exists)", "", err)
	_, err = handler.Apply(object2)
	checkErr("apply deployment from runtime.Object (deployment exists)", "", err)

	// 6. apply deployment from unstructured data, aka map[string]interface{}, it will updates the deployment, if already exist, or creates it.
	handler.Delete(unstructName)
	_, err = handler.Apply(unstructData)
	checkErr("apply deployment from unstructured data (deployment not exists)", "", err)
	_, err = handler.Apply(unstructData)
	checkErr("apply deployment from unstructured data (deployment exists)", "", err)

	// Output:

	//2022/07/22 22:52:49 apply deployment from file (deployment not exists) success:
	//2022/07/22 22:52:49 apply deployment from file (deployment exists) success:
	//2022/07/22 22:52:49 apply deployment from bytes (deployment not exists) success:
	//2022/07/22 22:52:49 apply deployment from bytes (deployment exists) success:
	//2022/07/22 22:52:49 apply deployment from *appsv1.Deployment (deployment not exists) success:
	//2022/07/22 22:52:49 apply deployment from *appsv1.Deployment (deployment exists) success:
	//2022/07/22 22:52:49 apply deployment from appsv1.Deployment (deployment not exists) success:
	//2022/07/22 22:52:50 apply deployment from appsv1.Deployment (deployment exists) success:
	//2022/07/22 22:52:50 apply deployment from runtime.Object (deployment not exists) success:
	//2022/07/22 22:52:51 apply deployment from runtime.Object (deployment exists) success:
	//2022/07/22 22:52:51 apply deployment from unstructured data (deployment not exists) success:
	//2022/07/22 22:52:51 apply deployment from unstructured data (deployment exists) success:
}
