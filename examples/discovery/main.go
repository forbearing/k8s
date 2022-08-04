package main

import (
	"context"
	"fmt"

	"github.com/forbearing/k8s/pod"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	handler, err := pod.New(context.TODO(), clientcmd.RecommendedHomeFile, "")
	if err != nil {
		panic(err)
	}

	fmt.Println(handler.DiscoveryClient().ServerGroups())
	fmt.Println()
	fmt.Println(handler.DiscoveryClient().ServerGroupsAndResources())
	fmt.Println()
	fmt.Println(handler.DiscoveryClient().ServerPreferredNamespacedResources())
	fmt.Println()
	fmt.Println(handler.DiscoveryClient().ServerPreferredResources())
	fmt.Println()
	fmt.Println(handler.DiscoveryClient().ServerResourcesForGroupVersion("apps/v1"))
}
