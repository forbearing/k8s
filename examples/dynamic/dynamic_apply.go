package main

import (
	"context"

	"github.com/forbearing/k8s/clusterrole"
	"github.com/forbearing/k8s/deployment"
	"github.com/forbearing/k8s/dynamic"
	"github.com/forbearing/k8s/namespace"
	"github.com/forbearing/k8s/persistentvolume"
	"github.com/forbearing/k8s/pod"
	"k8s.io/client-go/tools/clientcmd"
)

func Dynamic_Apply() {
	handler, err := dynamic.New(context.TODO(), clientcmd.RecommendedHomeFile, "", deployment.GVR())
	if err != nil {
		panic(err)
	}
	defer cleanup(handler)

	// apply deployment
	_, err = handler.WithNamespace("test").Apply(deployUnstructData)
	checkErr("apply deployment", "", err)

	// apply pod
	_, err = handler.WithNamespace("test").WithGVR(pod.GVR()).Apply(podUnstructData)
	checkErr("apply pod", "", err)

	// apply namespace
	_, err = handler.WithGVR(namespace.GVR()).Apply(nsUnstructData)
	checkErr("apply namespace", "", err)

	// apply persistentvolume
	_, err = handler.WithGVR(persistentvolume.GVR()).Apply(pvUnstructData)
	checkErr("apply persistentvolume", "", err)

	// apply clusterrole
	_, err = handler.WithGVR(clusterrole.GVR()).Apply(crUnstructData)
	checkErr("apply clusterrole", "", err)

	// Output:

	//2022/08/10 13:55:00 apply deployment success:
	//2022/08/10 13:55:00 apply pod success:
	//2022/08/10 13:55:00 apply namespace success:
	//2022/08/10 13:55:00 apply persistentvolume success:
	//2022/08/10 13:55:00 apply clusterrole success:
}
