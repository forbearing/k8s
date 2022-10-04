package main

import (
	"fmt"

	"github.com/forbearing/k8s/pod"
)

func Pod_Others() {
	fmt.Println("Group: ", pod.Group)
	fmt.Println("Version: ", pod.Version)
	fmt.Println("Resource: ", pod.Resource)
	// Output:
	//Group:
	//Version:  v1
	//Resource:  pods
}
