package main

import (
	"fmt"
	"log"

	"github.com/forbearing/k8s"
	"github.com/forbearing/k8s/clusterrole"
	"github.com/forbearing/k8s/clusterrolebinding"
	"github.com/forbearing/k8s/cronjob"
	"github.com/forbearing/k8s/role"
	"github.com/forbearing/k8s/rolebinding"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/discovery"
)

// ServerGroupsAndResources returns the supported resources for all groups and versions.

func ServerGroupsAndResources() {
	discoveryClient := k8s.DiscoveryClientOrDie("")
	//allResources(discoveryClient)

	_, resourceLists, err := discoveryClient.ServerGroupsAndResources()
	if err != nil {
		log.Fatal(err)
	}

	var gv string
	if gv, err = GroupVersionForResource(resourceLists, cronjob.Resource()); err != nil {
		log.Fatal(err)
	}
	fmt.Println("CronJob Group And Version: ", gv)
	if gv, err = GroupVersionForResource(resourceLists, clusterrole.Resource()); err != nil {
		log.Fatal(err)
	}
	fmt.Println("ClusterRole Group And Version: ", gv)
	if gv, err = GroupVersionForResource(resourceLists, clusterrolebinding.Resource()); err != nil {
		log.Fatal(err)
	}
	fmt.Println("ClusterRoleBinding Group And Version: ", gv)
	if gv, err = GroupVersionForResource(resourceLists, role.Resource()); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Role Group And Version: ", gv)
	if gv, err = GroupVersionForResource(resourceLists, rolebinding.Resource()); err != nil {
		log.Fatal(err)
	}
	fmt.Println("RoleBinding Group And Version: ", gv)

	// Output:

	//CronJob Group And Version:  batch/v1beta1
	//ClusterRole Group And Version:  rbac.authorization.k8s.io/v1
	//ClusterRoleBinding Group And Version:  rbac.authorization.k8s.io/v1
	//Role Group And Version:  rbac.authorization.k8s.io/v1
	//RoleBinding Group And Version:  rbac.authorization.k8s.io/v1
}

func GroupVersionForResource(resourceLists []*metav1.APIResourceList, resourceName string) (string, error) {
	for _, resourceList := range resourceLists {
		for _, resource := range resourceList.APIResources {
			if resource.Name == resourceName {
				return resourceList.GroupVersion, nil
			}
		}
	}
	return "", fmt.Errorf("not found %s GroupVersion", resourceName)
}

