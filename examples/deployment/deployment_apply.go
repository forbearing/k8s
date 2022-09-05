package main

import (
	"io/ioutil"

	"github.com/forbearing/k8s/deployment"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

// ApplyXXX method will creates deployment if specified deployment not exist and
// it will updates deployment if already exist.
func Deployment_Apply() {
	// New returns a handler used to multiples deployment.
	handler, err := deployment.New(ctx, kubeconfig, namespace)
	if err != nil {
		panic(err)
	}
	defer cleanup(handler)

	// 1. apply deployment from file, it will updates the deployment, if already exist, or creates it.
	handler.Delete(name)
	deploy, err := handler.Apply(filename)
	checkErr("apply deployment from file (deployment not exists)", deploy.Name, err)
	deploy, err = handler.Apply(updateFile)
	checkErr("apply deployment from file (deployment exists)", deploy.Name, err)

	// 2. apply deployment from bytes, it will updates the deployment, if already exist, or creates it.
	handler.Delete(name)
	var data, data2 []byte
	if data, err = ioutil.ReadFile(filename); err != nil {
		panic(err)
	}
	if data2, err = ioutil.ReadFile(updateFile); err != nil {
		panic(err)
	}
	deploy2, err := handler.Apply(data)
	checkErr("apply deployment from bytes (deployment not exists)", deploy2.Name, err)
	deploy2, err = handler.Apply(data2)
	checkErr("apply deployment from bytes (deployment exists)", deploy2.Name, err)

	// 3. apply deployment from *appsv1.Deployment, it will updates the deployment, if already exist, or creates it.
	replicas := *deploy2.Spec.Replicas
	handler.Delete(name)
	deploy3, err := handler.Apply(deploy2)
	checkErr("apply deployment from *appsv1.Deployment (deployment not exists)", deploy3.Name, err)
	replicas += 1
	deploy.Spec.Replicas = &replicas
	deploy3, err = handler.Apply(deploy)
	checkErr("apply deployment from *appsv1.Deployment (deployment exists)", deploy3.Name, err)

	// 4. apply deployment from appsv1.Deployment, it will updates the deployment, if already exist, or creates it.
	handler.Delete(name)
	deploy4, err := handler.Apply(*deploy)
	checkErr("apply deployment from appsv1.Deployment (deployment not exists)", deploy4.Name, err)
	replicas += 1
	deploy.Spec.Replicas = &replicas
	deploy4, err = handler.Apply(*deploy)
	checkErr("apply deployment from appsv1.Deployment (deployment exists)", deploy4.Name, err)

	// 5. apply deployment from runtime.Object, it will updates the deployment, if already exist, or creates it.
	handler.Delete(name)
	deploy5, err := handler.Apply(runtime.Object(deploy))
	checkErr("apply deployment from runtime.Object (deployment not exists)", deploy5.Name, err)
	replicas -= 1
	deploy.Spec.Replicas = &replicas
	deploy5, err = handler.Apply(runtime.Object(deploy))
	checkErr("apply deployment from runtime.Object (deployment exists)", deploy5.Name, err)

	// 6. apply deployment from *unstructured.Unstructured, it will updates the deployment, if already exist, or creates it.
	handler.Delete(unstructName)
	deploy6, err := handler.Apply(&unstructured.Unstructured{Object: unstructData})
	checkErr("apply deployment from *unstructured.Unstructured (deployment not exists)", deploy6.Name, err)
	deploy6, err = handler.Apply(&unstructured.Unstructured{Object: unstructData})
	checkErr("apply deployment from *unstructured.Unstructured (deployment exists)", deploy6.Name, err)

	// 7. apply deployment from unstructured.Unstructured, it will updates the deployment, if already exist, or creates it.
	handler.Delete(unstructName)
	deploy7, err := handler.Apply(unstructured.Unstructured{Object: unstructData})
	checkErr("apply deployment from unstructured.Unstructured (deployment not exists)", deploy7.Name, err)
	deploy7, err = handler.Apply(unstructured.Unstructured{Object: unstructData})
	checkErr("apply deployment from unstructured.Unstructured (deployment exists)", deploy7.Name, err)

	// 8. apply deployment map[string]interface{}, it will updates the deployment, if already exist, or creates it.
	handler.Delete(unstructName)
	deploy8, err := handler.Apply(unstructData)
	checkErr("apply deployment from map[string]interface{} (deployment not exists)", deploy8.Name, err)
	deploy8, err = handler.Apply(unstructData)
	checkErr("apply deployment from map[string]interface{} (deployment exists)", deploy8.Name, err)

	// Output:

	//2022/09/05 16:53:25 apply deployment from file (deployment not exists) success: mydep
	//2022/09/05 16:53:25 apply deployment from file (deployment exists) success: mydep
	//2022/09/05 16:53:25 apply deployment from bytes (deployment not exists) success: mydep
	//2022/09/05 16:53:25 apply deployment from bytes (deployment exists) success: mydep
	//2022/09/05 16:53:25 apply deployment from *appsv1.Deployment (deployment not exists) success: mydep
	//2022/09/05 16:53:25 apply deployment from *appsv1.Deployment (deployment exists) success: mydep
	//2022/09/05 16:53:25 apply deployment from appsv1.Deployment (deployment not exists) success: mydep
	//2022/09/05 16:53:26 apply deployment from appsv1.Deployment (deployment exists) success: mydep
	//2022/09/05 16:53:26 apply deployment from runtime.Object (deployment not exists) success: mydep
	//2022/09/05 16:53:27 apply deployment from runtime.Object (deployment exists) success: mydep
	//2022/09/05 16:53:27 apply deployment from *unstructured.Unstructured (deployment not exists) success: mydep-unstruct
	//2022/09/05 16:53:27 apply deployment from *unstructured.Unstructured (deployment exists) success: mydep-unstruct
	//2022/09/05 16:53:28 apply deployment from unstructured.Unstructured (deployment not exists) success: mydep-unstruct
	//2022/09/05 16:53:28 apply deployment from unstructured.Unstructured (deployment exists) success: mydep-unstruct
	//2022/09/05 16:53:29 apply deployment from map[string]interface{} (deployment not exists) success: mydep-unstruct
	//2022/09/05 16:53:29 apply deployment from map[string]interface{} (deployment exists) success: mydep-unstruct
}
