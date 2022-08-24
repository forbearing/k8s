package main

import (
	"log"

	"github.com/forbearing/k8s"
	"github.com/forbearing/k8s/replicationcontroller"
	corev1 "k8s.io/api/core/v1"
)

func ReplicationController_Tools() {
	handler := replicationcontroller.NewOrDie(ctx, kubeconfig, namespace)
	//defer cleanup(handler)

	if _, err := handler.Apply(filename); err != nil {
		log.Fatal(err)
	}
	if err := k8s.ApplyF(ctx, "", filename2); err != nil {
		log.Fatal(err)
	}
	log.Println(handler.IsReady(name))  // false
	log.Println(handler.IsReady(name2)) // false
	handler.WaitReady(name)             // block until the replicationcontroller ready
	handler.WaitReady(name2)            // block until the replicationcontroller ready
	log.Println(handler.IsReady(name))  // true
	log.Println(handler.IsReady(name2)) // true

	//GetPods
	podList, err := handler.GetPods(name2)
	checkErr("GetPods", outputPods(podList), err)

	//GetPVC
	pvcList, err := handler.GetPVC(name2)
	checkErr("GetPVC", pvcList, err)

	//GetPV
	pvList, err := handler.GetPV(name2)
	checkErr("GtPV", pvList, err)
	handler.Delete(name)

	//Output:

	//2022/08/24 17:58:10 false
	//2022/08/24 17:58:10 false
	//2022/08/24 17:58:16 true
	//2022/08/24 17:58:16 true
	//2022/08/24 17:58:16 GetPods success: [nginx-rc-9swl4 nginx-rc-pgm6g nginx-rc-pthb7]
	//2022/08/24 17:58:16 GetPVC success: [rc-nginx-data rc-nginx-html]
	//2022/08/24 17:58:16 GtPV success: [pvc-29bd4c74-e464-4ed7-a55a-7dbcea428325 pvc-ed802c83-18eb-4728-8a26-ec760717701e]
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
