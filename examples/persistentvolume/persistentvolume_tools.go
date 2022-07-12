package main

import (
	"log"
	"time"

	"github.com/forbearing/k8s/persistentvolume"
)

func PersistentVolume_Tools() {
	handler, err := persistentvolume.New(ctx, kubeconfig)
	if err != nil {
		panic(err)
	}
	handler.Apply(filename)

	pv, err := handler.Get(name)
	if err != nil {
		panic(err)
	}

	getByName := func() {
		log.Println("===== Get PersistentVolume By Name")
		begin := time.Now()

		cap, err := handler.GetCapacity(name)
		checkErr("GetCapacity", cap, err)

		accessModes, err := handler.GetAccessModes(name)
		checkErr("GetAccessModes", accessModes, err)

		reclaimPolicy, err := handler.GetReclaimPolicy(name)
		checkErr("GetReclaimPolicy", reclaimPolicy, err)

		status, err := handler.GetStatus(name)
		checkErr("GetStatus", status, err)

		pvc, err := handler.GetPVC(name)
		checkErr("GetPVC", pvc, err)

		sc, err := handler.GetStorageClass(name)
		checkErr("GetStorageClass", sc, err)

		vs, err := handler.GetVolumeSource(name)
		checkErr("GetVolumeSource", vs, err)

		vm, err := handler.GetVolumeMode(name)
		checkErr("GetVolumeMode", vm, err)

		age, err := handler.GetAge(name)
		checkErr("GetAge", age, err)

		end := time.Now()
		log.Println("===== Get PersistentVolume By Name Cost Time:", end.Sub(begin))
		log.Println()
	}

	getByObj := func() {
		log.Println("===== Get PersistentVolume By Object")
		begin := time.Now()

		cap, err := handler.GetCapacity(pv)
		checkErr("GetCapacity", cap, err)

		accessModes, err := handler.GetAccessModes(pv)
		checkErr("GetAccessModes", accessModes, err)

		reclaimPolicy, err := handler.GetReclaimPolicy(pv)
		checkErr("GetReclaimPolicy", reclaimPolicy, err)

		status, err := handler.GetStatus(pv)
		checkErr("GetStatus", status, err)

		pvc, err := handler.GetPVC(pv)
		checkErr("GetPVC", pvc, err)

		sc, err := handler.GetStorageClass(pv)
		checkErr("GetStorageClass", sc, err)

		vs, err := handler.GetVolumeSource(pv)
		checkErr("GetVolumeSource", vs, err)

		vm, err := handler.GetVolumeMode(pv)
		checkErr("GetVolumeMode", vm, err)

		age, err := handler.GetAge(pv)
		checkErr("GetAge", age, err)

		end := time.Now()
		log.Println("===== Get PersistentVolume By Object Cost Time:", end.Sub(begin))
	}

	getByName()
	getByObj()

	// Output:

	//2022/07/10 19:27:14 ===== Get PersistentVolume By Name
	//2022/07/10 19:27:14 GetCapacity success: 8000000000
	//2022/07/10 19:27:14 GetAccessModes success: [ReadWriteOnce ReadWriteMany ReadOnlyMany]
	//2022/07/10 19:27:14 GetReclaimPolicy success: Delete
	//2022/07/10 19:27:14 GetStatus success: Available
	//2022/07/10 19:27:14 GetPVC success:
	//2022/07/10 19:27:14 GetStorageClass success: nfs-sc
	//2022/07/10 19:27:14 GetVolumeSource success: nfs
	//2022/07/10 19:27:14 GetVolumeMode success: Filesystem
	//2022/07/10 19:27:14 GetAge success: 23h24m36.466542s
	//2022/07/10 19:27:14 ===== Get PersistentVolume By Name Cost Time: 389.794647ms
	//2022/07/10 19:27:14
	//2022/07/10 19:27:14 ===== Get PersistentVolume By Object
	//2022/07/10 19:27:14 GetCapacity success: 8000000000
	//2022/07/10 19:27:14 GetAccessModes success: [ReadWriteOnce ReadWriteMany ReadOnlyMany]
	//2022/07/10 19:27:14 GetReclaimPolicy success: Delete
	//2022/07/10 19:27:14 GetStatus success: Available
	//2022/07/10 19:27:14 GetPVC success:
	//2022/07/10 19:27:14 GetStorageClass success: nfs-sc
	//2022/07/10 19:27:14 GetVolumeSource success: nfs
	//2022/07/10 19:27:14 GetVolumeMode success: Filesystem
	//2022/07/10 19:27:14 GetAge success: 23h24m36.466592s
	//2022/07/10 19:27:14 ===== Get PersistentVolume By Object Cost Time: 25.003µs
}
