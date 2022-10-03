package main

import (
	"context"

	"github.com/forbearing/k8s/dynamic"
)

func Dynamic_Apply() {
	handler := dynamic.NewOrDie(context.TODO(), "", namespace)
	defer cleanup(handler)

	// apply deployment
	deployObj, err := handler.Apply(deployUnstructData)
	checkErr("apply deployment from map[string]interface{}", deployObj.GetName(), err)
	deployObj2, err := handler.Apply("../../testdata/examples/deployment.yaml")
	checkErr("apply deployment from yaml file", deployObj2.GetName(), err)
	deployObj3, err := handler.Apply(("../../testdata/examples/deployment.json"))
	checkErr("apply deployment from json file", deployObj3.GetName(), err)

	// apply pod
	podObj, err := handler.Apply(podUnstructData)
	checkErr("apply pod", podObj.GetName(), err)

	// apply namespace
	nsObj, err := handler.Apply(nsUnstructData)
	checkErr("apply namespace", nsObj.GetName(), err)

	// apply persistentvolume
	pvObj, err := handler.Apply(pvUnstructData)
	checkErr("apply persistentvolume", pvObj.GetName(), err)

	// apply clusterrole
	crObj, err := handler.Apply(crUnstructData)
	checkErr("apply clusterrole", crObj.GetName(), err)

	// Output:

	//2022/10/04 00:15:26 apply deployment from map[string]interface{} success: mydep-unstruct
	//2022/10/04 00:15:26 apply deployment from yaml file success: mydep
	//2022/10/04 00:15:26 apply deployment from json file success: mydep-json
	//2022/10/04 00:15:26 apply pod success: pod-unstruct
	//2022/10/04 00:15:26 apply namespace success: ns-unstruct
	//2022/10/04 00:15:26 apply persistentvolume success: pv-unstruct
	//2022/10/04 00:15:26 apply clusterrole success: cr-unstruct
}
