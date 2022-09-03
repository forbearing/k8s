package main

import (
	"context"

	"github.com/forbearing/k8s/dynamic"
)

func Dynamic_Create() {
	handler := dynamic.NewOrDie(context.TODO(), "")
	defer cleanup(handler)

	// create deployment
	_, err := handler.WithNamespace("test").Create(deployUnstructData)
	checkErr("create deployment", "", err)

	// create pod
	_, err = handler.WithNamespace("test").Create(podUnstructData)
	checkErr("create pod", "", err)

	// create namespace
	_, err = handler.Create(nsUnstructData)
	checkErr("create namespace", "", err)

	// create persistentvolume
	_, err = handler.Create(pvUnstructData)
	checkErr("create persistentvolume", "", err)

	// create clusterrole
	_, err = handler.Create(crUnstructData)
	checkErr("create clusterrole", "", err)

	// Output:

	//2022/08/10 13:58:21 create deployment success:
	//2022/08/10 13:58:21 create pod success:
	//2022/08/10 13:58:21 create namespace success:
	//2022/08/10 13:58:21 create persistentvolume success:
	//2022/08/10 13:58:21 create clusterrole success:
}
