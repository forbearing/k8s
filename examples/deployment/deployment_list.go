package main

import (
	"log"

	"github.com/forbearing/k8s"
	"github.com/forbearing/k8s/deployment"
	appsv1 "k8s.io/api/apps/v1"
)

func Deployment_List() {
	// New returns a handler used to multiples deployment.
	handler, err := deployment.New(ctx, namespace, kubeconfig)
	if err != nil {
		panic(err)
	}
	defer cleanup(handler)

	k8s.ApplyF(ctx, kubeconfig, filename2)

	// ListByLabel list deployment by label.
	deployList1, err := handler.ListByLabel(label)
	checkErr("ListByLabel", "", err)
	outputDeploy(*deployList1)

	// List list deployment by label, it's alias to "ListByLabel".
	deployList2, err := handler.List(label)
	checkErr("List", "", err)
	outputDeploy(*deployList2)

	// ListByNamespace list all deployments in the namespace where the deployment is running.
	deployList3, err := handler.ListByNamespace(namespace)
	checkErr("ListByNamespace", "", err)
	outputDeploy(*deployList3)

	// ListAll list all deployments in the k8s cluster.
	deployList4, err := handler.ListAll()
	checkErr("ListAll", "", err)
	outputDeploy(*deployList4)

	// Output:

	//2022/07/04 21:43:09 ListByLabel success.
	//2022/07/04 21:43:09 [mydep-2 nginx-deploy]
	//2022/07/04 21:43:09 List success.
	//2022/07/04 21:43:09 [mydep-2 nginx-deploy]
	//2022/07/04 21:43:09 ListByNamespace success.
	//2022/07/04 21:43:09 [mydep-2 nginx-deploy]
	//2022/07/04 21:43:09 ListAll success.
	//2022/07/04 21:43:09 [calico-kube-controllers coredns metrics-server local-path-provisioner nfs-provisioner-nfs-subdir-external-provisioner mydep-2 nginx-deploy]

}

func outputDeploy(deployList appsv1.DeploymentList) {
	var dl []string
	for _, deploy := range deployList.Items {
		dl = append(dl, deploy.Name)
	}
	log.Println(dl)
}
