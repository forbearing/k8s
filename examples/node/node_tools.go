package main

import (
	"log"
	"time"

	"github.com/forbearing/k8s/node"
	corev1 "k8s.io/api/core/v1"
)

func Node_Tools() {
	handler, err := node.New(ctx, kubeconfig)
	if err != nil {
		panic(err)
	}
	masterObj, err := handler.Get(masterName)
	if err != nil {
		panic(err)
	}
	workerObj, err := handler.Get(workerName)
	if err != nil {
		panic(err)
	}

	getByName := func() {
		log.Println("===== Get Node By Name")
		begin := time.Now()

		//log.Println(handler.IsReady(masterName))
		//log.Println(handler.IsReady(workerName))
		log.Println(handler.IsMaster(masterName))
		log.Println(handler.IsMaster(workerName))
		log.Println(handler.IsControlPlane(masterName))
		log.Println(handler.IsControlPlane(workerName))

		//phase1, err := handler.GetPhase(masterName)
		//checkErr("GetPhase", err)
		//phase2, err := handler.GetPhase(workerName)
		//checkErr("GetPhase", err)
		//log.Println(phase1)
		//log.Println(phase2)

		hostname1, err := handler.GetHostname(masterName)
		checkErr("GetHostname", err)
		hostname2, err := handler.GetHostname(workerName)
		checkErr("GetHostname", err)
		log.Println(hostname1)
		log.Println(hostname2)

		internalIP1, err := handler.GetInternalIP(masterName)
		checkErr("GetInternalIP", err)
		internalIP2, err := handler.GetInternalIP(workerName)
		checkErr("GetInternalIP", err)
		log.Println(internalIP1)
		log.Println(internalIP2)

		roles1, err := handler.GetRoles(masterName)
		checkErr("GetRoles", err)
		roles2, err := handler.GetRoles(workerName)
		checkErr("GetRoles", err)
		log.Println(roles1)
		log.Println(roles2)

		podList1, err := handler.GetPods(masterName)
		checkErr("GetPods", err)
		podList2, err := handler.GetPods(workerName)
		checkErr("GetPods", err)
		printPods(podList1)
		printPods(podList2)

		cidr1, err := handler.GetCIDR(masterName)
		checkErr("GetCIDR", err)
		cidr2, err := handler.GetCIDR(workerName)
		checkErr("GetCIDR", err)
		log.Println(cidr1)
		log.Println(cidr2)

		cidrs1, err := handler.GetCIDRs(masterName)
		checkErr("GetCIDRs", err)
		cidrs2, err := handler.GetCIDRs(workerName)
		checkErr("GetCIDRs", err)
		log.Println(cidrs1)
		log.Println(cidrs2)

		masterInfo, err := handler.GetNodeInfo(masterName)
		checkErr("GetNodeInfo", err)
		workerInfo, err := handler.GetNodeInfo(workerName)
		checkErr("GetNodeInfo", err)
		log.Println(masterInfo)
		log.Println(workerInfo)

		age1, err := handler.GetAge(masterName)
		checkErr("GetAge", err)
		age2, err := handler.GetAge(workerName)
		checkErr("GetAge", err)
		log.Println(age1)
		log.Println(age2)

		end := time.Now()
		log.Println("===== Get Node By Name Cost Time:", end.Sub(begin))
		log.Println()
	}

	getByObj := func() {
		log.Println("===== Get Node By Object")
		begin := time.Now()

		log.Println(handler.IsMaster(masterObj))
		log.Println(handler.IsMaster(workerObj))
		log.Println(handler.IsControlPlane(masterObj))
		log.Println(handler.IsControlPlane(workerObj))

		//phase1, err := handler.GetPhase(masterObj)
		//checkErr("GetPhase", err)
		//phase2, err := handler.GetPhase(workerObj)
		//checkErr("GetPhase", err)
		//log.Println(phase1)
		//log.Println(phase2)

		hostname1, err := handler.GetHostname(masterObj)
		checkErr("GetHostname", err)
		hostname2, err := handler.GetHostname(workerObj)
		checkErr("GetHostname", err)
		log.Println(hostname1)
		log.Println(hostname2)

		internalIP1, err := handler.GetInternalIP(masterObj)
		checkErr("GetInternalIP", err)
		internalIP2, err := handler.GetInternalIP(workerObj)
		checkErr("GetInternalIP", err)
		log.Println(internalIP1)
		log.Println(internalIP2)

		roles1, err := handler.GetRoles(masterObj)
		checkErr("GetRoles", err)
		roles2, err := handler.GetRoles(workerObj)
		checkErr("GetRoles", err)
		log.Println(roles1)
		log.Println(roles2)

		podList1, err := handler.GetPods(masterObj)
		checkErr("GetPods", err)
		podList2, err := handler.GetPods(workerObj)
		checkErr("GetPods", err)
		printPods(podList1)
		printPods(podList2)

		cidr1, err := handler.GetCIDR(masterObj)
		checkErr("GetCIDR", err)
		cidr2, err := handler.GetCIDR(workerObj)
		checkErr("GetCIDR", err)
		log.Println(cidr1)
		log.Println(cidr2)

		cidrs1, err := handler.GetCIDRs(masterObj)
		checkErr("GetCIDRs", err)
		cidrs2, err := handler.GetCIDRs(workerObj)
		checkErr("GetCIDRs", err)
		log.Println(cidrs1)
		log.Println(cidrs2)

		masterInfo, err := handler.GetNodeInfo(masterObj)
		checkErr("GetNodeInfo", err)
		workerInfo, err := handler.GetNodeInfo(workerObj)
		checkErr("GetNodeInfo", err)
		log.Println(masterInfo)
		log.Println(workerInfo)

		age1, err := handler.GetAge(masterObj)
		checkErr("GetAge", err)
		age2, err := handler.GetAge(workerObj)
		checkErr("GetAge", err)
		log.Println(age1)
		log.Println(age2)

		end := time.Now()
		log.Println("===== Get Node By Object Cost Time:", end.Sub(begin))
		log.Println()
	}

	getByName()
	getByObj()
}

func printPods(podList []corev1.Pod) {
	var pl []string
	for _, pod := range podList {
		pl = append(pl, pod.Name)
	}
	log.Println(pl)
}
