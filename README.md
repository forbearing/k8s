## Introduction



The library implements various handlers to more easy manipulate k8s resources such as pods, deployments, etc, inside or outside k8s cluster. A program that uses the library and runs in a k8s pod meant to be inside k8s cluster. If you simply run [test code](./testCode) in your pc/mac or server, it meant outside k8s cluster. Both of inside and outside k8s cluster is supported by the library.

To create a handler for outside cluster just call `deployment.New(ctx, namespace, kubeconfig)`.
To create a handler for the inside cluster just call `deployment.New(ctx, namespace, "")`.

The Variable `namespace` is used to limit the scope of the handler. If `namespace=test`, the handler is only allowed to create/update/delete deployments in namespace/test. Of course, handler.WithNamespace(newNamespace) returns a new temporary handler that allowed to create/update/delete deployments in the new namespace, for examples:

```go
namespace1 := "test"
namespace2 := "test2"
// Inside cluster. the program run within k8s pod.
handler, _ := deployment.New(ctx, namespace1, "")
// handler is only allowed to create/update/delete deployment in namespace/test.
handler.Create(filename)
// handler is only allowed to create/update/delete deployment in namespace/test2.
handler.WithNamespace(namespace2).Create(filename)
// handler is only allowed to create/update/delete deployment in namespace/test (not namespace/test2).
handler.Create(filename)
```



For more examples on how to use this library, you can refer to the [testCode](testCode) folder or related test code.

## Installation

```go get -u github.com/forbearing/k8s```



## Planning

- Simplify the use of client-go informer, lister
- 完成测试代码的编写和使用示例, 使用示例放在 testCode 中
- 代码优化

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
	handler, err := pod.New(ctx, namespace, kubeconfig)
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
	handler.Execute(name, "", command1) // execute success.
	handler.Execute(name, "", command2) // execute success.
	handler.Execute(name, "", command3) // execute success.
	handler.Execute(name, "", command4) // execute success.
	handler.Execute(name, "", command5) // execute failed.
	handler.Execute(name, "", command6) // execute success.
	handler.Execute(name, "", command7) // execute success, but may be cancelled by context timeout.

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



## How to create k8s deployment by handler.

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
	corev1 "k8s.io/api/core/v1"
)

var (
	ctx, cancel = context.WithTimeout(context.Background(), time.Minute*10)
	namespace   = "test"
	kubeconfig  = filepath.Join(os.Getenv("HOME"), ".kube/config")
	filename    = "../../testData/examples/deployment.yaml"
	filename2   = "../../testData/examples/deployment-2.yaml"
	update1File = "../../testData/examples/deployment-update1.yaml"
	update2File = "../../testData/examples/deployment-update2.yaml"
	update3File = "../../testData/examples/deployment-update3.yaml"
	nginxFile   = "../../testData/nginx/nginx-deploy.yaml"
	name        = "mydep"
	label       = "type=deployment"
)

