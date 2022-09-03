package main

import (
	"context"

	"github.com/forbearing/k8s/dynamic"
	"k8s.io/client-go/tools/clientcmd"
)

func Dynamic_Update() {
	handler := dynamic.NewOrDie(context.TODO(), clientcmd.RecommendedHomeFile)
	defer cleanup(handler)

	// update deployment.
	if _, err := handler.WithNamespace("test").Create(deployUnstructData); err != nil {
		panic(err)
	}
	_, err := handler.WithNamespace("test").Update(deployUnstructData)
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

	//2022/08/10 13:57:40 update deployment success:
	//2022/08/10 13:57:41 update namespace success:
	//2022/08/10 13:57:41 update persistentvolume success:
	//2022/08/10 13:57:41 update clusterrole success:
}
