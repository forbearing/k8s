package main

import (
	"fmt"

	"github.com/forbearing/k8s/pod"
)

func Pod_Others() {
	fmt.Println("Group: ", pod.GVR().Group)
	fmt.Println("Version: ", pod.GVR().Version)
	fmt.Println("Resource: ", pod.GVR().Resource)
	// Output:
	//Group:
	//Version:  v1
	//Resource:  pods
}
