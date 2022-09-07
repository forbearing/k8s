package main

import (
	"log"
	"time"

	"github.com/forbearing/k8s"
	"github.com/forbearing/k8s/pod"
)

func Pod_Tools() {
	handler, err := pod.New(ctx, kubeconfig, namespace)
	if err != nil {
		panic(err)
	}
	defer cleanup(handler)

	handler.Apply(filename)
	k8s.ApplyF(ctx, kubeconfig, filename2, namespace, k8s.IgnoreInvalid)

	log.Println(handler.IsReady(name))  // false
	log.Println(handler.IsReady(name2)) // false
	handler.WaitReady(name)
	handler.WaitReady(name2)
	log.Println(handler.IsReady(name))  // true
	log.Println(handler.IsReady(name2)) // true

	pod, err := handler.Get(name)
	if err != nil {
		panic(err)
	}
	pod2, err := handler.Get(name2)
	if err != nil {
		panic(err)
	}
	_, _ = pod, pod2

	//Pod_Execute(handler)

	pvList1, err := handler.GetPV(name2)
	checkErr("GetPV", pvList1, err)
	pvList2, err := handler.GetPV(pod2)
	checkErr("GetPV", pvList2, err)
	log.Println()

	getByName := func() {
		log.Println("===== Get Pod By Name")
		begin := time.Now()

		uid, err := handler.GetUID(name)
		checkErr("GetUID", uid, err)
		podIP, err := handler.GetIP(name)
		checkErr("GetIP", podIP, err)
		nodeIP, err := handler.GetNodeIP(name)
		checkErr("GetNodeIP", nodeIP, err)
		nodeName, err := handler.GetNodeName(name)
		checkErr("GetNodeName", nodeName, err)
		age, err := handler.GetAge(name)
		checkErr("GetAge", age, err)
		status, err := handler.GetStatus(name)
		checkErr("GetStatus", status, err)
		qos, err := handler.GetQosClass(name)
		checkErr("GetQosClass", qos, err)
		pvcList, err := handler.GetPVC(name2)
		checkErr("GetPVC", pvcList, err)
		controller, err := handler.GetController(name)
		checkErr("GetContainers", controller, err)
		containers, err := handler.GetContainers(name)
		checkErr("GetContainers", containers, err)
		initContainers, err := handler.GetInitContainers(name)
		checkErr("GetInitContainers", initContainers, err)
		readyContainers, err := handler.GetReadyContainers(name)
		checkErr("GetReadyContainers", readyContainers, err)

		end := time.Now()
		log.Println("===== Get Pod By Name Cost Time:", end.Sub(begin))
		log.Println()
	}

	getByObj := func() {
		log.Println("===== Get Pod By Object")
		begin := time.Now()

		uid, err := handler.GetUID(pod)
		checkErr("GetUID", uid, err)
		podIP, err := handler.GetIP(pod)
		checkErr("GetIP", podIP, err)
		nodeIP, err := handler.GetNodeIP(pod)
		checkErr("GetNodeIP", nodeIP, err)
		nodeName, err := handler.GetNodeName(pod)
		checkErr("GetNodeName", nodeName, err)
		age, err := handler.GetAge(pod)
		checkErr("GetAge", age, err)
		status, err := handler.GetStatus(pod)
		checkErr("GetStatus", status, err)
		qos, err := handler.GetQosClass(pod)
		checkErr("GetQosClass", qos, err)
		pvcList, err := handler.GetPVC(pod2)
		checkErr("GetPVC", pvcList, err)
		controller, err := handler.GetController(pod)
		checkErr("GetContainers", controller, err)
		containers, err := handler.GetContainers(pod)
		checkErr("GetContainers", containers, err)
		initContainers, err := handler.GetInitContainers(pod)
		checkErr("GetInitContainers", initContainers, err)
		readyContainers, err := handler.GetReadyContainers(pod)
		checkErr("GetReadyContainers", readyContainers, err)

		end := time.Now()
		log.Println("===== Get Pod By Object Cost Time:", end.Sub(begin))
	}

	getByName()
	getByObj()

	// Output:

	//2022/07/11 18:09:20 false
	//2022/07/11 18:09:20 false
	//2022/07/11 18:09:43 true
	//2022/07/11 18:09:43 true
	//2022/07/11 18:09:44 GetPV success: [pvc-2cd03beb-26df-48a3-8ca2-a35a759fff65 pvc-92618243-68ee-46af-ac6f-72f99727a9ca pvc-4820754c-4954-4f46-92f2-f5387fcad452 pvc-f81ecb29-570b-4986-8cc1-b0dbc0b930a7 pvc-76e2b987-4e25-460b-b879-e46072706ff1].
	//2022/07/11 18:09:45 GetPV success: [pvc-2cd03beb-26df-48a3-8ca2-a35a759fff65 pvc-92618243-68ee-46af-ac6f-72f99727a9ca pvc-4820754c-4954-4f46-92f2-f5387fcad452 pvc-f81ecb29-570b-4986-8cc1-b0dbc0b930a7 pvc-76e2b987-4e25-460b-b879-e46072706ff1].
	//2022/07/11 18:09:45
	//2022/07/11 18:09:45 ===== Get Pod By Name
	//2022/07/11 18:09:45 GetUID success: f5b5ccc0-053d-42f4-8783-0ce0be5875ca.
	//2022/07/11 18:09:45 GetIP success: 192.169.172.32.
	//2022/07/11 18:09:45 GetNodeIP success: 10.240.2.22.
	//2022/07/11 18:09:45 GetNodeName success: sh-u18-k8s-cicd-worker2.
	//2022/07/11 18:09:46 GetAge success: 26.102993s.
	//2022/07/11 18:09:46 GetStatus success: Running.
	//2022/07/11 18:09:46 GetQosClass success: BestEffort.
	//2022/07/11 18:09:46 GetPVC success: [pod-k8s-tools-data-rbd pod-k8s-tools-data-cephfs pod-k8s-tools-data-nfs pod-nginx-data pod-nginx-html].
	//2022/07/11 18:09:46 GetContainers failed: the pod "mypod" doesn't have controller
	//2022/07/11 18:09:47 GetContainers success: [{nginx nginx}].
	//2022/07/11 18:09:47 GetInitContainers success: [{busybox busybox}].
	//2022/07/11 18:09:47 GetReadyContainers success: [{nginx nginx:latest}].
	//2022/07/11 18:09:47 ===== Get Pod By Name Cost Time: 2.40130122s
	//2022/07/11 18:09:47
	//2022/07/11 18:09:47 ===== Get Pod By Object
	//2022/07/11 18:09:47 GetUID success: f5b5ccc0-053d-42f4-8783-0ce0be5875ca.
	//2022/07/11 18:09:47 GetIP success: 192.169.172.32.
	//2022/07/11 18:09:47 GetNodeIP success: 10.240.2.22.
	//2022/07/11 18:09:47 GetNodeName success: sh-u18-k8s-cicd-worker2.
	//2022/07/11 18:09:47 GetAge success: 27.504643s.
	//2022/07/11 18:09:47 GetStatus success: Running.
	//2022/07/11 18:09:47 GetQosClass success: BestEffort.
	//2022/07/11 18:09:47 GetPVC success: [pod-k8s-tools-data-rbd pod-k8s-tools-data-cephfs pod-k8s-tools-data-nfs pod-nginx-data pod-nginx-html].
	//2022/07/11 18:09:47 GetContainers failed: the pod "mypod" doesn't have controller
	//2022/07/11 18:09:47 GetContainers success: [{nginx nginx}].
	//2022/07/11 18:09:47 GetInitContainers success: [{busybox busybox}].
	//2022/07/11 18:09:47 GetReadyContainers success: [{nginx nginx:latest}].
	//2022/07/11 18:09:47 ===== Get Pod By Object Cost Time: 49.573µs
}
