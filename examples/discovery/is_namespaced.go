package main

import (
	"fmt"
	"log"

	"github.com/forbearing/k8s"
	"github.com/forbearing/k8s/clusterrole"
	"github.com/forbearing/k8s/clusterrolebinding"
	"github.com/forbearing/k8s/configmap"
	"github.com/forbearing/k8s/cronjob"
	"github.com/forbearing/k8s/daemonset"
	"github.com/forbearing/k8s/deployment"
	"github.com/forbearing/k8s/ingress"
	"github.com/forbearing/k8s/ingressclass"
	"github.com/forbearing/k8s/job"
	"github.com/forbearing/k8s/namespace"
	"github.com/forbearing/k8s/networkpolicy"
	"github.com/forbearing/k8s/node"
	"github.com/forbearing/k8s/persistentvolume"
	"github.com/forbearing/k8s/persistentvolumeclaim"
	"github.com/forbearing/k8s/pod"
	"github.com/forbearing/k8s/replicaset"
	"github.com/forbearing/k8s/replicationcontroller"
	"github.com/forbearing/k8s/role"
	"github.com/forbearing/k8s/rolebinding"
	"github.com/forbearing/k8s/secret"
	"github.com/forbearing/k8s/service"
	"github.com/forbearing/k8s/serviceaccount"
	"github.com/forbearing/k8s/statefulset"
	"github.com/forbearing/k8s/storageclass"
	"github.com/forbearing/k8s/util/conversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func Is_Namespaced() {
	discoveryClient := k8s.DiscoveryClientOrDie("")
	_, allAPIResourceList, err := discoveryClient.ServerGroupsAndResources()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(isNamespaced(allAPIResourceList, clusterrole.GVR()))           // false
	fmt.Println(isNamespaced(allAPIResourceList, clusterrolebinding.GVR()))    // false
	fmt.Println(isNamespaced(allAPIResourceList, configmap.GVR()))             // true
	fmt.Println(isNamespaced(allAPIResourceList, cronjob.GVR()))               // true
	fmt.Println(isNamespaced(allAPIResourceList, daemonset.GVR()))             // true
	fmt.Println(isNamespaced(allAPIResourceList, deployment.GVR()))            // true
	fmt.Println(isNamespaced(allAPIResourceList, ingress.GVR()))               // true
	fmt.Println(isNamespaced(allAPIResourceList, ingressclass.GVR()))          // false
	fmt.Println(isNamespaced(allAPIResourceList, job.GVR()))                   // true
	fmt.Println(isNamespaced(allAPIResourceList, namespace.GVR()))             // false
	fmt.Println(isNamespaced(allAPIResourceList, networkpolicy.GVR()))         // true
	fmt.Println(isNamespaced(allAPIResourceList, node.GVR()))                  // false
	fmt.Println(isNamespaced(allAPIResourceList, persistentvolume.GVR()))      // false
	fmt.Println(isNamespaced(allAPIResourceList, persistentvolumeclaim.GVR())) // true
	fmt.Println(isNamespaced(allAPIResourceList, pod.GVR()))                   // true
	fmt.Println(isNamespaced(allAPIResourceList, replicaset.GVR()))            // true
	fmt.Println(isNamespaced(allAPIResourceList, replicationcontroller.GVR())) // true
	fmt.Println(isNamespaced(allAPIResourceList, role.GVR()))                  // true
	fmt.Println(isNamespaced(allAPIResourceList, rolebinding.GVR()))           // true
	fmt.Println(isNamespaced(allAPIResourceList, secret.GVR()))                // true
	fmt.Println(isNamespaced(allAPIResourceList, service.GVR()))               // true
	fmt.Println(isNamespaced(allAPIResourceList, serviceaccount.GVR()))        // true
	fmt.Println(isNamespaced(allAPIResourceList, statefulset.GVR()))           // true
	fmt.Println(isNamespaced(allAPIResourceList, storageclass.GVR()))          // false

}

func isNamespaced(allAPIResourceList []*metav1.APIResourceList, gvr schema.GroupVersionResource) bool {
	for _, apiResourceList := range allAPIResourceList {
		for _, apiResource := range apiResourceList.APIResources {
			//fmt.Println(apiResourceList.GroupVersion, gvr.GroupVersion().String())
			//fmt.Println(apiResource.Kind, conversion.ResourceToKind(gvr.Resource))
			if apiResourceList.GroupVersion == gvr.GroupVersion().String() && apiResource.Kind == conversion.ResourceToKind(gvr.Resource) {
				return apiResource.Namespaced
			}
		}
	}
	log.Println("not found")
	fmt.Println()
	fmt.Println()
	return false
}
