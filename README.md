## Introduction



The library implements various handlers to more easy manipulate k8s resources such as pods, deployments, etc, inside or outside k8s cluster. A program that uses the library and runs in a k8s pod meant to be inside k8s cluster. If you simply run [examples](./examples) in your pc/mac or server, it meant outside k8s cluster. Both of inside and outside k8s cluster is supported by the library.

To create a handler for outside cluster just call `deployment.New(ctx, kubeconfig, namespace)`.
To create a handler for the inside cluster just call `deployment.New(ctx, "", namespace)`.

The variable `namespace` is used to limit the scope of the handler. If `namespace=test`, the handler is only allowed to create/update/delete deployments in namespace/test. Of course, handler.WithNamespace(newNamespace) returns a new temporary handler that allowed to create/update/delete deployments in the new namespace, for examples:

```go
namespace1 := "test"
namespace2 := "test2"
// Inside cluster. the program run within k8s pod.
handler, _ := deployment.New(ctx, "", namespace1)
// handler is only allowed to create/update/delete deployment in namespace/test.
handler.Create(filename)
// handler is only allowed to create/update/delete deployment in namespace/test2.
handler.WithNamespace(namespace2).Create(filename)
// handler is only allowed to create/update/delete deployment in namespace/test (not namespace/test2).
handler.Create(filename)
```

The library is used by another open source project that used to backup pv/pvc data attached by deployments/statefulsets/daemosnets/pods running in k8s cluster.

For more examples on how to use this library, you can refer to the [examples](./examples) folder or related test code.

## Installation

```go get github.com/forbearing/k8s@v0.4.4```



## Planning

- [ ] Simplify the use of client-go informer, lister
- [ ] create/delete/update/delete/get ... all kinds of k8s resources by dyanmic client.
- [ ] Support crate/update/delete/get... `Event` resources
- [ ] Support crate/update/delete/get... `Endpoint` resources
- [ ] Support crate/update/delete/get... `EndpointSlice` resources
- [ ] Support crate/update/delete/get... `LimitRange` resources
- [ ] Support crate/update/delete/get... `PriorityClass` resources
- [ ] Support crate/update/delete/get... `ResourceQuota` resources
- [ ] Add function: GVK() returns the Group, version, Resource name of k8s resources.

## How to execute command within pod by handler.

```golang
import (
	"github.com/forbearing/k8s/pod"
)

var (
	ctx, cancel = context.WithTimeout(context.Background(), time.Minute*10)
	namespace   = "test"
	kubeconfig  = filepath.Join(os.Getenv("HOME"), ".kube/config")
	filename    = "../../testData/examples/pod.yaml"
	name        = "mypod"
	label       = "type=pod"
)

func main() {
	defer cancel()

	Pod_Tools()
}

func cleanup(handler *pod.Handler) {
	handler.Delete(name)
}

func Pod_Tools() {
	handler, err := pod.New(ctx, kubeconfig, namespace)
	if err != nil {
		panic(err)
	}
	defer cleanup(handler)
	handler.Apply(filename)
	handler.WaitReady(name)

	command1 := []string{
		"hostname",
	}
	command2 := []string{
		"sh",
		"-c",
		"hostname",
	}
	command3 := []string{
		"/bin/sh",
		"-c",
		"hostname",
	}
	command4 := []string{
		"/bin/bash",
		"-c",
		"hostname",
	}
	command5 := []string{
		"cat /etc/os-release",
	}
	command6 := []string{
		"sh",
		"-c",
		"cat /etc/os-release",
	}
	command7 := []string{
		"sh",
		"-c",
		"apt update; apt upgrade -y",
	}
	handler.Execute(name, "", command1, nil) // execute success.
	handler.Execute(name, "", command2, nil) // execute success.
	handler.Execute(name, "", command3, nil) // execute success.
	handler.Execute(name, "", command4, nil) // execute success.
	handler.Execute(name, "", command5, nil) // execute failed.
	handler.Execute(name, "", command6, nil) // execute success.
	handler.Execute(name, "", command7, nil) // execute success, but may be cancelled by context timeout.

	// Output:

	//mypod
	//mypod
	//mypod
	//mypod
	//OCI runtime exec failed: exec failed: unable to start container process: exec: "cat /etc/os-release": stat cat /etc/os-release: no such file or directory: unknown
	//PRETTY_NAME="Debian GNU/Linux 11 (bullseye)"
	//NAME="Debian GNU/Linux"
	//VERSION_ID="11"
	//VERSION="11 (bullseye)"
	//VERSION_CODENAME=bullseye
	//ID=debian
	//HOME_URL="https://www.debian.org/"
	//SUPPORT_URL="https://www.debian.org/support"
	//BUG_REPORT_URL="https://bugs.debian.org/"
	//Get:1 http://security.debian.org/debian-security bullseye-security InRelease [44.1 kB]
	//Get:2 http://deb.debian.org/debian bullseye InRelease [116 kB]
	//Get:3 http://security.debian.org/debian-security bullseye-security/main amd64 Packages [163 kB]
	//Get:4 http://deb.debian.org/debian bullseye-updates InRelease [39.4 kB]
	//Get:5 http://deb.debian.org/debian bullseye/main amd64 Packages [8182 kB]
	//26% [5 Packages 750 kB/8182 kB 9%]                           36.7 kB/s 3min 22s
}
```



