package main

import (
	"context"

	"github.com/forbearing/k8s/dynamic"
)

func Dynamic_Apply() {
	handler := dynamic.NewOrDie(context.TODO(), "")
	defer cleanup(handler)

	// apply deployment
	_, err := handler.WithNamespace("test").Apply(deployUnstructData)
	checkErr("apply deployment", "", err)

	// apply pod
	_, err = handler.WithNamespace("test").Apply(podUnstructData)
	checkErr("apply pod", "", err)

	// apply namespace
	_, err = handler.Apply(nsUnstructData)
	checkErr("apply namespace", "", err)

	// apply persistentvolume
	_, err = handler.Apply(pvUnstructData)
	checkErr("apply persistentvolume", "", err)

	// apply clusterrole
	_, err = handler.Apply(crUnstructData)
	checkErr("apply clusterrole", "", err)

	// Output:

	//2022/08/10 13:55:00 apply deployment success:
	//2022/08/10 13:55:00 apply pod success:
	//2022/08/10 13:55:00 apply namespace success:
	//2022/08/10 13:55:00 apply persistentvolume success:
	//2022/08/10 13:55:00 apply clusterrole success:
}
