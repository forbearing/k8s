package main

import (
	"log"

	"github.com/forbearing/k8s/pod"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type updateObject struct {
	oldObj interface{}
	newObj interface{}
}

func Pod_Informer() {
	handler, err := pod.New(ctx, kubeconfig, namespace)
	if err != nil {
		panic(err)
	}
	defer cleanup(handler)

	//stopCh := make(chan struct{}, 1)
	//handler.TestInformer(stopCh)

	addQueue := make(chan interface{}, 100)
	updateQueue := make(chan updateObject, 100)
	deleteQueue := make(chan interface{}, 100)

	addFunc := func(obj interface{}) { addQueue <- obj }
	updateFunc := func(oldObj, newObj interface{}) {
		uo := updateObject{}
		uo.oldObj = oldObj
		uo.newObj = newObj
		updateQueue <- uo
	}
	deleteFunc := func(obj interface{}) { deleteQueue <- obj }
	stopCh := make(chan struct{}, 1)

	// RunInformer 必须开启一个新的 goroutine 来执行
	go func() {
		handler.RunInformer(addFunc, updateFunc, deleteFunc, stopCh)
	}()

	for {
		select {
		case obj := <-addQueue:
			myObj := obj.(metav1.Object)
			log.Printf("New Pod Added to Store: %s\n", myObj.GetName())
		case uo := <-updateQueue:
			// OnUpdate is called with all existing objects on the specific resync interval
			// OnUpdate is called with all existing objects when a watch connection
			// is dropped by the server and a full relist is required

			oldPod := uo.oldObj.(*corev1.Pod)
			newPod := uo.newObj.(*corev1.Pod)
			if oldPod.ResourceVersion != newPod.ResourceVersion {
				log.Printf("Pod Updated to Store: %s\n", newPod.Name)
			}

			//if !reflect.DeepEqual(uo.oldObj, uo.newObj) {
			//    log.Printf("Pod Updated to Store: %s\n", uo.newObj.(*corev1.Pod).Name)
			//}
		case obj := <-deleteQueue:
			myObj := obj.(metav1.Object)
			log.Printf("Pod Deleted from Store: %s\n", myObj.GetName())
		case <-stopCh:
			log.Println("informer stopped")
			return
		}
	}

	// Output:

	//2022/07/26 14:52:33 New Pod Added to Store: prometheus-adapter-7b5bfc99dc-nwpv8
	//2022/07/26 14:52:33 New Pod Added to Store: prometheus-operator-75d9b475d9-7nml6
	//2022/07/26 14:52:33 New Pod Added to Store: metrics-server-c8bb454c8-k9lb6
	//2022/07/26 14:52:33 New Pod Added to Store: kubeapps-internal-kubeappsapis-6f8bb67f5f-wfmfh
	//2022/07/26 14:52:33 New Pod Added to Store: node-exporter-mvmb6
	//2022/07/26 14:52:33 New Pod Added to Store: haproxy-exporter-zdqfs
	//2022/07/26 14:52:33 New Pod Added to Store: node-exporter-p97rv
	//2022/07/26 14:52:33 New Pod Added to Store: calico-node-lk256
	//2022/07/26 14:52:33 New Pod Added to Store: metrics-server-c8bb454c8-ndkh7
	//2022/07/26 14:52:33 New Pod Added to Store: apprepo-kubeapps-sync-bitnami-27646970-xmrhw
	//2022/07/26 14:52:33 New Pod Added to Store: local-path-provisioner-cc67d8db7-c9t62
	//2022/07/26 14:52:33 New Pod Added to Store: ingress-controller-ingress-nginx-controller-zqs8d
	//2022/07/26 14:52:33 New Pod Added to Store: kubeapps-internal-dashboard-5f8b767f86-2lxxq
	//2022/07/26 14:52:33 New Pod Added to Store: nginx-exporter-cicd-698c8dcccd-4qhlb
	//2022/07/26 14:52:33 New Pod Added to Store: kubeapps-75c5c49c8c-fdqhc
	//2022/07/26 14:52:33 New Pod Added to Store: kubeapps-internal-apprepository-controller-648f4f6494-pqx4q
	//2022/07/26 14:52:33 New Pod Added to Store: kube-state-metrics-5cf84d7f85-nc9xh
	//2022/07/26 14:52:33 New Pod Added to Store: prometheus-64fd8ccd65-78jwc
	//2022/07/26 14:52:33 New Pod Added to Store: calico-node-qh4c6
	//2022/07/26 14:52:33 New Pod Added to Store: kubeapps-75c5c49c8c-slzkv
	//2022/07/26 14:52:33 New Pod Added to Store: k8s-tools-7b9697899c-rfskk
	//2022/07/26 14:52:33 New Pod Added to Store: alertmanager-main-2
	//2022/07/26 14:52:33 New Pod Added to Store: node-exporter-2cgmk
	//2022/07/26 14:52:33 New Pod Added to Store: kubeapps-internal-kubeappsapis-6f8bb67f5f-44wzn
	//2022/07/26 14:52:33 New Pod Added to Store: node-exporter-gx98t
	//2022/07/26 14:52:33 New Pod Added to Store: node-exporter-n7pzr
	//2022/07/26 14:52:33 New Pod Added to Store: nginx-logs
	//2022/07/26 14:52:33 New Pod Added to Store: calico-node-hvd48
	//2022/07/26 14:52:33 New Pod Added to Store: calico-node-wn55k
	//2022/07/26 14:52:33 New Pod Added to Store: kubeapps-internal-dashboard-5f8b767f86-vmk9t
	//2022/07/26 14:52:33 New Pod Added to Store: calico-node-vtq5p
	//2022/07/26 14:52:33 New Pod Added to Store: nfs-provisioner-nfs-subdir-external-provisioner-5fc895f4d9fb2xb
	//2022/07/26 14:52:33 New Pod Added to Store: nfs-provisioner-nfs-subdir-external-provisioner-5fc895f4d9xr5v9
	//2022/07/26 14:52:33 New Pod Added to Store: gitlab-cf9cd9779-6jrr2
	//2022/07/26 14:52:33 New Pod Added to Store: openldap-0
	//2022/07/26 14:52:33 New Pod Added to Store: grafana-6ccd56f4b6-2cz6v
}
