package main

import (
	"fmt"

	"github.com/forbearing/k8s"
)

func All() {
	discoveryClient := k8s.DiscoveryClientOrDie("")

	// ServerGroups returns the supported groups, with information like supported versions and the
	// preferred version.
	fmt.Println(discoveryClient.ServerGroups())
	fmt.Println()

	// ServerGroupsAndResources returns the supported resources for all groups and versions.
	fmt.Println(discoveryClient.ServerGroupsAndResources())
	fmt.Println()

	// ServerPreferredNamespacedResources returns the supported namespaced resources with the
	// version preferred by the server.
	fmt.Println(discoveryClient.ServerPreferredNamespacedResources())
	fmt.Println()

	// ServerPreferredResources returns the supported resources with the version preferred by the
	// server.
	fmt.Println(discoveryClient.ServerPreferredResources())
	fmt.Println()

	// ServerResourcesForGroupVersion returns the supported resources for a group and version.
	fmt.Println(discoveryClient.ServerResourcesForGroupVersion("apps/v1"))
}
