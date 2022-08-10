package main

import (
	"github.com/forbearing/k8s"
	"github.com/forbearing/k8s/deployment"
	appsv1 "k8s.io/api/apps/v1"
)

func Deployment_List() {
	// New returns a handler used to multiples deployment.
	handler, err := deployment.New(ctx, kubeconfig, namespace)
	if err != nil {
		panic(err)
	}
	defer cleanup(handler)

	k8s.ApplyF(ctx, kubeconfig, filename2)

	// ListByLabel list deployment by label.
	deployList, err := handler.WithNamespace("kube-system").ListByLabel("k8s-app=kube-dns")
	checkErr("ListByLabel", outputDeploy(deployList), err)

	// List list deployment by label, it simply call `ListByLabel`.
	deployList2, err := handler.WithNamespace("kube-system").ListByLabel("k8s-app=kube-dns")
	checkErr("List", outputDeploy(deployList2), err)

	// ListByNamespace list all deployments in the namespace where the deployment is running.
	deployList3, err := handler.ListByNamespace("kube-system")
	checkErr("ListByNamespace", outputDeploy(deployList3), err)

	// ListAll list all deployments in the k8s cluster.
	deployList4, err := handler.ListAll()
	checkErr("ListAll", outputDeploy(deployList4), err)

	// Output:

	//2022/08/10 14:36:52 ListByLabel success: [coredns]
	//2022/08/10 14:36:52 List success: [coredns]
	//2022/08/10 14:36:52 ListByNamespace success: [calico-kube-controllers coredns metrics-server]
	//2022/08/10 14:36:53 ListAll success: [k8s-tools ......]
}

func outputDeploy(deployList []*appsv1.Deployment) []string {
	var dl []string
	for _, deploy := range deployList {
		dl = append(dl, deploy.Name)
	}
	return dl
}
