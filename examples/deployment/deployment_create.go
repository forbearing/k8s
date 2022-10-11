package main

import (
	"io/ioutil"

	"github.com/forbearing/k8s/deployment"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
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
	checkErr("create deployment from file", deploy.Name, err)
	handler.Delete(name)

	// 2. create deployment from bytes.
	var data []byte
	if data, err = ioutil.ReadFile(filename); err != nil {
		panic(err)
	}
	deploy2, err := handler.Create(data)
	checkErr("create deployment from bytes", deploy2.Name, err)
	handler.Delete(name)

	// 3. create deployment from *appsv1.Deployment.
	deploy3, err := handler.Create(deploy2)
	checkErr("create deployment from *appsv1.Deployment", deploy3.Name, err)
	handler.Delete(name)

	// 4. create deployment from appsv1.Deployment.
	deploy4, err := handler.Create(*deploy3)
	checkErr("create deployment from appsv1.Deployment", deploy4.Name, err)
	handler.Delete(name)

	// 5. create deployment from runtime.Object.
	deploy5, err := handler.Create(runtime.Object(deploy4))
	checkErr("create deployment from runtime.Object", deploy5.Name, err)
	handler.Delete(name)

	// 6. create deployment from metav1.Object.
	deploy6, err := handler.Create(metav1.Object(deploy5))
	checkErr("create deployment from metav1.Object", deploy6.Name, err)
	handler.Delete(name)

	// 7. create deployment from *unstructured.Unstructured
	deploy7, err := handler.Create(&unstructured.Unstructured{Object: unstructData})
	checkErr("create deployment from *unstructured.Unstructured", deploy7.Name, err)
	handler.Delete(unstructName)

	// 8. create deployment from unstructured.Unstructured
	deploy8, err := handler.Create(unstructured.Unstructured{Object: unstructData})
	checkErr("create deployment from *unstructured.Unstructured", deploy8.Name, err)
	handler.Delete(unstructName)

	// 9. create deployment from map[string]interface{}.
	deploy9, err := handler.Create(unstructData)
	checkErr("create deployment from map[string]interface{}", deploy9.Name, err)
	handler.Delete(unstructData)

	// Output:

	//2022/10/08 22:20:24 create deployment from file success: mydep
	//2022/10/08 22:20:24 create deployment from bytes success: mydep
	//2022/10/08 22:20:24 create deployment from *appsv1.Deployment success: mydep
	//2022/10/08 22:20:24 create deployment from appsv1.Deployment success: mydep
	//2022/10/08 22:20:24 create deployment from runtime.Object success: mydep
	//2022/10/08 22:20:24 create deployment from metav1.Object success: mydep
	//2022/10/08 22:20:25 create deployment from *unstructured.Unstructured success: mydep-unstruct
	//2022/10/08 22:20:25 create deployment from *unstructured.Unstructured success: mydep-unstruct
	//2022/10/08 22:20:26 create deployment from map[string]interface{} success: mydep-unstruct
}