## How to create deployment resources by handler.

```golang
package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/forbearing/k8s"
	"github.com/forbearing/k8s/deployment"
)

var (
	ctx, cancel = context.WithTimeout(context.Background(), time.Minute*10)
	namespace   = "test"
	kubeconfig  = filepath.Join(os.Getenv("HOME"), ".kube/config")
	filename    = "../../testdata/examples/deployment.yaml"
	updateFile  = "../../testdata/examples/deployment-update1.yaml"
	filename2   = "../../testdata/nginx/nginx-deploy.yaml"
	name        = "mydep"
	name2       = "nginx-deploy"
	label       = "type=deployment"
)

var (
	unstructName = "mydep-unstruct"
	unstructData = map[string]interface{}{
		"apiVersion": "apps/v1",
		"kind":       "Deployment",
		"metadata": map[string]interface{}{
			"name": unstructName,
			"labels": map[string]interface{}{
				"app":  unstructName,
				"type": "deployment",
			},
		},
		"spec": map[string]interface{}{
			// replicas type is int32, not string.
			"replicas": 1,
			"selector": map[string]interface{}{
				"matchLabels": map[string]interface{}{
					"app":  unstructName,
					"type": "deployment",
				},
			},
			"template": map[string]interface{}{
				"metadata": map[string]interface{}{
					"labels": map[string]interface{}{
						"app":  unstructName,
						"type": "deployment",
					},
				},
				"spec": map[string]interface{}{
					"containers": []map[string]interface{}{
						{
							"name":  "nginx",
							"image": "nginx",
							"resources": map[string]interface{}{
								"limits": map[string]interface{}{
									"cpu": "100m",
								},
							},
						},
					},
				},
			},
		},
	}
)

func main() {
	Deployment_Create()
	Deployment_Update()
	Deployment_Apply()
	Deployment_Delete()
	Deployment_Get()
	Deployment_List()
	//Deployment_Watch()
	//Deployment_Informer()
	//Deployment_Tools()

	// Output:

	//2022/07/22 22:52:47 create deployment from file success:
	//2022/07/22 22:52:47 create deployment from bytes success:
	//2022/07/22 22:52:47 create deployment from *appsv1.Deployment success:
	//2022/07/22 22:52:47 create deployment from appsv1.Deployment success:
	//2022/07/22 22:52:47 create deployment from runtime.Object success:
	//2022/07/22 22:52:48 create deployment from unstructured data success:
	//2022/07/22 22:52:48 update deployment from file success:
	//2022/07/22 22:52:48 update deployment from bytes success:
	//2022/07/22 22:52:48 update deployment from *appsv1.Deployment success:
	//2022/07/22 22:52:48 update deployment from appsv1.Deployment success:
	//2022/07/22 22:52:48 update deployment from runtime.Object success:
	//2022/07/22 22:52:48 update deployment from unstructured data success:
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
	//2022/07/22 22:52:52 delete deployment by name success:
	//2022/07/22 22:52:52 delete deployment from file success:
	//2022/07/22 22:52:52 delete deployment from bytes success:
	//2022/07/22 22:52:52 delete deployment from *appsv1.Deployment success:
	//2022/07/22 22:52:52 delete deployment from appsv1.Deployment success:
	//2022/07/22 22:52:53 delete deployment from runtime.Object success:
	//2022/07/22 22:52:53 delete deployment from unstructured data success:
	//2022/07/22 22:52:54 get deployment by name success: mydep
	//2022/07/22 22:52:54 get deployment from file success: mydep
	//2022/07/22 22:52:54 get deployment from bytes success: mydep
	//2022/07/22 22:52:54 get deployment from *appsv1.Deployment success: mydep
	//2022/07/22 22:52:54 get deployment from appsv1.Deployment success: mydep
	//2022/07/22 22:52:54 get deployment from runtime.Object success: mydep
	//2022/07/22 22:52:54 get deployment from unstructured data success: mydep-unstruct
	//2022/07/22 22:52:54 ListByLabel success:
	//2022/07/22 22:52:54 [nginx-deploy]
	//2022/07/22 22:52:54 List success:
	//2022/07/22 22:52:54 [nginx-deploy]
	//2022/07/22 22:52:54 ListByNamespace success:
	//2022/07/22 22:52:54 [nginx-deploy]
	//2022/07/22 22:52:54 ListAll success:
	//2022/07/22 22:52:54 [horus-operator ingress-controller calico-kube-controllers coredns metrics-server dashboard-metrics-scraper kubernetes-dashboard local-path-provisioner nfs-provisioner-nfs-subdir-external-provisioner nginx-deploy]
}

func checkErr(name string, val interface{}, err error) {
	if err != nil {
		log.Printf("%s failed: %v\n", name, err)
	} else {
		log.Printf("%s success: %v\n", name, val)
	}
}

// cleanup will delete or prune created deployments.
func cleanup(handler *deployment.Handler) {
	handler.Delete(name)
	handler.Delete(unstructName)
	handler.DeleteFromFile(updateFile)
	k8s.DeleteF(ctx, kubeconfig, filename2)
}

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
```



