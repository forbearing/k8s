package main

import (
	"fmt"

	"github.com/forbearing/k8s/deployment"
)

func Deployment_Others() {
	fmt.Println("Group: ", deployment.Group)
	fmt.Println("Version: ", deployment.Version)
	fmt.Println("Resource: ", deployment.Resource)

	handler, err := deployment.New(ctx, kubeconfig, namespace)
	if err != nil {
		panic(err)
	}
	if _, err := handler.Apply(filename); err != nil {
		panic(err)
	}

	deploy, err := handler.Get(name)
	if err != nil {
		panic(err)
	}

	fmt.Println(deploy.GetObjectKind())

	// Output:
	//Group:  apps
	//Version:  v1
	//Resource:  deployments
}
