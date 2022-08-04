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
	handler, err := dynamic.New(context.TODO(), clientcmd.RecommendedHomeFile, "", "apps", "v1", "deployments")
	if err != nil {
		panic(err)
	}
	defer cleanup(handler)

	//_, err = handler.Namespace("test").Create(deployUnstructData)
	//checkErr("create deployment", "", err)
	//_, err = handler.Namespace("test").Group("").Resource("pods").Create(podUnstructData)
	//checkErr("create pod", "", err)
	//_, err = handler.Group("").Resource("namespaces").Create(nsUnstructData)
	//checkErr("create namespace", "", err)
	//_, err = handler.Group("").Resource("persistentvolumes").Create(pvUnstructData)
	//checkErr("create persistentvolume", "", err)
	//_, err = handler.Group("rbac.authorization.k8s.io").Resource("clusterroles").Create(crUnstructData)
	//checkErr("create clusterrole", "", err)

	_, err = handler.Namespace("test").GVR(deployment.GVR()).Create(deployUnstructData)
	checkErr("create deployment", "", err)
	_, err = handler.Namespace("test").GVR(pod.GVR()).Create(podUnstructData)
	checkErr("create pod", "", err)
	_, err = handler.GVR(namespace.GVR()).Create(nsUnstructData)
	checkErr("create namespace", "", err)
	_, err = handler.GVR(persistentvolume.GVR()).Create(pvUnstructData)
	checkErr("create persistentvolume", "", err)
	_, err = handler.GVR(clusterrole.GVR()).Create(crUnstructData)
	checkErr("create clusterrole", "", err)

	// Output:

	//2022/07/29 17:05:05 create deployment success:
	//2022/07/29 17:05:05 create pod success:
	//2022/07/29 17:05:05 create namespace success:
	//2022/07/29 17:05:05 create persistentvolume success:
	//2022/07/29 17:05:05 create clusterrole success:
}