var (
	rawName = "mydep-raw"
	rawData = map[string]interface{}{
		"apiVersion": "apps/v1",
		"kind":       "Deployment",
		"metadata": map[string]interface{}{
			"name": rawName,
			"labels": map[string]interface{}{
				"app":  rawName,
				"type": "deployment",
			},
		},
		"spec": map[string]interface{}{
			// replicas type is int32, not string.
			"replicas": 1,
			"selector": map[string]interface{}{
				"matchLabels": map[string]interface{}{
					"app":  rawName,
					"type": "deployment",
				},
			},
			"template": map[string]interface{}{
				"metadata": map[string]interface{}{
					"labels": map[string]interface{}{
						"app":  rawName,
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
	//Deployment_Create()
	//Deployment_Update()
	//Deployment_Apply()
	//Deployment_Delete()
	//Deployment_Get()
	//Deployment_List()
	//Deployment_Watch()
	Deployment_Tools()
}

func myerr(name string, err error) {
	if err != nil {
		log.Printf("%s failed: %v\n", name, err)
	} else {
		log.Printf("%s success.\n", name)
	}
}

// cleanup will delete or prune created deployments.
func cleanup(handler *deployment.Handler) {
	handler.Delete(name)
	handler.Delete(rawName)
	handler.DeleteFromFile(filename2)
	handler.DeleteFromFile(update1File)
	handler.DeleteFromFile(update2File)
	handler.DeleteFromFile(update3File)
	//k8s.DeleteF(ctx, kubeconfig, nginxFile)
}

func Deployment_Create() {
	// New returns a handler used to multiples deployment.
	handler, err := deployment.New(ctx, namespace, kubeconfig)
	if err != nil {
		panic(err)
	}
	defer cleanup(handler)

	// CreateFromRaw creates a deploymnt from map[string]interface.
	_, err = handler.CreateFromRaw(rawData)
	myerr("CreateFromRaw", err)
	handler.Delete(name)

	// CreateFromFile creates a deploymnt from file.
	_, err = handler.CreateFromFile(filename)
	myerr("CreateFromFile", err)
	handler.Delete(name)

	// CreateFromBytes creates a deploymnt from bytes.
	var data []byte
	if data, err = ioutil.ReadFile(filename); err != nil {
		panic(err)
	}
	_, err = handler.CreateFromBytes(data)
	myerr("CreateFromBytes", err)
	handler.Delete(name)

	// Create creates a deploymnt from file, it's alias to "CreateFromFile".
	_, err = handler.Create(filename)
	myerr("Create", err)

	// Output:

	//2022/07/04 21:43:04 CreateFromRaw success.
	//2022/07/04 21:43:04 CreateFromFile success.
	//2022/07/04 21:43:04 CreateFromBytes success.
	//2022/07/04 21:43:04 Create success.
}
```



## How to update k8s deployment by handler.

```golang
func Deployment_Update() {
	// New returns a handler used to multiples deployment.
	handler, err := deployment.New(ctx, namespace, kubeconfig)
	if err != nil {
		panic(err)
	}
	defer cleanup(handler)
	handler.Apply(filename)
	handler.ApplyFromRaw(rawData)

	// UpdateFromRaw updates deployment from map[string]interface.
	_, err = handler.UpdateFromRaw(rawData)
	myerr("UpdateFromRaw", err)

	// UpdateFromFile updates deployment from file.
	_, err = handler.UpdateFromFile(update1File)
	myerr("UpdateFromFile", err)

	// UpdateFromBytes updates deployment from bytes.
	var data []byte
	if data, err = ioutil.ReadFile(update2File); err != nil {
		panic(err)
	}
	_, err = handler.UpdateFromBytes(data)
	myerr("UpdateFromBytes", err)

	// Update updates deployment from file, it's alias to "UpdateFromFile".
	_, err = handler.Update(update3File)
	myerr("Update", err)

	// Output:

	//2022/07/04 21:43:05 UpdateFromRaw success.
	//2022/07/04 21:43:05 UpdateFromFile success.
	//2022/07/04 21:43:05 UpdateFromBytes success.
	//2022/07/04 21:43:05 Update success.
}
```



## How to apply k8s deployment by handler.

```go
// ApplyXXX method will creates deployment if specified deployment not exist and
// it will updates deployment if already exist.
func Deployment_Apply() {
	// New returns a handler used to multiples deployment.
	handler, err := deployment.New(ctx, namespace, kubeconfig)
	if err != nil {
		panic(err)
	}
	defer cleanup(handler)

	handler.Delete(name)
	handler.Delete(rawName)

	// ApplyFromRaw apply a deployment from map[string]interface.
	// it will updates the deployment, if already exist, or creates it.
	handler.CreateFromRaw(rawData)
	_, err = handler.ApplyFromRaw(rawData)
	myerr("ApplyFromRaw", err)

	handler.Delete(rawName)
	_, err = handler.ApplyFromRaw(rawData)
	myerr("ApplyFromRaw", err)

	// ApplyFromFile apply a deployment from file.
	// it will updates the deployment, if already exist, or creates it.
	handler.CreateFromFile(update1File)
	_, err = handler.ApplyFromFile(update1File)
	myerr("ApplyFromFile", err)
	handler.DeleteFromFile(update1File)
	_, err = handler.ApplyFromFile(update1File)
	myerr("ApplyFromFile", err)

	// ApplyFromBytes apply a deployment from bytes.
	// it will updates the deployment, if already exist, or creates it.
	var data []byte
	if data, err = ioutil.ReadFile(update2File); err != nil {
		panic(err)
	}
	handler.CreateFromFile(update2File)
	_, err = handler.ApplyFromBytes(data)
	myerr("ApplyFromBytes", err)
	handler.DeleteFromFile(update2File)
	_, err = handler.ApplyFromBytes(data)
	myerr("ApplyFromBytes", err)

	// Apply apply a deployment from file, it's alias to "ApplyFromFile".
	// it will updates the deployment, if already exist, or creates it.
	handler.CreateFromFile(update3File)
	_, err = handler.Apply(update3File)
	myerr("Apply", err)
	handler.DeleteFromFile(update3File)
	myerr("Apply", err)

	// Output:

	//2022/07/04 21:43:05 ApplyFromRaw success.
	//2022/07/04 21:43:05 ApplyFromRaw success.
	//2022/07/04 21:43:05 ApplyFromFile success.
	//2022/07/04 21:43:05 ApplyFromFile success.
	//2022/07/04 21:43:06 ApplyFromBytes success.
	//2022/07/04 21:43:06 ApplyFromBytes success.
	//2022/07/04 21:43:07 Apply success.
	//2022/07/04 21:43:07 Apply success.
}
```



## How to delete k8s deployment by handler.

```go
func Deployment_Delete() {
	handler, err := deployment.New(ctx, namespace, kubeconfig)
	if err != nil {
		panic(err)
	}
	defer cleanup(handler)

	// DeleteByName delete a deployment by name.
	handler.Apply(filename)
	myerr("DeleteByName", handler.DeleteByName(name))

	// Delete delete a deployment by name, it's alias to "DeleteByName".
	handler.Apply(filename)
	myerr("Delete", handler.Delete(name))

	// DeleteFromFile delete a deployment from yaml file.
	handler.Apply(filename)
	myerr("DeleteFromFile", handler.DeleteFromFile(filename))

	// DeleteFromBytes delete a deployment from bytes.
	var data []byte
	if data, err = ioutil.ReadFile(filename); err != nil {
		panic(err)
	}
	handler.Apply(filename)
	myerr("DeleteFromBytes", handler.DeleteFromBytes(data))

	// Output:

	//2022/07/04 21:43:08 DeleteByName success.
	//2022/07/04 21:43:08 Delete success.
	//2022/07/04 21:43:08 DeleteFromFile success.
	//2022/07/04 21:43:08 DeleteFromBytes success.
}
```



## How to get k8s deployment by handler.

```go
func Deployment_Get() {
	// New returns a handler used to multiples deployment.
	handler, err := deployment.New(ctx, namespace, kubeconfig)
	if err != nil {
		panic(err)
	}
	defer cleanup(handler)
	handler.Apply(filename)

	deploy1, err := handler.GetByName(name)
	myerr("GetByName", err)

	deploy2, err := handler.Get(name)
	myerr("Get", err)

	deploy3, err := handler.GetFromFile(filename)
	myerr("GetFromFile", err)

	var data []byte
	if data, err = ioutil.ReadFile(filename); err != nil {
		panic(err)
	}
	deploy4, err := handler.GetFromBytes(data)
	myerr("GetFromBytes", err)

	log.Println(deploy1.Name, deploy2.Name, deploy3.Name, deploy4.Name)

	// Output:

	//2022/07/04 21:51:05 GetByName success.
	//2022/07/04 21:51:05 Get success.
	//2022/07/04 21:51:05 GetFromFile success.
	//2022/07/04 21:51:05 GetFromBytes success.
	//2022/07/04 21:51:05 mydep mydep mydep mydep
}
```



## How to list k8s deployment by handler.

```go
func Deployment_List() {
	// New returns a handler used to multiples deployment.
	handler, err := deployment.New(ctx, namespace, kubeconfig)
	if err != nil {
		panic(err)
	}
	defer cleanup(handler)

	handler.Apply(filename2)

	// ListByLabel list deployment by label.
	deployList1, err := handler.ListByLabel(label)
	myerr("ListByLabel", err)
	outputDeploy(*deployList1)

	// List list deployment by label, it's alias to "ListByLabel".
	deployList2, err := handler.List(label)
	myerr("List", err)
	outputDeploy(*deployList2)

	// ListByNamespace list all deployments in the namespace where the deployment is running.
	deployList3, err := handler.ListByNamespace(namespace)
	myerr("ListByNamespace", err)
	outputDeploy(*deployList3)

	// ListAll list all deployments in the k8s cluster.
	deployList4, err := handler.ListAll()
	myerr("ListAll", err)
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



## How to watch k8s deployment by handler.

```go
func Deployment_Watch() {
	// New returns a handler used to multiples deployment.
	handler, err := deployment.New(ctx, namespace, kubeconfig)
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



## How to get a deployment more information.

```go
func Deployment_Tools() {
	handler, err := deployment.New(ctx, namespace, kubeconfig)
	if err != nil {
		panic(err)
	}
	defer cleanup(handler)

	k8s.ApplyF(ctx, kubeconfig, nginxFile)
	name = "nginx-deploy"

	log.Println(handler.IsReady(name)) // false
	handler.WaitReady(name)            // block until the deployment is ready and available.
	log.Println(handler.IsReady(name)) // true

	// GetPods get all pods that generated by the deployments.
	podList, err := handler.GetPods(name)
	outputPods(podList)

	// GetPVC get all persistentvolumeclaim that attached by the deployment.
	pvcList, err := handler.GetPVC(name)
	log.Println(pvcList)

	// GetPV get all persistentvolume that attached by the deployment.
	pvList, err := handler.GetPV(name)
	log.Println(pvList)

	// cleanup
	k8s.DeleteF(ctx, kubeconfig, nginxFile)

	// Output:

	//2022/07/05 08:33:03 false
	//2022/07/05 08:33:45 true
	//2022/07/05 08:33:45 [nginx-deploy-79979d95dd-v4bkt]
	//2022/07/05 08:33:45 [deploy-k8s-tools-data deploy-nginx-data deploy-nginx-html]
	//2022/07/05 08:33:45 [pvc-fd08839d-eadd-47ea-a21e-380e6c4b7227 pvc-749ab51d-210d-4b1e-bafa-5b8ac758e2e3 pvc-db135fc5-543e-4d20-a7ca-e9afac5882dc]
}

func outputPods(podList []corev1.Pod) {
	var pl []string

	for _, pod := range podList {
		pl = append(pl, pod.Name)
	}
	log.Println(pl)
}
```

