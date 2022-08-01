package main

import (
	"fmt"

	"github.com/forbearing/k8s/deployment"
)

func Deployment_Others() {
	fmt.Println("Group: ", deployment.GVR().Group)
	fmt.Println("Version: ", deployment.GVR().Version)
	fmt.Println("Resource: ", deployment.GVR().Resource)
	// Output:
	//Group:  apps
	//Version:  v1
	//Resource:  deployments
}
