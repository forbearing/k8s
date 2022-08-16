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

`go get github.com/forbearing/k8s@v0.8.2`



## Planning

- [x] Simplify the use of client-go informer, lister
- [x] create/delete/update/delete/get ... all kinds of k8s resources by dyanmic client.
- [ ] Support crate/update/delete/get... `Event` resources
- [ ] Support crate/update/delete/get... `Endpoint` resources
- [ ] Support crate/update/delete/get... `EndpointSlice` resources
- [ ] Support crate/update/delete/get... `LimitRange` resources
- [ ] Support crate/update/delete/get... `PriorityClass` resources
- [ ] Support crate/update/delete/get... `ResourceQuota` resources
- [ ] Support crate/update/delete/get... `Lease` resources
- [x] Add function: GVK(), GVR(), Kind(), Group(), Version(), Resource(), KindToResource(), ResourceToKind()
- [x] signal handler
- [x] Finalizers
- [x] controller and owner
- [ ] UpdateStatus: update Deployment/StatefulSet... status
- [ ] UpdateScale: scale Deployment/StatefulSet...
- [ ] DeleteCollection
- [ ] Leader Election
- [ ] Recoder
- [ ] Replace interface{} -> any
- [ ] fack client
- [ ] EnvTest
- [ ] healthz
- [ ] metrics
- [ ] port-forward for pod, deployment and service
- [ ] proxy
- [ ] operators refer to https://sdk.operatorframework.io/docs/building-operators/golang/references/client/

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

	// Output:

	//2022/08/03 20:56:05 create deployment from file success: mydep
	//2022/08/03 20:56:05 create deployment from bytes success: mydep
	//2022/08/03 20:56:05 create deployment from *appsv1.Deployment success: mydep
	//2022/08/03 20:56:05 create deployment from appsv1.Deployment success: mydep
	//2022/08/03 20:56:05 create deployment from runtime.Object success: mydep
	//2022/08/03 20:56:05 create deployment from *unstructured.Unstructured success: mydep-unstruct
	//2022/08/03 20:56:06 create deployment from *unstructured.Unstructured success: mydep-unstruct
	//2022/08/03 20:56:06 create deployment from map[string]interface{} success: mydep-unstruct
	//2022/08/03 20:56:07 update deployment from file success: mydep
	//2022/08/03 20:56:07 update deployment from bytes success: mydep
	//2022/08/03 20:56:07 update deployment from *appsv1.Deployment success: mydep
	//2022/08/03 20:56:07 update deployment from appsv1.Deployment success: mydep
	//2022/08/03 20:56:07 update deployment from runtime.Object success: mydep
	//2022/08/03 20:56:07 update deployment from *unstructured.Unstructured success: mydep-unstruct
	//2022/08/03 20:56:07 update deployment from unstructured.Unstructured success: mydep-unstruct
	//2022/08/03 20:56:07 update deployment from map[string]interface{} success: mydep-unstruct
	//2022/08/03 20:56:08 apply deployment from file (deployment not exists) success: mydep
	//2022/08/03 20:56:08 apply deployment from file (deployment exists) success: mydep
	//2022/08/03 20:56:08 apply deployment from bytes (deployment not exists) success: mydep
	//2022/08/03 20:56:08 apply deployment from bytes (deployment exists) success: mydep
	//2022/08/03 20:56:08 apply deployment from *appsv1.Deployment (deployment not exists) success: mydep
	//2022/08/03 20:56:08 apply deployment from *appsv1.Deployment (deployment exists) success: mydep
	//2022/08/03 20:56:08 apply deployment from appsv1.Deployment (deployment not exists) success: mydep
	//2022/08/03 20:56:09 apply deployment from appsv1.Deployment (deployment exists) success: mydep
	//2022/08/03 20:56:09 apply deployment from runtime.Object (deployment not exists) success: mydep
	//2022/08/03 20:56:10 apply deployment from runtime.Object (deployment exists) success: mydep
	//2022/08/03 20:56:10 apply deployment from *unstructured.Unstructured (deployment not exists) success: mydep-unstruct
	//2022/08/03 20:56:10 apply deployment from *unstructured.Unstructured (deployment exists) success: mydep-unstruct
	//2022/08/03 20:56:11 apply deployment from unstructured.Unstructured (deployment not exists) success: mydep-unstruct
	//2022/08/03 20:56:11 apply deployment from unstructured.Unstructured (deployment exists) success: mydep-unstruct
	//2022/08/03 20:56:12 apply deployment from map[string]interface{} (deployment not exists) success: mydep-unstruct
	//2022/08/03 20:56:12 apply deployment from map[string]interface{} (deployment exists) success: mydep-unstruct
	//2022/08/03 20:56:13 delete deployment by name success:
	//2022/08/03 20:56:13 delete deployment from file success:
	//2022/08/03 20:56:13 delete deployment from bytes success:
	//2022/08/03 20:56:13 delete deployment from *appsv1.Deployment success:
	//2022/08/03 20:56:13 delete deployment from appsv1.Deployment success:
	//2022/08/03 20:56:13 delete deployment from runtime.Object success:
	//2022/08/03 20:56:14 delete deployment from *unstructured.Unstructured success:
	//2022/08/03 20:56:14 delete deployment from unstructured.Unstructured success:
	//2022/08/03 20:56:14 delete deployment from map[string]interface{} success:
	//2022/08/03 20:56:15 get deployment by name success: mydep
	//2022/08/03 20:56:15 get deployment from file success: mydep
	//2022/08/03 20:56:15 get deployment from bytes success: mydep
	//2022/08/03 20:56:15 get deployment from *appsv1.Deployment success: mydep
	//2022/08/03 20:56:15 get deployment from appsv1.Deployment success: mydep
	//2022/08/03 20:56:15 get deployment from runtime.Object success: mydep
	//2022/08/03 20:56:15 get deployment from *unstructured.Unstructured success: mydep-unstruct
	//2022/08/03 20:56:15 get deployment from unstructured.Unstructured success: mydep-unstruct
	//2022/08/03 20:56:15 get deployment from map[string]interface{} success: mydep-unstruct
	//2022/08/10 14:36:52 ListByLabel success: [coredns]
	//2022/08/10 14:36:52 ListByNamespace success: [calico-kube-controllers coredns metrics-server]
	//2022/08/10 14:36:53 ListAll success: [k8s-tools ......]
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

	// 6. create deployment from *unstructured.Unstructured
	deploy6, err := handler.Create(&unstructured.Unstructured{Object: unstructData})
	checkErr("create deployment from *unstructured.Unstructured", deploy6.Name, err)
	handler.Delete(unstructName)

	// 7. create deployment from unstructured.Unstructured
	deploy7, err := handler.Create(unstructured.Unstructured{Object: unstructData})
	checkErr("create deployment from *unstructured.Unstructured", deploy7.Name, err)
	handler.Delete(unstructName)

	// 8. create deployment from map[string]interface{}.
	deploy8, err := handler.Create(unstructData)
	checkErr("create deployment from map[string]interface{}", deploy8.Name, err)
	handler.Delete(unstructData)

	// Output:

	//2022/08/03 19:56:41 create deployment from file success: mydep
	//2022/08/03 19:56:41 create deployment from bytes success: mydep
	//2022/08/03 19:56:41 create deployment from *appsv1.Deployment success: mydep
	//2022/08/03 19:56:41 create deployment from appsv1.Deployment success: mydep
	//2022/08/03 19:56:41 create deployment from runtime.Object success: mydep
	//2022/08/03 19:56:41 create deployment from *unstructured.Unstructured success: mydep-unstruct
	//2022/08/03 19:56:42 create deployment from *unstructured.Unstructured success: mydep-unstruct
	//2022/08/03 19:56:42 create deployment from map[string]interface{} success: mydep-unstruct
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

	//2022/08/03 20:07:08 update deployment from file success: mydep
	//2022/08/03 20:07:08 update deployment from bytes success: mydep
	//2022/08/03 20:07:08 update deployment from *appsv1.Deployment success: mydep
	//2022/08/03 20:07:08 update deployment from appsv1.Deployment success: mydep
	//2022/08/03 20:07:08 update deployment from runtime.Object success: mydep
	//2022/08/03 20:07:08 update deployment from *unstructured.Unstructured success: mydep-unstruct
	//2022/08/03 20:07:08 update deployment from unstructured.Unstructured success: mydep-unstruct
	//2022/08/03 20:07:08 update deployment from map[string]interface{} success: mydep-unstruct
}
```



## How to apply deployment resources by handler.

```go
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

	//2022/08/03 20:24:08 apply deployment from file (deployment not exists) success: mydep
	//2022/08/03 20:24:08 apply deployment from file (deployment exists) success: mydep
	//2022/08/03 20:24:08 apply deployment from bytes (deployment not exists) success: mydep
	//2022/08/03 20:24:08 apply deployment from bytes (deployment exists) success: mydep
	//2022/08/03 20:24:08 apply deployment from *appsv1.Deployment (deployment not exists) success: mydep
	//2022/08/03 20:24:08 apply deployment from *appsv1.Deployment (deployment exists) success: mydep
	//2022/08/03 20:24:08 apply deployment from appsv1.Deployment (deployment not exists) success: mydep
	//2022/08/03 20:24:09 apply deployment from appsv1.Deployment (deployment exists) success: mydep
	//2022/08/03 20:24:09 apply deployment from runtime.Object (deployment not exists) success: mydep
	//2022/08/03 20:24:09 apply deployment from runtime.Object (deployment exists) success: mydep
	//2022/08/03 20:24:10 apply deployment from *unstructured.Unstructured (deployment not exists) success: mydep-unstruct
	//2022/08/03 20:24:10 apply deployment from *unstructured.Unstructured (deployment exists) success: mydep-unstruct
	//2022/08/03 20:24:11 apply deployment from unstructured.Unstructured (deployment not exists) success: mydep-unstruct
	//2022/08/03 20:24:11 apply deployment from unstructured.Unstructured (deployment exists) success: mydep-unstruct
	//2022/08/03 20:24:11 apply deployment from map[string]interface{} (deployment not exists) success: mydep-unstruct
	//2022/08/03 20:24:12 apply deployment from map[string]interface{} (deployment exists) success: mydep-unstruct
}
```



## How to delete deployment resources by handler.

```go
package main

