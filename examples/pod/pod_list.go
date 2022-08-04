package main

import (
	"fmt"
	"log"

	"github.com/forbearing/k8s/pod"
	corev1 "k8s.io/api/core/v1"
)

func Pod_List() {
	handler, err := pod.New(ctx, kubeconfig, namespace)
	if err != nil {
		panic(err)
	}
	defer cleanup(handler)

	podList, err := handler.ListAll()
	checkErr("ListAll", "", err)
	outputPods(podList)
	fmt.Println()

	podList2, err := handler.ListByNamespace("default")
	checkErr("ListByNamespace", "", err)
	outputPods(podList2)
	fmt.Println()

	podList3, err := handler.WithNamespace("default").ListByLabel("")
	checkErr("ListByLabel", "", err)
	outputPods(podList3)
}

//func outputPods(podList *corev1.PodList) {
//    var pl []string
//    for _, pod := range podList.Items {
//        pl = append(pl, pod.Name)
//    }
//    log.Println(pl)
//}
func outputPods(podList []*corev1.Pod) {
	var pl []string
	for _, pod := range podList {
		pl = append(pl, pod.Name)
	}
	log.Println(pl)
	//for _, pod := range podList {
	//    log.Println(pod.Name)
	//}
}
