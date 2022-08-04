package main

import (
	"context"

	"github.com/forbearing/k8s/dynamic"
	"k8s.io/client-go/tools/clientcmd"
)

func Dynamic_Apply() {
	handler, err := dynamic.New(context.TODO(), clientcmd.RecommendedHomeFile, "", "apps", "v1", "deployments")
	if err != nil {
		panic(err)
	}
	defer cleanup(handler)

	_, err = handler.Namespace("test").Apply(deployUnstructData)
	checkErr("apply deployment", "", err)
	_, err = handler.Namespace("test").Group("").Resource("pods").Apply(podUnstructData)
	checkErr("apply pod", "", err)
	_, err = handler.Group("").Resource("namespaces").Apply(nsUnstructData)
	checkErr("apply namespace", "", err)
	_, err = handler.Group("").Resource("persistentvolumes").Apply(pvUnstructData)
	checkErr("apply persistentvolume", "", err)
	_, err = handler.Group("rbac.authorization.k8s.io").Resource("clusterroles").Apply(crUnstructData)
	checkErr("apply clusterrole", "", err)

	// Output:

	//2022/07/29 17:04:11 apply deployment success:
	//2022/07/29 17:04:11 apply pod success:
	//2022/07/29 17:04:11 apply namespace success:
	//2022/07/29 17:04:11 apply persistentvolume success:
	//2022/07/29 17:04:11 apply clusterrole success:
}