import (
	"io/ioutil"

	"github.com/forbearing/k8s/deployment"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
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
	checkErr("delete deployment from runtime.Object", "", handler.Delete(runtime.Object(deploy)))

	// 7. delete deployment from *unstructured.Unstructured
	handler.Apply(unstructData)
	checkErr("delete deployment from *unstructured.Unstructured", "", handler.Delete(&unstructured.Unstructured{Object: unstructData}))

	// 8. delete deployment from unstructured.Unstructured
	handler.Apply(unstructData)
	checkErr("delete deployment from unstructured.Unstructured", "", handler.Delete(unstructured.Unstructured{Object: unstructData}))

	// 9. delete deployment from map[string]interface{}.
	handler.Apply(unstructData)
	checkErr("delete deployment from map[string]interface{}", "", handler.Delete(unstructData))

	// Output:

	//2022/08/03 20:55:20 delete deployment by name success:
	//2022/08/03 20:55:20 delete deployment from file success:
	//2022/08/03 20:55:20 delete deployment from bytes success:
	//2022/08/03 20:55:20 delete deployment from *appsv1.Deployment success:
	//2022/08/03 20:55:20 delete deployment from appsv1.Deployment success:
	//2022/08/03 20:55:20 delete deployment from runtime.Object success:
	//2022/08/03 20:55:20 delete deployment from *unstructured.Unstructured success:
	//2022/08/03 20:55:21 delete deployment from unstructured.Unstructured success:
	//2022/08/03 20:55:21 delete deployment from map[string]interface{} success:
}
```



## How to get deployment resources by handler.

```go
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

	//2022/08/03 20:43:38 get deployment by name success: mydep
	//2022/08/03 20:43:38 get deployment from file success: mydep
	//2022/08/03 20:43:38 get deployment from bytes success: mydep
	//2022/08/03 20:43:38 get deployment from *appsv1.Deployment success: mydep
	//2022/08/03 20:43:38 get deployment from appsv1.Deployment success: mydep
	//2022/08/03 20:43:38 get deployment from runtime.Object success: mydep
	//2022/08/03 20:43:38 get deployment from *unstructured.Unstructured success: mydep-unstruct
	//2022/08/03 20:43:38 get deployment from unstructured.Unstructured success: mydep-unstruct
	//2022/08/03 20:43:38 get deployment from map[string]interface{} success: mydep-unstruct
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

	//2022/08/10 14:36:52 ListByLabel success: [coredns]
	//2022/08/10 14:36:52 List success: [coredns]
	//2022/08/10 14:36:52 ListByNamespace success: [calico-kube-controllers coredns metrics-server]
	//2022/08/10 14:36:53 ListAll success: [k8s-tools ......]
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
