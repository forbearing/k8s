package main

import (
	"context"

	"github.com/forbearing/k8s/dynamic"
	"k8s.io/client-go/tools/clientcmd"
)

func Dynamic_Get() {
	handler, err := dynamic.New(context.TODO(), clientcmd.RecommendedHomeFile, "", "apps", "v1", "deployments")
	if err != nil {
		panic(err)
	}
	defer cleanup(handler)

	if _, err := handler.Namespace("test").Apply(deployUnstructData); err != nil {
		panic(err)
	}
	u1, err := handler.Namespace("test").Get(deployUnstructData)
	checkErr("get deployment", u1.GetName(), err)

	if _, err := handler.Group("").Resource("namespaces").Apply(nsUnstructData); err != nil {
		panic(err)
	}
	if _, err := handler.Namespace("test").Group("").Resource("pods").Apply(podUnstructData); err != nil {
		panic(err)
	}
	u2, err := handler.Namespace("test").Group("").Resource("pods").Get(podUnstructName)
	checkErr("get pod", u2.GetName(), err)
	u3, err := handler.Group("").Resource("namespaces").Get(nsUnstructData)
	checkErr("get namespace", u3.GetName(), err)

	if _, err := handler.Group("").Resource("persistentvolumes").Apply(pvUnstructData); err != nil {
		panic(err)
	}
	u4, err := handler.Group("").Resource("persistentvolumes").Get(pvUnstructData)
	checkErr("get persistentvolume", u4.GetName(), err)

	if _, err := handler.Group("rbac.authorization.k8s.io").Resource("clusterroles").Apply(crUnstructData); err != nil {
		panic(err)
	}
	u5, err := handler.Group("rbac.authorization.k8s.io").Resource("clusterroles").Get(crUnstructData)
	checkErr("get clusterrole", u5.GetName(), err)

	// Output:

	//2022/07/29 17:16:31 get deployment success: mydep-unstruct
	//2022/07/29 17:16:31 get pod success: pod-unstruct
	//2022/07/29 17:16:31 get namespace success: ns-unstruct
	//2022/07/29 17:16:31 get persistentvolume success: pv-unstruct
	//2022/07/29 17:16:31 get clusterrole success: cr-unstruct
}
