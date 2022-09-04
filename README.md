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

`go get github.com/forbearing/k8s@v0.10.4`

## Documents

### pod handler examples

- [How to execute command within pod](./examples/pod/pod_execute.go)
- [How to port-forward a local port to pod](./examples/port-forward/portforward_pod.go)

### deployment handler examples

- [How to create deployment resources inside cluster or outside cluster](./examples/deployment/deployment_create.go)
- [How to update deployment resources inside cluster or outside cluster](./examples/deployment/deployment_update.go)
- [How to apply deployment resources](./examples/deployment/deployment_apply.go)
- [How to delete deployment resources](./examples/deployment/deployment_delete.go)
- [How to get deployment resources](./examples/deployment/deployment_get.go)
- [How to list deployment resources](./examples/deployment/deployment_list.go)
- [How to watch deployment resources](./examples/deployment/deployment_watch.go)
- [deployment informer](./examples/deployment/deployment_informer.go)
- [more usage for informer](./deployment/informer.go)
- [Tools](./examples/deployment/deployment_tools.go)
    - GetPods(): get all pods ownerd by a deployment
    - GetRS(): get all replicaset ownerd by a deployment
    - GetPVC()/GetPV(): Get PVC/PV mounted by a deployment
    - IsReady(): check if a deployment is ready/available/rollout finished.
    - WaitReady(): block here until a deployment is ready/available/rollout finished.

### dynamic handler

### k8s handler

### more

- [ApplyF()/DeleteF() apply/delete various k8s resource from a yaml file.](./k8s_test.go)
- [Check whether the k8s resources has the specifed label, Get/Set/Remove labels of k8s resources](./examples/labels/main.go)
- [Check whether the k8s resources has the specifed annotation, Get/Set/Remove labels of k8s annotations](./examples/annotations/main.go)
- [Find k8s resource's GroupVersionKind from yaml file, json file, bytes data, map[string]interface{}, etc.](./examples/restmapper/find_gvk.go)
- [Find k8s resource's GroupVersionResources from yaml file, json file, bytes data, map[string]interface{}, etc.](./examples/restmapper/find_gvr.go)
- [Check if the k8s resource is namespace scope from yaml file, json file, bytes data, map[string]interface{}, etc.](./examples/restmapper/is_namespaced.go)

## TODO

- [ ] https://github.com/kubernetes/kubectl/tree/master/pkg
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
- [x] port-forward for pod, deployment and service
- [ ] proxy
- [ ] operators refer to https://sdk.operatorframework.io/docs/building-operators/golang/references/client/
- [x] Add/Set/Remove/Has Labels and Annotations
