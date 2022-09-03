package main

import (
	"context"

	"github.com/forbearing/k8s/dynamic"
)

func Dynamic_Update() {
	namespace := "test"
	handler := dynamic.NewOrDie(context.TODO(), "", namespace)
	defer cleanup(handler)

	// update deployment.
	if _, err := handler.Create(deployUnstructData); err != nil {
		panic(err)
	}
	_, err := handler.Update(deployUnstructData)
	checkErr("update deployment", "", err)

	// update namespace
	if _, err := handler.Create(nsUnstructData); err != nil {
		panic(err)
	}
	_, err = handler.Update(nsUnstructData)
	checkErr("update namespace", "", err)

	// update persistentvolume
	if _, err := handler.Create(pvUnstructData); err != nil {
		panic(err)
	}
	_, err = handler.Update(pvUnstructData)
	checkErr("update persistentvolume", "", err)

	// update clusterrole
	if _, err := handler.Create(crUnstructData); err != nil {
		panic(err)
	}
	_, err = handler.Update(crUnstructData)
	checkErr("update clusterrole", "", err)

	// Output:

	//2022/09/03 22:00:38 update deployment success:
	//2022/09/03 22:00:38 update namespace success:
	//2022/09/03 22:00:38 update persistentvolume success:
	//2022/09/03 22:00:38 update clusterrole success:

}
