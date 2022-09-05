package main

import (
	"io/ioutil"

	"github.com/forbearing/k8s/deployment"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

func Deployment_Update() {
	// New returns a handler used to multiples deployment.
	handler, err := deployment.New(ctx, kubeconfig, namespace)
	if err != nil {
		panic(err)
	}
	defer cleanup(handler)

	if _, err := handler.Apply(filename); err != nil {
		panic(err)
	}
	if _, err := handler.Apply(unstructData); err != nil {
		panic(err)
	}

	// 1. update deployment from file.
	deploy, err := handler.Update(updateFile)
	checkErr("update deployment from file", deploy.Name, err)

	// 2. update deployment from bytes.
	var data []byte
	if data, err = ioutil.ReadFile(updateFile); err != nil {
		panic(err)
	}
	deploy2, err := handler.Update(data)
	checkErr("update deployment from bytes", deploy2.Name, err)

	// 3. update deployment from *appsv1.Deployment
	deploy3, err := handler.Update(deploy2)
	checkErr("update deployment from *appsv1.Deployment", deploy3.Name, err)

	// 4. update deployment from appsv1.Deployment
	deploy4, err := handler.Update(*deploy3)
	checkErr("update deployment from appsv1.Deployment", deploy4.Name, err)

	// 5. update deployment from runtime.Object.
	deploy5, err := handler.Update(runtime.Object(deploy4))
	checkErr("update deployment from runtime.Object", deploy5.Name, err)

	// 6. update deployment from *unstructured.Unstructured
	deploy6, err := handler.Update(&unstructured.Unstructured{Object: unstructData})
	checkErr("update deployment from *unstructured.Unstructured", deploy6.Name, err)

	// 7. update deployment from unstructured.Unstructured
	deploy7, err := handler.Update(unstructured.Unstructured{Object: unstructData})
	checkErr("update deployment from unstructured.Unstructured", deploy7.Name, err)

	// 8. update deployment from map[string]interface{}.
	deploy8, err := handler.Update(unstructData)
	checkErr("update deployment from map[string]interface{}", deploy8.Name, err)

	// Output:

	//2022/09/05 16:52:13 update deployment from file success: mydep
	//2022/09/05 16:52:13 update deployment from bytes success: mydep
	//2022/09/05 16:52:13 update deployment from *appsv1.Deployment success: mydep
	//2022/09/05 16:52:13 update deployment from appsv1.Deployment success: mydep
	//2022/09/05 16:52:13 update deployment from runtime.Object success: mydep
	//2022/09/05 16:52:13 update deployment from *unstructured.Unstructured success: mydep-unstruct
	//2022/09/05 16:52:13 update deployment from unstructured.Unstructured success: mydep-unstruct
	//2022/09/05 16:52:13 update deployment from map[string]interface{} success: mydep-unstruct
}