## How to update deployment resources by handler.

```golang
package main

import (
	"io/ioutil"

	"github.com/forbearing/k8s/deployment"
	"k8s.io/apimachinery/pkg/runtime"
)

func Deployment_Update() {
	// New returns a handler used to multiples deployment.
	handler, err := deployment.New(ctx, kubeconfig, namespace)
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

	//2022/07/22 23:09:31 update deployment from file success:
	//2022/07/22 23:09:31 update deployment from bytes success:
	//2022/07/22 23:09:31 update deployment from *appsv1.Deployment success:
	//2022/07/22 23:09:31 update deployment from appsv1.Deployment success:
	//2022/07/22 23:09:31 update deployment from runtime.Object success:
	//2022/07/22 23:09:31 update deployment from unstructured data success:
}
```



## How to apply deployment resources by handler.

```go
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
	handler, err := deployment.New(ctx, kubeconfig, namespace)
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
```



## How to delete deployment resources by handler.

```go
package main

import (
	"io/ioutil"

	"github.com/forbearing/k8s/deployment"
	"k8s.io/apimachinery/pkg/runtime"
)

func Deployment_Delete() {
	handler, err := deployment.New(ctx, kubeconfig, namespace)
	if err != nil {
		panic(err)
	}
	defer cleanup(handler)

	// 1. delete deployment by name.
	handler.Apply(filename)
	checkErr("delete deployment by name", "", handler.Delete(name))

	// 2. delete deployment from file.
	// You should always explicitly call DeleteFromFile to delete a deployment
	// from file. if the parameter type passed to the `Delete` method is string,
	// the `Delete` method will call DeleteByName not DeleteFromFile.
	handler.Apply(filename)
	checkErr("delete deployment from file", "", handler.DeleteFromFile(filename))

	// 3. delete deployment from bytes.
	var data []byte
	if data, err = ioutil.ReadFile(filename); err != nil {
		panic(err)
	}
	handler.Apply(filename)
	checkErr("delete deployment from bytes", "", handler.Delete(data))

	// 4. delete deployment from *appsv1.deployment
	deploy, _ := handler.Apply(filename)
	checkErr("delete deployment from *appsv1.Deployment", "", handler.Delete(deploy))

	// 5. delete deployment from appsv1.deployment
	deploy, _ = handler.Apply(filename)
	checkErr("delete deployment from appsv1.Deployment", "", handler.Delete(*deploy))

	// 6. delete deployment from runtime.Object.
	deploy, _ = handler.Apply(filename)
	object := runtime.Object(deploy)
	checkErr("delete deployment from runtime.Object", "", handler.Delete(object))

	// 7. delete deployment from unstructured data, aka map[string]interface{}.
	handler.Apply(unstructData)
	checkErr("delete deployment from unstructured data", "", handler.Delete(unstructData))

	// Output:

	//2022/07/22 22:20:18 delete deployment by name success:
	//2022/07/22 22:20:18 delete deployment from file success:
	//2022/07/22 22:20:18 delete deployment from bytes success:
	//2022/07/22 22:20:18 delete deployment from *appsv1.Deployment success:
	//2022/07/22 22:20:18 delete deployment from appsv1.Deployment success:
	//2022/07/22 22:20:18 delete deployment from runtime.Object success:
	//2022/07/22 22:20:19 delete deployment from unstructured data success:
}
```



