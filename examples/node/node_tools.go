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
		checkErr("GetHostname", hostname1, err)
		hostname2, err := handler.GetHostname(workerName)
		checkErr("GetHostname", hostname2, err)

		internalIP1, err := handler.GetInternalIP(masterName)
		checkErr("GetInternalIP", internalIP1, err)
		internalIP2, err := handler.GetInternalIP(workerName)
		checkErr("GetInternalIP", internalIP2, err)

		roles1, err := handler.GetRoles(masterName)
		checkErr("GetRoles", roles1, err)
		roles2, err := handler.GetRoles(workerName)
		checkErr("GetRoles", roles2, err)

		podList1, err := handler.GetPods(masterName)
		checkErr("GetPods", outputPods(podList1), err)
		podList2, err := handler.GetPods(workerName)
		checkErr("GetPods", outputPods(podList2), err)

		cidr1, err := handler.GetCIDR(masterName)
		checkErr("GetCIDR", cidr1, err)
		cidr2, err := handler.GetCIDR(workerName)
		checkErr("GetCIDR", cidr2, err)

		cidrs1, err := handler.GetCIDRs(masterName)
		checkErr("GetCIDRs", cidrs1, err)
		cidrs2, err := handler.GetCIDRs(workerName)
		checkErr("GetCIDRs", cidrs2, err)

		masterInfo, err := handler.GetNodeInfo(masterName)
		checkErr("GetNodeInfo", masterInfo, err)
		workerInfo, err := handler.GetNodeInfo(workerName)
		checkErr("GetNodeInfo", workerInfo, err)

		age1, err := handler.GetAge(masterName)
		checkErr("GetAge", age1, err)
		age2, err := handler.GetAge(workerName)
		checkErr("GetAge", age2, err)

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
		checkErr("GetHostname", hostname1, err)
		hostname2, err := handler.GetHostname(workerObj)
		checkErr("GetHostname", hostname2, err)

		internalIP1, err := handler.GetInternalIP(masterObj)
		checkErr("GetInternalIP", internalIP1, err)
		internalIP2, err := handler.GetInternalIP(workerObj)
		checkErr("GetInternalIP", internalIP2, err)

		roles1, err := handler.GetRoles(masterObj)
		checkErr("GetRoles", roles1, err)
		roles2, err := handler.GetRoles(workerObj)
		checkErr("GetRoles", roles2, err)

		podList1, err := handler.GetPods(masterObj)
		checkErr("GetPods", outputPods(podList1), err)
		podList2, err := handler.GetPods(workerObj)
		checkErr("GetPods", outputPods(podList2), err)

		cidr1, err := handler.GetCIDR(masterObj)
		checkErr("GetCIDR", cidr1, err)
		cidr2, err := handler.GetCIDR(workerObj)
		checkErr("GetCIDR", cidr2, err)

		cidrs1, err := handler.GetCIDRs(masterObj)
		checkErr("GetCIDRs", cidrs1, err)
		cidrs2, err := handler.GetCIDRs(workerObj)
		checkErr("GetCIDRs", cidrs2, err)

		masterInfo, err := handler.GetNodeInfo(masterObj)
		checkErr("GetNodeInfo", masterInfo, err)
		workerInfo, err := handler.GetNodeInfo(workerObj)
		checkErr("GetNodeInfo", workerInfo, err)

		age1, err := handler.GetAge(masterObj)
		checkErr("GetAge", age1, err)
		age2, err := handler.GetAge(workerObj)
		checkErr("GetAge", age2, err)

		end := time.Now()
		log.Println("===== Get Node By Object Cost Time:", end.Sub(begin))
		log.Println()
	}

	getByName()
	getByObj()

	// Output:

	//2022/08/24 13:32:51 ===== Get Node By Name
	//2022/08/24 13:32:51 true
	//2022/08/24 13:32:51 false
	//2022/08/24 13:32:51 true
	//2022/08/24 13:32:51 false
	//2022/08/24 13:32:51 GetHostname success: [d11-k8s-master1]
	//2022/08/24 13:32:51 GetHostname success: [d11-k8s-worker1]
	//2022/08/24 13:32:51 GetInternalIP success: [10.250.16.11]
	//2022/08/24 13:32:51 GetInternalIP success: [10.250.16.21]
	//2022/08/24 13:32:51 GetRoles success: [control-plane master]
	//2022/08/24 13:32:52 GetRoles success: []
	//2022/08/24 13:32:52 GetPods success: [nginx calico-node-smlk6 metrics-server-6c4b586fbb-zxmkf]
	//2022/08/24 13:32:52 GetPods success: [kustomize-guestbook-ui-65f8795578-pmnhj ingress-controller-ingress-nginx-controller-jb2w9 calico-kube-controllers-756b87c6c4-f62gt calico-node-sw8hn calico-typha-6dd58fffb8-9hqws coredns-6bf8d6b86b-z7d42 metrics-server-6c4b586fbb-rwbpq local-path-provisioner-7fdb4745c6-r5nbt nfs-provisioner-nfs-subdir-external-provisioner-6668cc6487d6wjt]
	//2022/08/24 13:32:53 GetCIDR success: 192.169.3.0/24
	//2022/08/24 13:32:53 GetCIDR success: 192.169.0.0/24
	//2022/08/24 13:32:53 GetCIDRs success: [192.169.3.0/24]
	//2022/08/24 13:32:53 GetCIDRs success: [192.169.0.0/24]
	//2022/08/24 13:32:54 GetNodeInfo success: &{d11-k8s-master1 [10.250.16.11] 2 3915380Ki 59217979294 2 4017780Ki 64255620Ki amd64 9788d8c4-42de-45e5-b23f-9ae944a23bdd containerd://1.6.7 5.10.0-17-amd64 v1.24.3 v1.24.3 82e56fd2e3a04ed693fbd3412d3b1c60 linux Debian GNU/Linux 11 (bullseye) 054a110d-4cef-3f4d-9c74-2f7ebf7faf59}
	//2022/08/24 13:32:54 GetNodeInfo success: &{d11-k8s-worker1 [10.250.16.21] 2 3915380Ki 59217979294 2 4017780Ki 64255620Ki amd64 bbd3aafa-70df-4fd3-af7a-c87bb23f5ee5 containerd://1.6.7 5.10.0-17-amd64 v1.24.3 v1.24.3 82e56fd2e3a04ed693fbd3412d3b1c60 linux Debian GNU/Linux 11 (bullseye) 054a110d-4cef-3f4d-9c74-2f7ebf7faf59}
	//2022/08/24 13:32:54 GetAge success: 51h35m37.62745s
	//2022/08/24 13:32:54 GetAge success: 51h35m50.827587s
	//2022/08/24 13:32:54 ===== Get Node By Name Cost Time: 3.195244474s
	//2022/08/24 13:32:54
	//2022/08/24 13:32:54 ===== Get Node By Object
	//2022/08/24 13:32:54 true
	//2022/08/24 13:32:54 false
	//2022/08/24 13:32:54 true
	//2022/08/24 13:32:54 false
	//2022/08/24 13:32:54 GetHostname success: [d11-k8s-master1]
	//2022/08/24 13:32:54 GetHostname success: [d11-k8s-worker1]
	//2022/08/24 13:32:54 GetInternalIP success: [10.250.16.11]
	//2022/08/24 13:32:54 GetInternalIP success: [10.250.16.21]
	//2022/08/24 13:32:54 GetRoles success: [control-plane master]
	//2022/08/24 13:32:54 GetRoles success: []
	//2022/08/24 13:32:55 GetPods success: [nginx calico-node-smlk6 metrics-server-6c4b586fbb-zxmkf]
	//2022/08/24 13:32:55 GetPods success: [kustomize-guestbook-ui-65f8795578-pmnhj ingress-controller-ingress-nginx-controller-jb2w9 calico-kube-controllers-756b87c6c4-f62gt calico-node-sw8hn calico-typha-6dd58fffb8-9hqws coredns-6bf8d6b86b-z7d42 metrics-server-6c4b586fbb-rwbpq local-path-provisioner-7fdb4745c6-r5nbt nfs-provisioner-nfs-subdir-external-provisioner-6668cc6487d6wjt]
	//2022/08/24 13:32:55 GetCIDR success: 192.169.3.0/24
	//2022/08/24 13:32:55 GetCIDR success: 192.169.0.0/24
	//2022/08/24 13:32:55 GetCIDRs success: [192.169.3.0/24]
	//2022/08/24 13:32:55 GetCIDRs success: [192.169.0.0/24]
	//2022/08/24 13:32:55 GetNodeInfo success: &{d11-k8s-master1 [10.250.16.11] 2 3915380Ki 59217979294 2 4017780Ki 64255620Ki amd64 9788d8c4-42de-45e5-b23f-9ae944a23bdd containerd://1.6.7 5.10.0-17-amd64 v1.24.3 v1.24.3 82e56fd2e3a04ed693fbd3412d3b1c60 linux Debian GNU/Linux 11 (bullseye) 054a110d-4cef-3f4d-9c74-2f7ebf7faf59}
	//2022/08/24 13:32:55 GetNodeInfo success: &{d11-k8s-worker1 [10.250.16.21] 2 3915380Ki 59217979294 2 4017780Ki 64255620Ki amd64 bbd3aafa-70df-4fd3-af7a-c87bb23f5ee5 containerd://1.6.7 5.10.0-17-amd64 v1.24.3 v1.24.3 82e56fd2e3a04ed693fbd3412d3b1c60 linux Debian GNU/Linux 11 (bullseye) 054a110d-4cef-3f4d-9c74-2f7ebf7faf59}
	//2022/08/24 13:32:55 GetAge success: 51h35m38.623516s
	//2022/08/24 13:32:55 GetAge success: 51h35m51.62352s
	//2022/08/24 13:32:55 ===== Get Node By Object Cost Time: 795.900374ms
	//2022/08/24 13:32:55
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
