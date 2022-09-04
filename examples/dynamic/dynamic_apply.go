package main

import (
	"context"

	"github.com/forbearing/k8s/dynamic"
)

func Dynamic_Apply() {
	handler := dynamic.NewOrDie(context.TODO(), "", namespace)
	defer cleanup(handler)

	// apply deployment
	_, err := handler.Apply(deployUnstructData)
	checkErr("apply deployment from map[string]interface{}", "", err)
	_, err = handler.Apply("../../testdata/examples/deployment.yaml")
	checkErr("apply deployment from yaml file", "", err)
	_, err = handler.Apply(("../../testdata/examples/deployment.json"))
	checkErr("apply deployment from json file", "", err)

	// apply pod
	_, err = handler.Apply(podUnstructData)
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

	//2022/09/03 22:03:30 apply deployment success:
	//2022/09/03 22:03:30 apply pod success:
	//2022/09/03 22:03:30 apply namespace success:
	//2022/09/03 22:03:30 apply persistentvolume success:
	//2022/09/03 22:03:30 apply clusterrole success:
}
