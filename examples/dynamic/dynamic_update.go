package main

import (
	"context"

	"github.com/forbearing/k8s/dynamic"
	"k8s.io/client-go/tools/clientcmd"
)

func Dynamic_Update() {
	handler, err := dynamic.New(context.TODO(), clientcmd.RecommendedHomeFile, "", "apps", "v1", "deployments")
	if err != nil {
		panic(err)
	}
	defer cleanup(handler)

	if _, err := handler.Namespace("test").Create(deployUnstructData); err != nil {
		panic(err)
	}
	_, err = handler.Namespace("test").Update(deployUnstructData)
	checkErr("update deployment", "", err)

	if _, err := handler.Group("").Resource("namespaces").Create(nsUnstructData); err != nil {
		panic(err)
	}
	_, err = handler.Group("").Resource("namespaces").Update(nsUnstructData)
	checkErr("update namespace", "", err)

	if _, err := handler.Group("").Resource("persistentvolumes").Create(pvUnstructData); err != nil {
		panic(err)
	}
	_, err = handler.Group("").Resource("persistentvolumes").Update(pvUnstructData)
	checkErr("update persistentvolume", "", err)

	if _, err := handler.Group("rbac.authorization.k8s.io").Resource("clusterroles").Create(crUnstructData); err != nil {
		panic(err)
	}
	_, err = handler.Group("rbac.authorization.k8s.io").Resource("clusterroles").Update(crUnstructData)
	checkErr("update clusterrole", "", err)

	// Output:

	//2022/07/29 17:10:40 update deployment success:
	//2022/07/29 17:10:40 update namespace success:
	//2022/07/29 17:10:40 update persistentvolume success:
	//2022/07/29 17:10:40 update clusterrole success:
}
