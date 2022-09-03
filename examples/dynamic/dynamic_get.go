package main

import (
	"context"

	"github.com/forbearing/k8s/dynamic"
	"github.com/forbearing/k8s/pod"
)

func Dynamic_Get() {
	namespace := "test"
	handler := dynamic.NewOrDie(context.TODO(), "", namespace)
	defer cleanup(handler)

	// get deployment
	if _, err := handler.Apply(deployUnstructData); err != nil {
		panic(err)
	}
	// if the Get() parameter is []byte which containing the resource defination.
	// it's not necessarily to provides GroupVersionKind with WithGVK() method.
	u1, err := handler.Get(deployUnstructData)
	checkErr("get deployment", u1.GetName(), err)

	// get pod
	if _, err := handler.Apply(podUnstructData); err != nil {
		panic(err)
	}
	// if the Get() parameter is resource name, you should call WithGVK()
	// to provides the GroupVersionkind of this resource.
	u2, err := handler.WithGVK(pod.GVK()).Get(podUnstructName)
	checkErr("get pod", u2.GetName(), err)

	// get namespace
	if _, err := handler.Apply(nsUnstructData); err != nil {
		panic(err)
	}
	u3, err := handler.Get(nsUnstructData)
	checkErr("get namespace", u3.GetName(), err)

	// get persistentvolume
	if _, err := handler.Apply(pvUnstructData); err != nil {
		panic(err)
	}
	u4, err := handler.Get(pvUnstructData)
	checkErr("get persistentvolume", u4.GetName(), err)

	// get clusterrole
	if _, err := handler.Apply(crUnstructData); err != nil {
		panic(err)
	}
	u5, err := handler.Get(crUnstructData)
	checkErr("get clusterrole", u5.GetName(), err)

	// Output:

	//2022/09/03 22:04:10 get deployment success: mydep-unstruct
	//2022/09/03 22:04:10 get pod success: pod-unstruct
	//2022/09/03 22:04:10 get namespace success: ns-unstruct
	//2022/09/03 22:04:11 get persistentvolume success: pv-unstruct
	//2022/09/03 22:04:11 get clusterrole success: cr-unstruct
}
