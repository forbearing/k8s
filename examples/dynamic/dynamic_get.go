package main

import (
	"context"

	"github.com/forbearing/k8s/dynamic"
	"github.com/forbearing/k8s/pod"
)

func Dynamic_Get() {
	handler := dynamic.NewOrDie(context.TODO(), "")
	defer cleanup(handler)

	// get deployment
	if _, err := handler.WithNamespace("test").Apply(deployUnstructData); err != nil {
		panic(err)
	}
	u1, err := handler.WithNamespace("test").Get(deployUnstructData)
	checkErr("get deployment", u1.GetName(), err)

	// get pod
	if _, err := handler.WithNamespace("test").Apply(podUnstructData); err != nil {
		panic(err)
	}
	u2, err := handler.WithGVK(pod.GVK()).WithNamespace("test").Get(podUnstructName)
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

	//2022/08/10 13:57:04 get deployment success: mydep-unstruct
	//2022/08/10 13:57:04 get pod success: pod-unstruct
	//2022/08/10 13:57:05 get namespace success: ns-unstruct
	//2022/08/10 13:57:05 get persistentvolume success: pv-unstruct
	//2022/08/10 13:57:05 get clusterrole success: cr-unstruct
}
