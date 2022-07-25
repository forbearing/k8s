package main

import (
	"log"
	"time"

	"github.com/forbearing/k8s"
	"github.com/forbearing/k8s/statefulset"
	corev1 "k8s.io/api/core/v1"
)

func StatefulSet_Tools() {
	handler, err := statefulset.New(ctx, kubeconfig, namespace)
	if err != nil {
		panic(err)
	}
	handler.Apply(filename)
	k8s.ApplyF(ctx, filename2, kubeconfig)

	sts, err := handler.Get(name)
	if err != nil {
		panic(err)
	}
	sts2, err := handler.Get(name2)
	if err != nil {
		panic(err)
	}

	getByName := func() {
		log.Println("===== Get StatefulSet By Name")
		begin := time.Now()

		pods, err := handler.GetPods(name)
		checkErr("GetPods", printPods(pods), err)
		pvcList, err := handler.GetPVC(name2)
		checkErr("GetPVC", pvcList, err)
		pvList, err := handler.GetPV(name2)
		checkErr("GetPV", pvList, err)
		containers, err := handler.GetContainers(name)
		checkErr("GetContainers", containers, err)
		images, err := handler.GetImages(name)
		checkErr("GetImages", images, err)

		end := time.Now()
		log.Println("===== Get StatefulSet By Name Cost Time:", end.Sub(begin))
		log.Println()
	}
	getByObj := func() {
		log.Println("===== Get StatefulSet By Object")
		begin := time.Now()

		pods, err := handler.GetPods(sts)
		checkErr("GetPods", printPods(pods), err)
		pvcList, err := handler.GetPVC(sts2)
		checkErr("GetPVC", pvcList, err)
		pvList, err := handler.GetPV(sts2)
		checkErr("GetPV", pvList, err)
		containers, err := handler.GetContainers(sts)
		checkErr("GetContainers", containers, err)
		images, err := handler.GetImages(sts)
		checkErr("GetImages", images, err)

		end := time.Now()
		log.Println("===== Get StatefulSet By Object Cost Time:", end.Sub(begin))
	}

	getByName()
	time.Sleep(time.Second * 5)
	getByObj()

	// Output:

	//2022/07/11 10:51:44 ===== Get StatefulSet By Name
	//2022/07/11 10:51:44 GetPods success: [mysts-0 mysts-1 mysts-2]
	//2022/07/11 10:51:44 GetPVC success: [data-nginx-sts-0 data-nginx-sts-1 data-nginx-sts-2 data-nginx-sts-3 html-nginx-sts-0 html-nginx-sts-1 html-nginx-sts-2 html-nginx-sts-3]
	//2022/07/11 10:51:44 GetPV success: [pvc-42d9110e-e0b2-476a-8905-6cbab69bf6c9 pvc-f2f9e388-9460-43da-a141-0af8f7a4719b pvc-3ea4cc9e-5dcd-4913-ac5a-52a9151c1bd3 pvc-bd0e819c-13f1-4eb1-8d3f-b36b13381ef7 pvc-ac6c44e2-4411-45d2-acb9-137d8344e231 pvc-4bac3464-3bf6-4583-b52b-dbe462a66b5d pvc-4764fac1-a1fa-46bb-bafd-08f2a486ad11 pvc-5186d5f9-e184-4e9e-ac2b-b8c63e700ef6]
	//2022/07/11 10:51:44 GetContainers success: [nginx busybox]
	//2022/07/11 10:51:44 GetImages success: [nginx:v1.20 busybox]
	//2022/07/11 10:51:44 ===== Get StatefulSet By Name Cost Time: 242.682331ms
	//2022/07/11 10:51:44
	//2022/07/11 10:51:49 ===== Get StatefulSet By Object
	//2022/07/11 10:51:49 GetPods success: [mysts-0 mysts-1 mysts-2]
	//2022/07/11 10:51:49 GetPVC success: [data-nginx-sts-0 data-nginx-sts-1 data-nginx-sts-2 data-nginx-sts-3 html-nginx-sts-0 html-nginx-sts-1 html-nginx-sts-2 html-nginx-sts-3]
	//2022/07/11 10:51:49 GetPV success: [pvc-42d9110e-e0b2-476a-8905-6cbab69bf6c9 pvc-f2f9e388-9460-43da-a141-0af8f7a4719b pvc-3ea4cc9e-5dcd-4913-ac5a-52a9151c1bd3 pvc-bd0e819c-13f1-4eb1-8d3f-b36b13381ef7 pvc-ac6c44e2-4411-45d2-acb9-137d8344e231 pvc-4bac3464-3bf6-4583-b52b-dbe462a66b5d pvc-4764fac1-a1fa-46bb-bafd-08f2a486ad11 pvc-5186d5f9-e184-4e9e-ac2b-b8c63e700ef6]
	//2022/07/11 10:51:49 GetContainers success: [nginx busybox]
	//2022/07/11 10:51:49 GetImages success: [nginx:v1.20 busybox]
	//2022/07/11 10:51:49 ===== Get StatefulSet By Object Cost Time: 106.91293ms
}

func printPods(podList []corev1.Pod) []string {
	var pl []string
	for _, pod := range podList {
		pl = append(pl, pod.Name)
	}
	return pl
}