## How to get deployment resources by handler.

```go
package main

import (
	"io/ioutil"

	"github.com/forbearing/k8s/deployment"
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
	object := runtime.Object(deploy5)
	deploy6, err := handler.Get(object)
	checkErr("get deployment from runtime.Object", deploy6.Name, err)

	// 7. get deployment from unstructured data, aka map[string]interface{}.
	deploy7, err := handler.Get(unstructData)
	checkErr("get deployment from unstructured data", deploy7.Name, err)

	// Output:

	//2022/07/22 22:47:33 get deployment by name success: mydep
	//2022/07/22 22:47:33 get deployment from file success: mydep
	//2022/07/22 22:47:33 get deployment from bytes success: mydep
	//2022/07/22 22:47:33 get deployment from *appsv1.Deployment success: mydep
	//2022/07/22 22:47:33 get deployment from appsv1.Deployment success: mydep
	//2022/07/22 22:47:33 get deployment from runtime.Object success: mydep
	//2022/07/22 22:47:33 get deployment from unstructured data success: mydep-unstruct
}
```



## How to list deployment resources by handler.

```go
package main

import (
	"log"

	"github.com/forbearing/k8s"
	"github.com/forbearing/k8s/deployment"
	appsv1 "k8s.io/api/apps/v1"
)

func Deployment_List() {
	// New returns a handler used to multiples deployment.
	handler, err := deployment.New(ctx, kubeconfig, namespace)
	if err != nil {
		panic(err)
	}
	defer cleanup(handler)

	k8s.ApplyF(ctx, kubeconfig, filename2)

	// ListByLabel list deployment by label.
	deployList1, err := handler.ListByLabel(label)
	checkErr("ListByLabel", "", err)
	outputDeploy(*deployList1)

	// List list deployment by label, it's alias to "ListByLabel".
	deployList2, err := handler.List(label)
	checkErr("List", "", err)
	outputDeploy(*deployList2)

	// ListByNamespace list all deployments in the namespace where the deployment is running.
	deployList3, err := handler.ListByNamespace(namespace)
	checkErr("ListByNamespace", "", err)
	outputDeploy(*deployList3)

	// ListAll list all deployments in the k8s cluster.
	deployList4, err := handler.ListAll()
	checkErr("ListAll", "", err)
	outputDeploy(*deployList4)

	// Output:

	//2022/07/04 21:43:09 ListByLabel success.
	//2022/07/04 21:43:09 [mydep-2 nginx-deploy]
	//2022/07/04 21:43:09 List success.
	//2022/07/04 21:43:09 [mydep-2 nginx-deploy]
	//2022/07/04 21:43:09 ListByNamespace success.
	//2022/07/04 21:43:09 [mydep-2 nginx-deploy]
	//2022/07/04 21:43:09 ListAll success.
	//2022/07/04 21:43:09 [calico-kube-controllers coredns metrics-server local-path-provisioner nfs-provisioner-nfs-subdir-external-provisioner mydep-2 nginx-deploy]

}

func outputDeploy(deployList appsv1.DeploymentList) {
	var dl []string
	for _, deploy := range deployList.Items {
		dl = append(dl, deploy.Name)
	}
	log.Println(dl)
}
```



