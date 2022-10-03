package main

import (
	"context"

	"github.com/forbearing/k8s/dynamic"
)

func Dynamic_Update() {
	handler := dynamic.NewOrDie(context.TODO(), "", namespace)
	defer cleanup(handler)

	// update deployment.
	if _, err := handler.Create(deployUnstructData); err != nil {
		panic(err)
	}
	deployObj, err := handler.Update(deployUnstructData)
	checkErr("update deployment", deployObj.GetName(), err)

	// update namespace
	if _, err := handler.Create(nsUnstructData); err != nil {
		panic(err)
	}
	nsObj, err := handler.Update(nsUnstructData)
	checkErr("update namespace", nsObj.GetName(), err)

	// update persistentvolume
	if _, err := handler.Create(pvUnstructData); err != nil {
		panic(err)
	}
	pvObj, err := handler.Update(pvUnstructData)
	checkErr("update persistentvolume", pvObj.GetName(), err)

	// update clusterrole
	if _, err := handler.Create(crUnstructData); err != nil {
		panic(err)
	}
	crObj, err := handler.Update(crUnstructData)
	checkErr("update clusterrole", crObj.GetName(), err)

	// Output:

	//2022/10/04 00:12:10 update deployment success: mydep-unstruct
	//2022/10/04 00:12:10 update namespace success: ns-unstruct
	//2022/10/04 00:12:10 update persistentvolume success: pv-unstruct
	//2022/10/04 00:12:10 update clusterrole success: cr-unstruct
}
