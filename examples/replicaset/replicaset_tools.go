package main

import (
	"log"
	"time"

	"github.com/forbearing/k8s"
	"github.com/forbearing/k8s/replicaset"
	corev1 "k8s.io/api/core/v1"
)

func ReplicaSet_Tools() {
	handler := replicaset.NewOrDie(ctx, "", namespace)
	defer cleanup(handler)

	if _, err := handler.Apply(filename); err != nil {
		log.Fatal(err)
	}
	if err := k8s.ApplyF(ctx, "", filename2); err != nil {
		log.Fatal(err)
	}
	log.Println(handler.IsReady(name))  // false
	log.Println(handler.IsReady(name2)) // false
	handler.WaitReady(name)             // block until replicaset is ready
	handler.WaitReady(name2)            // block until replicaset is ready
	log.Println(handler.IsReady(name))  // true
	log.Println(handler.IsReady(name2)) // true

	// GetPods
	podList, err := handler.GetPods(name2)
	checkErr("GetPods", outputPods(podList), err)

	// GetPVC
	pvcList, err := handler.GetPVC(name2)
	checkErr("GetPVC", pvcList, err)

	// GetPV
	pvList, err := handler.GetPV(name2)
	checkErr("GetPV", pvList, err)

	time.Sleep(time.Second * 3)

	// Output:

	//2022/08/24 13:18:29 false
	//2022/08/24 13:18:29 false
	//2022/08/24 13:18:35 true
	//2022/08/24 13:18:35 true
	//2022/08/24 13:18:35 GetPods success: [nginx-rs-2fqjd nginx-rs-p7tmt nginx-rs-wk9kg]
	//2022/08/24 13:18:35 GetPVC success: [rs-nginx-data rs-nginx-html]
	//2022/08/24 13:18:35 GetPV success: [pvc-be2c8ca4-f1dc-480f-8519-4020248448bc pvc-b2e70135-4dcd-4343-8262-c5d51f6a38af]
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
