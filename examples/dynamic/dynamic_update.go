package main

import (
	"context"

	"github.com/forbearing/k8s/clusterrole"
	"github.com/forbearing/k8s/deployment"
	"github.com/forbearing/k8s/dynamic"
	"github.com/forbearing/k8s/namespace"
	"github.com/forbearing/k8s/persistentvolume"
	"k8s.io/client-go/tools/clientcmd"
)

func Dynamic_Update() {
	handler, err := dynamic.New(context.TODO(), clientcmd.RecommendedHomeFile, "", deployment.GVR())
	if err != nil {
		panic(err)
	}
	defer cleanup(handler)

	// update deployment.
	if _, err := handler.WithNamespace("test").Create(deployUnstructData); err != nil {
		panic(err)
	}
	_, err = handler.WithNamespace("test").Update(deployUnstructData)
	checkErr("update deployment", "", err)

	// update namespace
	if _, err := handler.WithGVR(namespace.GVR()).Create(nsUnstructData); err != nil {
		panic(err)
	}
	_, err = handler.WithGVR(namespace.GVR()).Update(nsUnstructData)
	checkErr("update namespace", "", err)

	// update persistentvolume
	if _, err := handler.WithGVR(persistentvolume.GVR()).Create(pvUnstructData); err != nil {
		panic(err)
	}
	_, err = handler.WithGVR(persistentvolume.GVR()).Update(pvUnstructData)
	checkErr("update persistentvolume", "", err)

	// update clusterrole
	if _, err := handler.WithGVR(clusterrole.GVR()).Create(crUnstructData); err != nil {
		panic(err)
	}
	_, err = handler.WithGVR(clusterrole.GVR()).Update(crUnstructData)
	checkErr("update clusterrole", "", err)

	// Output:

	//2022/08/10 13:57:40 update deployment success:
	//2022/08/10 13:57:41 update namespace success:
	//2022/08/10 13:57:41 update persistentvolume success:
	//2022/08/10 13:57:41 update clusterrole success:
}