// allResources
func allResources(discoveryClient *discovery.DiscoveryClient) {
	// groups 相当于 discoveryClient.ServerGroups 获取到的 *metav1.APIGroupList.Items
	groups, resourceLists, err := discoveryClient.ServerGroupsAndResources()
	if err != nil {
		log.Fatal(err)
	}
	_, _ = groups, resourceLists

	for _, group := range groups {
		fmt.Printf("Group Name:        %v\n", group.Name)
		fmt.Printf("Preferred Version: %v\n", group.PreferredVersion)
		fmt.Printf("Available Version: %v\n", group.Versions)
		fmt.Println()
	}
	for _, resourceList := range resourceLists {
		fmt.Printf("Group Version: %v\n", resourceList.GroupVersion)
		for _, resource := range resourceList.APIResources {
			fmt.Printf("\tResource Name: %v\n", resource.Name)
		}
	}

	// Output:

	//Group Version: v1
	//        Resource Name: bindings
	//        Resource Name: componentstatuses
	//        Resource Name: configmaps
	//        Resource Name: endpoints
	//        Resource Name: events
	//        Resource Name: limitranges
	//        Resource Name: namespaces
	//        Resource Name: namespaces/finalize
	//        Resource Name: namespaces/status
	//        Resource Name: nodes
	//        Resource Name: nodes/proxy
	//        Resource Name: nodes/status
	//        Resource Name: persistentvolumeclaims
	//        Resource Name: persistentvolumeclaims/status
	//        Resource Name: persistentvolumes
	//        Resource Name: persistentvolumes/status
	//        Resource Name: pods
	//        Resource Name: pods/attach
	//        Resource Name: pods/binding
	//        Resource Name: pods/eviction
	//        Resource Name: pods/exec
	//        Resource Name: pods/log
	//        Resource Name: pods/portforward
	//        Resource Name: pods/proxy
	//        Resource Name: pods/status
	//        Resource Name: podtemplates
	//        Resource Name: replicationcontrollers
	//        Resource Name: replicationcontrollers/scale
	//        Resource Name: replicationcontrollers/status
	//        Resource Name: resourcequotas
	//        Resource Name: resourcequotas/status
	//        Resource Name: secrets
	//        Resource Name: serviceaccounts
	//        Resource Name: services
	//        Resource Name: services/proxy
	//        Resource Name: services/status
	//Group Version: apiregistration.k8s.io/v1
	//        Resource Name: apiservices
	//        Resource Name: apiservices/status
	//Group Version: apiregistration.k8s.io/v1beta1
	//        Resource Name: apiservices
	//        Resource Name: apiservices/status
	//Group Version: extensions/v1beta1
	//        Resource Name: ingresses
	//        Resource Name: ingresses/status
	//Group Version: apps/v1
	//        Resource Name: controllerrevisions
	//        Resource Name: daemonsets
	//        Resource Name: daemonsets/status
	//        Resource Name: deployments
	//        Resource Name: deployments/scale
	//        Resource Name: deployments/status
	//        Resource Name: replicasets
	//        Resource Name: replicasets/scale
	//        Resource Name: replicasets/status
	//        Resource Name: statefulsets
	//        Resource Name: statefulsets/scale
	//        Resource Name: statefulsets/status
	//Group Version: events.k8s.io/v1
	//        Resource Name: events
	//Group Version: events.k8s.io/v1beta1
	//        Resource Name: events
	//Group Version: authentication.k8s.io/v1
	//        Resource Name: tokenreviews
	//Group Version: authentication.k8s.io/v1beta1
	//        Resource Name: tokenreviews
	//Group Version: authorization.k8s.io/v1
	//        Resource Name: localsubjectaccessreviews
	//        Resource Name: selfsubjectaccessreviews
	//        Resource Name: selfsubjectrulesreviews
	//        Resource Name: subjectaccessreviews
	//Group Version: authorization.k8s.io/v1beta1
	//        Resource Name: localsubjectaccessreviews
	//        Resource Name: selfsubjectaccessreviews
	//        Resource Name: selfsubjectrulesreviews
	//        Resource Name: subjectaccessreviews
	//Group Version: autoscaling/v1
	//        Resource Name: horizontalpodautoscalers
	//        Resource Name: horizontalpodautoscalers/status
	//Group Version: autoscaling/v2beta1
	//        Resource Name: horizontalpodautoscalers
	//        Resource Name: horizontalpodautoscalers/status
	//Group Version: autoscaling/v2beta2
	//        Resource Name: horizontalpodautoscalers
	//        Resource Name: horizontalpodautoscalers/status
	//Group Version: batch/v1
	//        Resource Name: jobs
	//        Resource Name: jobs/status
	//Group Version: batch/v1beta1
	//        Resource Name: cronjobs
	//        Resource Name: cronjobs/status
	//Group Version: certificates.k8s.io/v1
	//        Resource Name: certificatesigningrequests
	//        Resource Name: certificatesigningrequests/approval
	//        Resource Name: certificatesigningrequests/status
	//Group Version: certificates.k8s.io/v1beta1
	//        Resource Name: certificatesigningrequests
	//        Resource Name: certificatesigningrequests/approval
	//        Resource Name: certificatesigningrequests/status
	//Group Version: networking.k8s.io/v1
	//        Resource Name: ingressclasses
	//        Resource Name: ingresses
	//        Resource Name: ingresses/status
	//        Resource Name: networkpolicies
	//Group Version: networking.k8s.io/v1beta1
	//        Resource Name: ingressclasses
	//        Resource Name: ingresses
	//        Resource Name: ingresses/status
	//Group Version: policy/v1beta1
	//        Resource Name: poddisruptionbudgets
	//        Resource Name: poddisruptionbudgets/status
	//        Resource Name: podsecuritypolicies
	//Group Version: rbac.authorization.k8s.io/v1
	//        Resource Name: clusterrolebindings
	//        Resource Name: clusterroles
	//        Resource Name: rolebindings
	//        Resource Name: roles
	//Group Version: rbac.authorization.k8s.io/v1beta1
	//        Resource Name: clusterrolebindings
	//        Resource Name: clusterroles
	//        Resource Name: rolebindings
	//        Resource Name: roles
	//Group Version: storage.k8s.io/v1
	//        Resource Name: csidrivers
	//        Resource Name: csinodes
	//        Resource Name: storageclasses
	//        Resource Name: volumeattachments
	//        Resource Name: volumeattachments/status
	//Group Version: storage.k8s.io/v1beta1
	//        Resource Name: csidrivers
	//        Resource Name: csinodes
	//        Resource Name: storageclasses
	//        Resource Name: volumeattachments
	//Group Version: admissionregistration.k8s.io/v1
	//        Resource Name: mutatingwebhookconfigurations
	//        Resource Name: validatingwebhookconfigurations
	//Group Version: admissionregistration.k8s.io/v1beta1
	//        Resource Name: mutatingwebhookconfigurations
	//        Resource Name: validatingwebhookconfigurations
	//Group Version: apiextensions.k8s.io/v1
	//        Resource Name: customresourcedefinitions
	//        Resource Name: customresourcedefinitions/status
	//Group Version: apiextensions.k8s.io/v1beta1
	//        Resource Name: customresourcedefinitions
	//        Resource Name: customresourcedefinitions/status
	//Group Version: scheduling.k8s.io/v1
	//        Resource Name: priorityclasses
	//Group Version: scheduling.k8s.io/v1beta1
	//        Resource Name: priorityclasses
	//Group Version: coordination.k8s.io/v1
	//        Resource Name: leases
	//Group Version: coordination.k8s.io/v1beta1
	//        Resource Name: leases
	//Group Version: node.k8s.io/v1beta1
	//        Resource Name: runtimeclasses
	//Group Version: discovery.k8s.io/v1beta1
	//        Resource Name: endpointslices

}
