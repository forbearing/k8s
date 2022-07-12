package main

import (
	"log"
	"time"

	"github.com/forbearing/k8s/storageclass"
)

func StorageClass_Tools() {
	handler, err := storageclass.New(ctx, kubeconfig)
	if err != nil {
		panic(err)
	}
	sc, err := handler.Get(name)
	if err != nil {
		panic(err)
	}

	getByName := func() {
		log.Println("===== Get StorageClass By Name")
		begin := time.Now()

		provisioner, err := handler.GetProvisioner(name)
		checkErr("GetProvisioner", provisioner, err)
		reclaimPolicy, err := handler.GetReclaimPolicy(name)
		checkErr("GetReclaimPolicy", reclaimPolicy, err)
		volumeExpansion, err := handler.IsAllowVolumeExpansion(name)
		checkErr("IsAllowVolumeExpansion", volumeExpansion, err)
		age, err := handler.GetAge(name)
		checkErr("GetAge", age, err)

		end := time.Now()
		log.Println("===== Get StorageClass By Name Cost Time:", end.Sub(begin))
		log.Println()
	}

	getByObj := func() {
		log.Println("===== Get StorageClass By Object")
		begin := time.Now()

		provisioner, err := handler.GetProvisioner(sc)
		checkErr("GetProvisioner", provisioner, err)
		reclaimPolicy, err := handler.GetReclaimPolicy(sc)
		checkErr("GetReclaimPolicy", reclaimPolicy, err)
		volumeExpansion, err := handler.IsAllowVolumeExpansion(sc)
		checkErr("IsAllowVolumeExpansion", volumeExpansion, err)
		age, err := handler.GetAge(sc)
		checkErr("GetAge", age, err)

		end := time.Now()
		log.Println("===== Get StorageClass By Object Cost Time:", end.Sub(begin))
	}

	getByName()
	getByObj()

	// Output:

	//2022/07/11 12:37:36 ===== Get StorageClass By Name
	//2022/07/11 12:37:36 GetProvisioner success: example.com/provisioner
	//2022/07/11 12:37:36 GetReclaimPolicy success: Delete
	//2022/07/11 12:37:36 IsAllowVolumeExpansion failed: AllowVolumeExpansion field not set
	//2022/07/11 12:37:36 GetAge success: 1h23m25.580644s
	//2022/07/11 12:37:36 ===== Get StorageClass By Name Cost Time: 39.020432ms
	//2022/07/11 12:37:36
	//2022/07/11 12:37:36 ===== Get StorageClass By Object
	//2022/07/11 12:37:36 GetProvisioner success: example.com/provisioner
	//2022/07/11 12:37:36 GetReclaimPolicy success: Delete
	//2022/07/11 12:37:36 IsAllowVolumeExpansion failed: AllowVolumeExpansion field not set
	//2022/07/11 12:37:36 GetAge success: 1h23m25.580673s
	//2022/07/11 12:37:36 ===== Get StorageClass By Object Cost Time: 8.72µs
}