## How to watch deployment resources by handler.

```go
package main

import (
	"context"
	"io/ioutil"
	"log"
	"time"

	"github.com/forbearing/k8s/deployment"
)

func Deployment_Watch() {
	// New returns a handler used to multiples deployment.
	handler, err := deployment.New(ctx, kubeconfig, namespace)
	if err != nil {
		panic(err)
	}
	defer cleanup(handler)
	handler.Apply(filename)

	addFunc := func(x interface{}) { log.Println("added deployment.") }
	modifyFunc := func(x interface{}) { log.Println("modified deployment.") }
	deleteFunc := func(x interface{}) { log.Println("deleted deployment.") }

	// WatchByLabel watchs a set of deployments by label.
	{
		ctx, cancel := context.WithCancel(ctx)

		go func(ctx context.Context) {
			handler.WatchByLabel(label, addFunc, modifyFunc, deleteFunc, nil)
		}(ctx)
		go func(ctx context.Context) {
			for {
				handler.Apply(filename)
				time.Sleep(time.Second * 5)
				handler.Delete(name)
			}
		}(ctx)

		timer := time.NewTimer(time.Second * 30)
		<-timer.C
		cancel()
	}

	// WatchByName watchs a deployment by label.
	// Watch simply calls WatchByName.
	ioutil.ReadFile(filename)
	{
		ctx, cancel := context.WithCancel(ctx)

		go func(ctx context.Context) {
			handler.WatchByName(name, addFunc, modifyFunc, deleteFunc, nil)
		}(ctx)
		go func(ctx context.Context) {
			for {
				handler.Apply(filename)
				time.Sleep(time.Second * 5)
				handler.Delete(name)
			}
		}(ctx)

		timer := time.NewTimer(time.Second * 30)
		<-timer.C
		cancel()
	}

	// Output:

	//2022/07/04 21:59:58 modified deployment.
	//2022/07/04 21:59:58 modified deployment.
	//2022/07/04 21:59:58 modified deployment.
	//2022/07/04 21:59:58 modified deployment.
	//2022/07/04 21:59:58 modified deployment.
	//2022/07/04 22:00:03 deleted deployment.
	//2022/07/04 22:00:03 added deployment.
	//2022/07/04 22:00:03 modified deployment.
	//2022/07/04 22:00:03 modified deployment.
	//2022/07/04 22:00:03 modified deployment.
	//2022/07/04 22:00:08 deleted deployment.
	//2022/07/04 22:00:08 added deployment.
	//2022/07/04 22:00:08 modified deployment.
	//2022/07/04 22:00:08 modified deployment.
	//2022/07/04 22:00:08 modified deployment.
	//2022/07/04 22:00:08 modified deployment.
	//2022/07/04 22:00:13 deleted deployment.
	//2022/07/04 22:00:13 added deployment.
	//2022/07/04 22:00:13 modified deployment.
	//2022/07/04 22:00:13 modified deployment.
	//2022/07/04 22:00:13 modified deployment.
	//2022/07/04 22:00:18 deleted deployment.
	//2022/07/04 22:00:18 added deployment.
	//2022/07/04 22:00:18 modified deployment.
	//2022/07/04 22:00:18 modified deployment.
	//2022/07/04 22:00:18 modified deployment.
	//2022/07/04 22:00:23 deleted deployment.
	//2022/07/04 22:00:23 added deployment.
	//2022/07/04 22:00:23 modified deployment.
	//2022/07/04 22:00:23 modified deployment.
	//2022/07/04 22:00:23 modified deployment.
	//2022/07/04 22:00:28 deleted deployment.
}
```



## How to get deployment resource more info.

