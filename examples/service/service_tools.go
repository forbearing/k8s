package main

import (
	"log"
	"time"

	"github.com/forbearing/k8s/service"
)

func Service_Tools() {
	handler, err := service.New(ctx, namespace, kubeconfig)
	if err != nil {
		panic(err)
	}
	handler.Apply(filename)
	handler.Apply(filenameNP)
	handler.Apply(filenameLB)
	handler.Apply(filenameEN)
	handler.Apply(filenameEI)

	svc, err := handler.Get(name)
	if err != nil {
		panic(err)
	}
	svcNP, err := handler.Get(nameNP)
	if err != nil {
		panic(err)
	}
	svcLB, err := handler.Get(nameLB)
	if err != nil {
		panic(err)
	}
	svcEN, err := handler.Get(nameEN)
	if err != nil {
		panic(err)
	}
	svcEI, err := handler.Get(nameEI)
	if err != nil {
		panic(err)
	}

	getByName := func() {
		log.Println("===== Get Service By Name")
		begin := time.Now()

		typ, err := handler.GetType(name)
		checkErr("GetType", typ, err)
		typNP, err := handler.GetType(nameNP)
		checkErr("GetType", typNP, err)
		typLB, err := handler.GetType(nameLB)
		checkErr("GetType", typLB, err)
		typEN, err := handler.GetType(nameEN)
		checkErr("GetType", typEN, err)
		typEI, err := handler.GetType(nameEI)
		checkErr("GetType", typEI, err)

		clusterip, err := handler.GetClusterIP(name)
		checkErr("GetClusterIP", clusterip, err)
		clusteripNP, err := handler.GetClusterIP(nameNP)
		checkErr("GetClusterIP", clusteripNP, err)
		clusteripLB, err := handler.GetClusterIP(nameLB)
		checkErr("GetClusterIP", clusteripLB, err)
		clusteripEN, err := handler.GetClusterIP(nameEN)
		checkErr("GetClusterIP", clusteripEN, err)
		clusteripEI, err := handler.GetClusterIP(nameEI)
		checkErr("GetClusterIP", clusteripEI, err)

		externalip, err := handler.GetExternalIPs(name)
		checkErr("GetExternalIPs", externalip, err)
		externalipNP, err := handler.GetExternalIPs(nameNP)
		checkErr("GetExternalIPs", externalipNP, err)
		externalipLB, err := handler.GetExternalIPs(nameLB)
		checkErr("GetExternalIPs", externalipLB, err)
		externalipEN, err := handler.GetExternalIPs(nameEN)
		checkErr("GetExternalIPs", externalipEN, err)
		externalipEI, err := handler.GetExternalIPs(nameEI)
		checkErr("GetExternalIPs", externalipEI, err)

		ports, err := handler.GetPorts(name)
		checkErr("GetPorts", ports, err)
		portsNP, err := handler.GetPorts(nameNP)
		checkErr("GetPorts", portsNP, err)
		portsLB, err := handler.GetPorts(nameLB)
		checkErr("GetPorts", portsLB, err)
		portsEN, err := handler.GetPorts(nameEN)
		checkErr("GetPorts", portsEN, err)
		portsEI, err := handler.GetPorts(nameEI)
		checkErr("GetPorts", portsEI, err)

		age, err := handler.GetAge(name)
		checkErr("GetAge", age, err)
		ageNP, err := handler.GetAge(nameNP)
		checkErr("GetAge", ageNP, err)
		ageLB, err := handler.GetAge(nameLB)
		checkErr("GetAge", ageLB, err)
		ageEN, err := handler.GetAge(nameEN)
		checkErr("GetAge", ageEN, err)
		ageEI, err := handler.GetAge(nameEI)
		checkErr("GetAge", ageEI, err)

		end := time.Now()
		log.Println("===== Get Service By Name Cost Time:", end.Sub(begin))
		log.Println()
	}

	getByObj := func() {
		log.Println("===== Get Service By Object")
		begin := time.Now()

		typ, err := handler.GetType(svc)
		checkErr("GetType", typ, err)
		typNP, err := handler.GetType(svcNP)
		checkErr("GetType", typNP, err)
		typLB, err := handler.GetType(svcLB)
		checkErr("GetType", typLB, err)
		typEN, err := handler.GetType(svcEN)
		checkErr("GetType", typEN, err)
		typEI, err := handler.GetType(svcEI)
		checkErr("GetType", typEI, err)

		clusterip, err := handler.GetClusterIP(svc)
		checkErr("GetClusterIP", clusterip, err)
		clusteripNP, err := handler.GetClusterIP(svcNP)
		checkErr("GetClusterIP", clusteripNP, err)
		clusteripLB, err := handler.GetClusterIP(svcLB)
		checkErr("GetClusterIP", clusteripLB, err)
		clusteripEN, err := handler.GetClusterIP(svcEN)
		checkErr("GetClusterIP", clusteripEN, err)
		clusteripEI, err := handler.GetClusterIP(svcEI)
		checkErr("GetClusterIP", clusteripEI, err)

		externalip, err := handler.GetExternalIPs(svc)
		checkErr("GetExternalIPs", externalip, err)
		externalipNP, err := handler.GetExternalIPs(svcNP)
		checkErr("GetExternalIPs", externalipNP, err)
		externalipLB, err := handler.GetExternalIPs(svcLB)
		checkErr("GetExternalIPs", externalipLB, err)
		externalipEN, err := handler.GetExternalIPs(svcEN)
		checkErr("GetExternalIPs", externalipEN, err)
		externalipEI, err := handler.GetExternalIPs(svcEI)
		checkErr("GetExternalIPs", externalipEI, err)

		ports, err := handler.GetPorts(svc)
		checkErr("GetPorts", ports, err)
		portsNP, err := handler.GetPorts(svcNP)
		checkErr("GetPorts", portsNP, err)
		portsLB, err := handler.GetPorts(svcLB)
		checkErr("GetPorts", portsLB, err)
		portsEN, err := handler.GetPorts(svcEN)
		checkErr("GetPorts", portsEN, err)
		portsEI, err := handler.GetPorts(svcEI)
		checkErr("GetPorts", portsEI, err)

		age, err := handler.GetAge(svc)
		checkErr("GetAge", age, err)
		ageNP, err := handler.GetAge(svcNP)
		checkErr("GetAge", ageNP, err)
		ageLB, err := handler.GetAge(svcLB)
		checkErr("GetAge", ageLB, err)
		ageEN, err := handler.GetAge(svcEN)
		checkErr("GetAge", ageEN, err)
		ageEI, err := handler.GetAge(svcEI)
		checkErr("GetAge", ageEI, err)

		end := time.Now()
		log.Println("===== Get Service By Object Cost Time:", end.Sub(begin))
	}

	getByName()
	getByObj()

	// Output:

	//2022/07/11 09:57:12 ===== Get Service By Name
	//2022/07/11 09:57:12 GetType success: ClusterIP
	//2022/07/11 09:57:12 GetType success: NodePort
	//2022/07/11 09:57:12 GetType success: LoadBalancer
	//2022/07/11 09:57:13 GetType success: ExternalName
	//2022/07/11 09:57:13 GetType success: ClusterIP
	//2022/07/11 09:57:13 GetClusterIP success: 172.18.168.97
	//2022/07/11 09:57:13 GetClusterIP success: 172.18.41.139
	//2022/07/11 09:57:13 GetClusterIP success: 172.18.1.147
	//2022/07/11 09:57:14 GetClusterIP success:
	//2022/07/11 09:57:14 GetClusterIP success: 172.18.136.51
	//2022/07/11 09:57:14 GetExternalIPs success: []
	//2022/07/11 09:57:14 GetExternalIPs success: []
	//2022/07/11 09:57:14 GetExternalIPs success: []
	//2022/07/11 09:57:15 GetExternalIPs success: []
	//2022/07/11 09:57:15 GetExternalIPs success: [1.1.1.1 2.2.2.2]
	//2022/07/11 09:57:15 GetPorts success: [{http TCP 80 {0 80 } 0}]
	//2022/07/11 09:57:15 GetPorts success: [{http TCP 80 {0 80 } 31608} {https TCP 443 {0 443 } 32022}]
	//2022/07/11 09:57:15 GetPorts success: [{http TCP 80 {0 80 } 31953} {web TCP 443 {0 443 } 31371}]
	//2022/07/11 09:57:16 GetPorts success: [{http TCP 80 {0 80 } 0} {https TCP 443 {0 443 } 0}]
	//2022/07/11 09:57:16 GetPorts success: [{http TCP 80 {0 80 } 0} {https TCP 443 {0 443 } 0}]
	//2022/07/11 09:57:16 GetAge success: 240h55m29.510636s
	//2022/07/11 09:57:16 GetAge success: 14m53.715186s
	//2022/07/11 09:57:16 GetAge success: 14m53.908694s
	//2022/07/11 09:57:17 GetAge success: 14m54.112511s
	//2022/07/11 09:57:17 GetAge success: 14m54.304982s
	//2022/07/11 09:57:17 ===== Get Service By Name Cost Time: 4.998930366s
	//2022/07/11 09:57:17
	//2022/07/11 09:57:17 ===== Get Service By Object
	//2022/07/11 09:57:17 GetType success: ClusterIP
	//2022/07/11 09:57:17 GetType success: NodePort
	//2022/07/11 09:57:17 GetType success: LoadBalancer
	//2022/07/11 09:57:17 GetType success: ExternalName
	//2022/07/11 09:57:17 GetType success: ClusterIP
	//2022/07/11 09:57:17 GetClusterIP success: 172.18.168.97
	//2022/07/11 09:57:17 GetClusterIP success: 172.18.41.139
	//2022/07/11 09:57:17 GetClusterIP success: 172.18.1.147
	//2022/07/11 09:57:17 GetClusterIP success:
	//2022/07/11 09:57:17 GetClusterIP success: 172.18.136.51
	//2022/07/11 09:57:17 GetExternalIPs success: []
	//2022/07/11 09:57:17 GetExternalIPs success: []
	//2022/07/11 09:57:17 GetExternalIPs success: []
	//2022/07/11 09:57:17 GetExternalIPs success: []
	//2022/07/11 09:57:17 GetExternalIPs success: [1.1.1.1 2.2.2.2]
	//2022/07/11 09:57:17 GetPorts success: [{http TCP 80 {0 80 } 0}]
	//2022/07/11 09:57:17 GetPorts success: [{http TCP 80 {0 80 } 31608} {https TCP 443 {0 443 } 32022}]
	//2022/07/11 09:57:17 GetPorts success: [{http TCP 80 {0 80 } 31953} {web TCP 443 {0 443 } 31371}]
	//2022/07/11 09:57:17 GetPorts success: [{http TCP 80 {0 80 } 0} {https TCP 443 {0 443 } 0}]
	//2022/07/11 09:57:17 GetPorts success: [{http TCP 80 {0 80 } 0} {https TCP 443 {0 443 } 0}]
	//2022/07/11 09:57:17 GetAge success: 240h55m30.305166s
	//2022/07/11 09:57:17 GetAge success: 14m54.305172s
	//2022/07/11 09:57:17 GetAge success: 14m54.305178s
	//2022/07/11 09:57:17 GetAge success: 14m54.305183s
	//2022/07/11 09:57:17 GetAge success: 14m54.305187s
	//2022/07/11 09:57:17 ===== Get Service By Object Cost Time: 155.932µs
}
