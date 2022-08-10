package main

import (
	"context"

	"github.com/forbearing/k8s/clusterrole"
	"github.com/forbearing/k8s/deployment"
	"github.com/forbearing/k8s/dynamic"
	"github.com/forbearing/k8s/namespace"
	"github.com/forbearing/k8s/persistentvolume"
	"github.com/forbearing/k8s/pod"
	"k8s.io/client-go/tools/clientcmd"
)

func Dynamic_Get() {
	handler, err := dynamic.New(context.TODO(), clientcmd.RecommendedHomeFile, "", deployment.GVR())
	if err != nil {
		panic(err)
	}
	defer cleanup(handler)

	// get deployment
	if _, err := handler.WithNamespace("test").Apply(deployUnstructData); err != nil {
		panic(err)
	}
	u1, err := handler.WithNamespace("test").Get(deployUnstructData)
	checkErr("get deployment", u1.GetName(), err)

	// get pod

	if _, err := handler.WithNamespace("test").WithGVR(pod.GVR()).Apply(podUnstructData); err != nil {
		panic(err)
	}
	u2, err := handler.WithNamespace("test").WithGVR(pod.GVR()).Get(podUnstructName)
	checkErr("get pod", u2.GetName(), err)

	// get namespace
	if _, err := handler.WithGVR(namespace.GVR()).Apply(nsUnstructData); err != nil {
		panic(err)
	}
	u3, err := handler.WithGVR(namespace.GVR()).Get(nsUnstructData)
	checkErr("get namespace", u3.GetName(), err)

	// get persistentvolume
	if _, err := handler.WithGVR(persistentvolume.GVR()).Apply(pvUnstructData); err != nil {
		panic(err)
	}
	u4, err := handler.WithGVR(persistentvolume.GVR()).Get(pvUnstructData)
	checkErr("get persistentvolume", u4.GetName(), err)

	// get clusterrole
	if _, err := handler.WithGVR(clusterrole.GVR()).Apply(crUnstructData); err != nil {
		panic(err)
	}
	u5, err := handler.WithGVR(clusterrole.GVR()).Get(crUnstructData)
	checkErr("get clusterrole", u5.GetName(), err)

	// Output:

	//2022/08/10 13:57:04 get deployment success: mydep-unstruct
	//2022/08/10 13:57:04 get pod success: pod-unstruct
	//2022/08/10 13:57:05 get namespace success: ns-unstruct
	//2022/08/10 13:57:05 get persistentvolume success: pv-unstruct
	//2022/08/10 13:57:05 get clusterrole success: cr-unstruct
}
