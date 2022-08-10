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

func Dynamic_Create() {
	handler, err := dynamic.New(context.TODO(), clientcmd.RecommendedHomeFile, "", deployment.GVR())
	if err != nil {
		panic(err)
	}
	defer cleanup(handler)

	// create deployment
	_, err = handler.WithNamespace("test").Create(deployUnstructData)
	checkErr("create deployment", "", err)

	// create pod
	_, err = handler.WithNamespace("test").WithGVR(pod.GVR()).Create(podUnstructData)
	checkErr("create pod", "", err)

	// create namespace
	_, err = handler.WithGVR(namespace.GVR()).Create(nsUnstructData)
	checkErr("create namespace", "", err)

	// create persistentvolume
	_, err = handler.WithGVR(persistentvolume.GVR()).Create(pvUnstructData)
	checkErr("create persistentvolume", "", err)

	// create clusterrole
	_, err = handler.WithGVR(clusterrole.GVR()).Create(crUnstructData)
	checkErr("create clusterrole", "", err)

	// Output:

	//2022/08/10 13:58:21 create deployment success:
	//2022/08/10 13:58:21 create pod success:
	//2022/08/10 13:58:21 create namespace success:
	//2022/08/10 13:58:21 create persistentvolume success:
	//2022/08/10 13:58:21 create clusterrole success:
}
