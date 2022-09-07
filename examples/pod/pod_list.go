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

	podList6, err := handler.ListRunning()
	checkErr("ListRunning", outputPods(podList6), err)

	podList7, err := handler.ListSucceeded()
	checkErr("ListSucceeded", outputPods(podList7), err)

	podList8, err := handler.ListFailed()
	checkErr("ListFailed", outputPods(podList8), err)

	podList9, err := handler.ListPending()
	checkErr("ListPending", outputPods(podList9), err)

	podList10, err := handler.ListUnknow()
	checkErr("ListUnknow", outputPods(podList10), err)

	// Ouptut:

	//022/09/07 18:37:59 ListByLabel success: [coredns-64897985d-467q7 coredns-64897985d-tgjph].
	//2022/09/07 18:37:59 ListByNamespace success: [coredns-64897985d-467q7 coredns-64897985d-tgjph ...]
	//2022/09/07 18:37:59 ListAll success: [nginx-85b98978db-tf4w7 coredns-64897985d-467q7 coredns-64897985d-tgjph ...]
	//2022/09/07 18:37:59 ListByNode success: [].
	//2022/09/07 18:37:59 ListByField success: [].
	//2022/09/07 18:37:59 ListRunning success: [nginx-85b98978db-tf4w7 coredns-64897985d-467q7 coredns-64897985d-tgjph ...]
	//2022/09/07 18:37:59 ListSucceeded success: [].
	//2022/09/07 18:37:59 ListFailed success: [].
	//2022/09/07 18:37:59 ListPending success: [].
	//2022/09/07 18:37:59 ListUnknow success: [].

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
