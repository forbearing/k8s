package main

import (
	"github.com/forbearing/k8s/pod"
	corev1 "k8s.io/api/core/v1"
)

func Pod_List() {
	handler, err := pod.New(ctx, kubeconfig, namespace)
	if err != nil {
		panic(err)
	}
	defer cleanup(handler)

	podList, err := handler.WithNamespace("kube-system").ListByLabel("k8s-app=kube-dns")
	checkErr("ListByLabel", outputPods(podList), err)

	podList2, err := handler.ListByNamespace("kube-system")
	checkErr("ListByNamespace", outputPods(podList2), err)

	podList3, err := handler.ListAll()
	checkErr("ListAll", outputPods(podList3), err)

	podList4, err := handler.ListByNode("d11-k8s-master1")
	checkErr("ListByNode", outputPods(podList4), err)

	podList5, err := handler.WithNamespace("").ListByField("spec.nodeName=d11-k8s-master1")
	checkErr("ListByField", outputPods(podList5), err)

	// Ouptut:

	//2022/08/08 22:51:24 ListByLabel success: [coredns-6fbf8b5fd4-wslqx].
	//2022/08/08 22:51:24 ListByNamespace success: [calico-kube-controllers-bfdd697d7-87fmm calico-node-8wzlg calico-node-hnxbd calico-node-w8qsn ......]
	//2022/08/08 22:51:24 ListAll success: [cert-manager-6544c44c6b-dzkrl cert-manager-cainjector-5687864d5f-vqr82 cert-manager-webhook-785bb86798-fk6lf ......]
	//2022/08/08 22:51:24 ListByNode success: [calico-node-x6b89 metrics-server-74bdd7786d-h67rw].
	//2022/08/08 22:51:24 ListByField success: [calico-node-x6b89 metrics-server-74bdd7786d-h67rw].
}

func outputPods(podList []*corev1.Pod) []string {
	if podList == nil {
		return nil
	}
	var pl []string
	for _, pod := range podList {
		pl = append(pl, pod.Name)
	}
	return pl
}
