package main

import (
	"io/ioutil"

	"github.com/forbearing/k8s/deployment"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

func Deployment_Get() {
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

	// 1. get deployment by name.
	deploy1, err := handler.Get(name)
	checkErr("get deployment by name", deploy1.Name, err)

	// 2. get deployment from file.
	// You should always explicitly call GetFromFile to get a deployment from file.
	// if the parameter type passed to the `Get` method is string, the `Get`
	// method will call GetByName not GetFromFile.
	deploy2, err := handler.GetFromFile(filename)
	checkErr("get deployment from file", deploy2.Name, err)

	// 3. get deployment from bytes.
	var data []byte
	if data, err = ioutil.ReadFile(filename); err != nil {
		panic(err)
	}
	deploy3, err := handler.Get(data)
	checkErr("get deployment from bytes", deploy3.Name, err)

	// 4. get deployment from *appv1.Deployment.
	deploy4, err := handler.Get(deploy3)
	checkErr("get deployment from *appsv1.Deployment", deploy4.Name, err)

	// 5. get deployment from appsv1.Deployment
	deploy5, err := handler.Get(*deploy4)
	checkErr("get deployment from appsv1.Deployment", deploy5.Name, err)

	// 6. get deployment from runtime.Object.
	deploy6, err := handler.Get(runtime.Object(deploy5))
	checkErr("get deployment from runtime.Object", deploy6.Name, err)

	// 7. get deployment from *unstructured.Unstructured
	deploy7, err := handler.Get(&unstructured.Unstructured{Object: unstructData})
	checkErr("get deployment from *unstructured.Unstructured", deploy7.Name, err)

	// 8. get deployment from unstructured.Unstructured
	deploy8, err := handler.Get(unstructured.Unstructured{Object: unstructData})
	checkErr("get deployment from unstructured.Unstructured", deploy8.Name, err)

	// 9. get deployment from map[string]interface{}.
	deploy9, err := handler.Get(unstructData)
	checkErr("get deployment from map[string]interface{}", deploy9.Name, err)

	// Output:

	//2022/09/05 16:56:17 get deployment by name success: mydep
	//2022/09/05 16:56:17 get deployment from file success: mydep
	//2022/09/05 16:56:17 get deployment from bytes success: mydep
	//2022/09/05 16:56:17 get deployment from *appsv1.Deployment success: mydep
	//2022/09/05 16:56:17 get deployment from appsv1.Deployment success: mydep
	//2022/09/05 16:56:17 get deployment from runtime.Object success: mydep
	//2022/09/05 16:56:17 get deployment from *unstructured.Unstructured success: mydep-unstruct
	//2022/09/05 16:56:17 get deployment from unstructured.Unstructured success: mydep-unstruct
	//2022/09/05 16:56:17 get deployment from map[string]interface{} success: mydep-unstruct
}
