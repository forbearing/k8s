package main

import (
	"context"

	"github.com/forbearing/k8s/dynamic"
)

func Dynamic_Create() {
	namespace := "test"
	handler := dynamic.NewOrDie(context.TODO(), "", namespace)
	defer cleanup(handler)

	// create deployment
	_, err := handler.Create(deployUnstructData)
	checkErr("create deployment", "", err)

	// create pod
	_, err = handler.Create(podUnstructData)
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

	//2022/09/03 21:58:35 create deployment success:
	//2022/09/03 21:58:35 create pod success:
	//2022/09/03 21:58:35 create namespace success:
	//2022/09/03 21:58:35 create persistentvolume success:
	//2022/09/03 21:58:35 create clusterrole success:

}
