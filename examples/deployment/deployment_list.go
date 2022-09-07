package main

import (
	"fmt"

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

	k8s.ApplyF(ctx, kubeconfig, filename2, namespace, k8s.IgnoreInvalid)

	label := "type=deployment"
	field := fmt.Sprintf("metadata.namespace=%s", namespace)

	dl, err := handler.List()
	checkErr("List()", outputDeploy(dl), err)
	dl2, err := handler.ListAll()
	checkErr("ListAll()", outputDeploy(dl2), err)
	dl3, err := handler.ListByLabel(label)
	checkErr("ListByLabel()", outputDeploy(dl3), err)
	dl4, err := handler.ListByField(field)
	checkErr("ListByField()", outputDeploy(dl4), err)
	dl5, err := handler.ListByNamespace("kube-system")
	checkErr("ListByNamespace()", outputDeploy(dl5), err)

	// Output:

	//2022/09/07 18:29:35 List() success: [nginx coredns local-path-provisioner nginx-deploy]
	//2022/09/07 18:29:35 ListAll() success: [nginx coredns local-path-provisioner nginx-deploy]
	//2022/09/07 18:29:35 ListByLabel() success: [nginx-deploy]
	//2022/09/07 18:29:35 ListByField() success: [nginx-deploy]
	//2022/09/07 18:29:35 ListByNamespace() success: [coredns]
}

func outputDeploy(deployList []*appsv1.Deployment) []string {
	if deployList == nil {
		return nil
	}
	var dl []string
	for _, deploy := range deployList {
		dl = append(dl, deploy.Name)
	}
	return dl
}
