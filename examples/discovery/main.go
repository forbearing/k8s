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

func main() {
	discoveryClient := k8s.DiscoveryClientOrDie("")
	_, allAPIResourceList, err := discoveryClient.ServerGroupsAndResources()
	if err != nil {
		log.Fatal(err)
	}

	//fmt.Println(discoveryClient.ServerGroups())
	//fmt.Println()
	//fmt.Println(discoveryClient.ServerGroupsAndResources())
	//fmt.Println()
	//fmt.Println(discoveryClient.ServerPreferredNamespacedResources())
	//fmt.Println(discoveryClient.ServerPreferredResources())
	//fmt.Println()
	//fmt.Println(discoveryClient.ServerResourcesForGroupVersion("apps/v1"))

	fmt.Println(IsNamespaced(allAPIResourceList, clusterrole.GVR()))           // false
	fmt.Println(IsNamespaced(allAPIResourceList, clusterrolebinding.GVR()))    // false
	fmt.Println(IsNamespaced(allAPIResourceList, configmap.GVR()))             // true
	fmt.Println(IsNamespaced(allAPIResourceList, cronjob.GVR()))               // true
	fmt.Println(IsNamespaced(allAPIResourceList, daemonset.GVR()))             // true
	fmt.Println(IsNamespaced(allAPIResourceList, deployment.GVR()))            // true
	fmt.Println(IsNamespaced(allAPIResourceList, ingress.GVR()))               // true
	fmt.Println(IsNamespaced(allAPIResourceList, ingressclass.GVR()))          // false
	fmt.Println(IsNamespaced(allAPIResourceList, job.GVR()))                   // true
	fmt.Println(IsNamespaced(allAPIResourceList, namespace.GVR()))             // false
	fmt.Println(IsNamespaced(allAPIResourceList, networkpolicy.GVR()))         // true
	fmt.Println(IsNamespaced(allAPIResourceList, node.GVR()))                  // false
	fmt.Println(IsNamespaced(allAPIResourceList, persistentvolume.GVR()))      // false
	fmt.Println(IsNamespaced(allAPIResourceList, persistentvolumeclaim.GVR())) // true
	fmt.Println(IsNamespaced(allAPIResourceList, pod.GVR()))                   // true
	fmt.Println(IsNamespaced(allAPIResourceList, replicaset.GVR()))            // true
	fmt.Println(IsNamespaced(allAPIResourceList, replicationcontroller.GVR())) // true
	fmt.Println(IsNamespaced(allAPIResourceList, role.GVR()))                  // true
	fmt.Println(IsNamespaced(allAPIResourceList, rolebinding.GVR()))           // true
	fmt.Println(IsNamespaced(allAPIResourceList, secret.GVR()))                // true
	fmt.Println(IsNamespaced(allAPIResourceList, service.GVR()))               // true
	fmt.Println(IsNamespaced(allAPIResourceList, serviceaccount.GVR()))        // true
	fmt.Println(IsNamespaced(allAPIResourceList, statefulset.GVR()))           // true
	fmt.Println(IsNamespaced(allAPIResourceList, storageclass.GVR()))          // false

}

func IsNamespaced(allAPIResourceList []*metav1.APIResourceList, gvr schema.GroupVersionResource) bool {
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
