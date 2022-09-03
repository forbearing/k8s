package main

import (
	"fmt"
	"log"

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
	utilrestmapper "github.com/forbearing/k8s/util/restmapper"
)

func Is_Namespaced() {
	restMapper, err := utilrestmapper.NewRESTMapper("")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(utilrestmapper.IsNamespaced(restMapper, clusterrole.GVK()))           // false
	fmt.Println(utilrestmapper.IsNamespaced(restMapper, clusterrolebinding.GVK()))    // false
	fmt.Println(utilrestmapper.IsNamespaced(restMapper, configmap.GVK()))             // true
	fmt.Println(utilrestmapper.IsNamespaced(restMapper, cronjob.GVK()))               // true
	fmt.Println(utilrestmapper.IsNamespaced(restMapper, daemonset.GVK()))             // true
	fmt.Println(utilrestmapper.IsNamespaced(restMapper, deployment.GVK()))            // true
	fmt.Println(utilrestmapper.IsNamespaced(restMapper, ingress.GVK()))               // true
	fmt.Println(utilrestmapper.IsNamespaced(restMapper, ingressclass.GVK()))          // false
	fmt.Println(utilrestmapper.IsNamespaced(restMapper, job.GVK()))                   // true
	fmt.Println(utilrestmapper.IsNamespaced(restMapper, namespace.GVK()))             // false
	fmt.Println(utilrestmapper.IsNamespaced(restMapper, networkpolicy.GVK()))         // true
	fmt.Println(utilrestmapper.IsNamespaced(restMapper, node.GVK()))                  // false
	fmt.Println(utilrestmapper.IsNamespaced(restMapper, persistentvolume.GVK()))      // false
	fmt.Println(utilrestmapper.IsNamespaced(restMapper, persistentvolumeclaim.GVK())) // true
	fmt.Println(utilrestmapper.IsNamespaced(restMapper, pod.GVK()))                   // true
	fmt.Println(utilrestmapper.IsNamespaced(restMapper, replicaset.GVK()))            // true
	fmt.Println(utilrestmapper.IsNamespaced(restMapper, replicationcontroller.GVK())) // true
	fmt.Println(utilrestmapper.IsNamespaced(restMapper, role.GVK()))                  // true
	fmt.Println(utilrestmapper.IsNamespaced(restMapper, rolebinding.GVK()))           // true
	fmt.Println(utilrestmapper.IsNamespaced(restMapper, secret.GVK()))                // true
	fmt.Println(utilrestmapper.IsNamespaced(restMapper, service.GVK()))               // true
	fmt.Println(utilrestmapper.IsNamespaced(restMapper, serviceaccount.GVK()))        // true
	fmt.Println(utilrestmapper.IsNamespaced(restMapper, statefulset.GVK()))           // true
	fmt.Println(utilrestmapper.IsNamespaced(restMapper, storageclass.GVK()))          // false
}
