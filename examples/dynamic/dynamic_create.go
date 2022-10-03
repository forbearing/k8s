package main

import (
	"context"

	"github.com/forbearing/k8s/dynamic"
)

func Dynamic_Create() {
	handler := dynamic.NewOrDie(context.TODO(), "", namespace)
	defer cleanup(handler)

	// create deployment
	deployObj, err := handler.Create(deployUnstructData)
	checkErr("create deployment", deployObj.GetName(), err)

	// create pod
	podObj, err := handler.Create(podUnstructData)
	checkErr("create pod", podObj.GetName(), err)

	// create namespace
	nsObj, err := handler.Create(nsUnstructData)
	checkErr("create namespace", nsObj.GetName(), err)

	// create persistentvolume
	pvObj, err := handler.Create(pvUnstructData)
	checkErr("create persistentvolume", pvObj.GetName(), err)

	// create clusterrole
	crObj, err := handler.Create(crUnstructData)
	checkErr("create clusterrole", crObj.GetName(), err)

	// Output:

	//2022/10/04 00:10:15 create deployment success: mydep-unstruct
	//2022/10/04 00:10:15 create pod success: pod-unstruct
	//2022/10/04 00:10:15 create namespace success: ns-unstruct
	//2022/10/04 00:10:15 create persistentvolume success: pv-unstruct
	//2022/10/04 00:10:15 create clusterrole success: cr-unstruct
}