```go
package main

import (
	"log"
	"time"

	"github.com/forbearing/k8s"
	"github.com/forbearing/k8s/deployment"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

func Deployment_Tools() {
	handler, err := deployment.New(ctx, kubeconfig, namespace)
	if err != nil {
		panic(err)
	}
	defer cleanup(handler)

	// kubectl apply -f ../../testdata/nginx/nginx-deploy.yaml
	k8s.ApplyF(ctx, kubeconfig, filename2)

	log.Println(handler.IsReady(name2)) // false
	handler.WaitReady(name2)            // block until the deployment is ready and available.
	log.Println(handler.IsReady(name2)) // true

	deploy, err := handler.Get(name2)
	if err != nil {
		panic(err)
	}

	getByName := func() {
		log.Println("===== Get Deployment By Name")
		begin := time.Now()

		// GetRS get all replicaset that created by the deployment.
		rsList, err := handler.GetRS(name2)
		checkErr("GetRS", printRS(rsList), err)
		// GetPods get all pods that created by the deployment.
		podList, err := handler.GetPods(name2)
		checkErr("GetRS", printPods(podList), err)
		// GetPV get all persistentvolume that attached by the deployment.
		pvList, err := handler.GetPV(name2)
		checkErr("GetPV", pvList, err)
		// GetPVC get all persistentvolumeclaim that attached by the deployment.
		pvcList, err := handler.GetPVC(name2)
		checkErr("GetPVC", pvcList, err)

		end := time.Now()
		log.Println("===== Get Deployment By Name Cost Time:", end.Sub(begin))
		log.Println()
	}

	getByObj := func() {
		log.Println("===== Get Deployment By Object")
		begin := time.Now()

		// GetRS get all replicaset that created by the deployment.
		rsList, err := handler.GetRS(deploy)
		checkErr("GetRS", printRS(rsList), err)
		// GetPods get all pods that created by the deployment.
		podList, err := handler.GetPods(deploy)
		checkErr("GetRS", printPods(podList), err)
		// GetPV get all persistentvolume that attached by the deployment.
		pvList, err := handler.GetPV(deploy)
		checkErr("GetPV", pvList, err)
		// GetPVC get all persistentvolumeclaim that attached by the deployment.
		pvcList, err := handler.GetPVC(deploy)
		checkErr("GetPVC", pvcList, err)

		end := time.Now()
		log.Println("===== Get Deployment By Object Cost Time:", end.Sub(begin))
	}

	getByName()
	getByObj()

	// Output:

	//2022/07/12 09:20:45 false
	//2022/07/12 09:22:16 true
	//2022/07/12 09:22:16 ===== Get Deployment By Name
	//2022/07/12 09:22:16 GetRS success: [nginx-deploy-79979d95dd]
	//2022/07/12 09:22:16 GetRS success: [nginx-deploy-79979d95dd-4lm8h nginx-deploy-79979d95dd-5l9rk nginx-deploy-79979d95dd-scjw9]
	//2022/07/12 09:22:16 GetPV success: [pvc-93ebe9a0-c464-439b-a252-51afb4d87069 pvc-c048ccf9-4d0c-4312-bb36-15e4fa7a1746 pvc-dc16fea0-f511-42d7-b78e-6fcac96fcc9b]
	//2022/07/12 09:22:16 GetPVC success: [deploy-k8s-tools-data deploy-nginx-data deploy-nginx-html]
	//2022/07/12 09:22:16 ===== Get Deployment By Name Cost Time: 82.467272ms
	//2022/07/12 09:22:16
	//2022/07/12 09:22:16 ===== Get Deployment By Object
	//2022/07/12 09:22:16 GetRS success: [nginx-deploy-79979d95dd]
	//2022/07/12 09:22:16 GetRS success: [nginx-deploy-79979d95dd-4lm8h nginx-deploy-79979d95dd-5l9rk nginx-deploy-79979d95dd-scjw9]
	//2022/07/12 09:22:17 GetPV success: [pvc-93ebe9a0-c464-439b-a252-51afb4d87069 pvc-c048ccf9-4d0c-4312-bb36-15e4fa7a1746 pvc-dc16fea0-f511-42d7-b78e-6fcac96fcc9b]
	//2022/07/12 09:22:17 GetPVC success: [deploy-k8s-tools-data deploy-nginx-data deploy-nginx-html]
	//2022/07/12 09:22:17 ===== Get Deployment By Object Cost Time: 134.639944ms
}

func printPods(podList []corev1.Pod) []string {
	var pl []string
	for _, pod := range podList {
		pl = append(pl, pod.Name)
	}
	return pl
}
func printRS(rsList []appsv1.ReplicaSet) []string {
	var rl []string
	for _, rs := range rsList {
		rl = append(rl, rs.Name)
	}
	return rl
}
```
