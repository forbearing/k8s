package main

import (
	"fmt"
	"log"

	"github.com/forbearing/k8s"
)

func ServerGroups() {
	discoveryClient := k8s.DiscoveryClientOrDie("")

	groups, err := discoveryClient.ServerGroups()
	if err != nil {
		log.Fatal(err)
	}

	for _, group := range groups.Groups {
		fmt.Printf("Group Name:        %v\n", group.Name)
		fmt.Printf("Preferred Version: %v\n", group.PreferredVersion)
		fmt.Printf("Available Version: %v\n", group.Versions)
	}

	// Output

	//Group Name:
	//Preferred Version: {v1 v1}
	//Available Version: [{v1 v1}]
	//Group Name:        apiregistration.k8s.io
	//Preferred Version: {apiregistration.k8s.io/v1 v1}
	//Available Version: [{apiregistration.k8s.io/v1 v1} {apiregistration.k8s.io/v1beta1 v1beta1}]
	//Group Name:        extensions
	//Preferred Version: {extensions/v1beta1 v1beta1}
	//Available Version: [{extensions/v1beta1 v1beta1}]
	//Group Name:        apps
	//Preferred Version: {apps/v1 v1}
	//Available Version: [{apps/v1 v1}]
	//Group Name:        events.k8s.io
	//Preferred Version: {events.k8s.io/v1 v1}
	//Available Version: [{events.k8s.io/v1 v1} {events.k8s.io/v1beta1 v1beta1}]
	//Group Name:        authentication.k8s.io
	//Preferred Version: {authentication.k8s.io/v1 v1}
	//Available Version: [{authentication.k8s.io/v1 v1} {authentication.k8s.io/v1beta1 v1beta1}]
	//Group Name:        authorization.k8s.io
	//Preferred Version: {authorization.k8s.io/v1 v1}
	//Available Version: [{authorization.k8s.io/v1 v1} {authorization.k8s.io/v1beta1 v1beta1}]
	//Group Name:        autoscaling
	//Preferred Version: {autoscaling/v1 v1}
	//Available Version: [{autoscaling/v1 v1} {autoscaling/v2beta1 v2beta1} {autoscaling/v2beta2 v2beta2}]
	//Group Name:        batch
	//Preferred Version: {batch/v1 v1}
	//Available Version: [{batch/v1 v1} {batch/v1beta1 v1beta1}]
	//Group Name:        certificates.k8s.io
	//Preferred Version: {certificates.k8s.io/v1 v1}
	//Available Version: [{certificates.k8s.io/v1 v1} {certificates.k8s.io/v1beta1 v1beta1}]
	//Group Name:        networking.k8s.io
	//Preferred Version: {networking.k8s.io/v1 v1}
	//Available Version: [{networking.k8s.io/v1 v1} {networking.k8s.io/v1beta1 v1beta1}]
	//Group Name:        policy
	//Preferred Version: {policy/v1beta1 v1beta1}
	//Available Version: [{policy/v1beta1 v1beta1}]
	//Group Name:        rbac.authorization.k8s.io
	//Preferred Version: {rbac.authorization.k8s.io/v1 v1}
	//Available Version: [{rbac.authorization.k8s.io/v1 v1} {rbac.authorization.k8s.io/v1beta1 v1beta1}]
	//Group Name:        storage.k8s.io
	//Preferred Version: {storage.k8s.io/v1 v1}
	//Available Version: [{storage.k8s.io/v1 v1} {storage.k8s.io/v1beta1 v1beta1}]
	//Group Name:        admissionregistration.k8s.io
	//Preferred Version: {admissionregistration.k8s.io/v1 v1}
	//Available Version: [{admissionregistration.k8s.io/v1 v1} {admissionregistration.k8s.io/v1beta1 v1beta1}]
	//Group Name:        apiextensions.k8s.io
	//Preferred Version: {apiextensions.k8s.io/v1 v1}
	//Available Version: [{apiextensions.k8s.io/v1 v1} {apiextensions.k8s.io/v1beta1 v1beta1}]
	//Group Name:        scheduling.k8s.io
	//Preferred Version: {scheduling.k8s.io/v1 v1}
	//Available Version: [{scheduling.k8s.io/v1 v1} {scheduling.k8s.io/v1beta1 v1beta1}]
	//Group Name:        coordination.k8s.io
	//Preferred Version: {coordination.k8s.io/v1 v1}
	//Available Version: [{coordination.k8s.io/v1 v1} {coordination.k8s.io/v1beta1 v1beta1}]
	//Group Name:        node.k8s.io
	//Preferred Version: {node.k8s.io/v1beta1 v1beta1}
	//Available Version: [{node.k8s.io/v1beta1 v1beta1}]
	//Group Name:        discovery.k8s.io
	//Preferred Version: {discovery.k8s.io/v1beta1 v1beta1}
	//Available Version: [{discovery.k8s.io/v1beta1 v1beta1}]
}
