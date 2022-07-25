package main

import (
	"log"
	"time"

	"github.com/forbearing/k8s/persistentvolumeclaim"
)

func PersistentVolumeClaim_Tools() {
	//handler, err := persisetntvolumeclaim.New()
	handler, err := persistentvolumeclaim.New(ctx, kubeconfig, namespace)
	if err != nil {
		panic(err)
	}
	pvc, err := handler.Get(name)
	if err != nil {
		panic(err)
	}

	getByName := func() {
		log.Println("===== Get PersistentVolumeClaim By Name")
		begin := time.Now()

		status, err := handler.GetStatus(name)
		checkErr("GetStatus", status, err)
		pv, err := handler.GetVolume(name)
		checkErr("GetVolume", pv, err)
		cap, err := handler.GetCapacity(name)
		checkErr("GetCapacity", cap, err)
		accessModes, err := handler.GetAccessModes(name)
		checkErr("GetAccessModes", accessModes, err)
		sc, err := handler.GetStorageClass(name)
		checkErr("GetStorageClass", sc, err)
		age, err := handler.GetAge(name)
		checkErr("GetAge", age, err)
		volumeMode, err := handler.GetVolumeMode(name)
		checkErr("GetVolumeMode", volumeMode, err)

		end := time.Now()
		log.Println("===== Get PersistentVolumeClaim By Name Cost Time:", end.Sub(begin))
		log.Println()
	}

	getByObj := func() {
		log.Println("===== Get PersistentVolumeClaim By Object")
		begin := time.Now()

		status, err := handler.GetStatus(pvc)
		checkErr("GetStatus", status, err)
		pv, err := handler.GetVolume(pvc)
		checkErr("GetVolume", pv, err)
		cap, err := handler.GetCapacity(pvc)
		checkErr("GetCapacity", cap, err)
		accessModes, err := handler.GetAccessModes(pvc)
		checkErr("GetAccessModes", accessModes, err)
		sc, err := handler.GetStorageClass(pvc)
		checkErr("GetStorageClass", sc, err)
		age, err := handler.GetAge(pvc)
		checkErr("GetAge", age, err)
		volumeMode, err := handler.GetVolumeMode(pvc)
		checkErr("GetVolumeMode", volumeMode, err)

		end := time.Now()
		log.Println("===== Get PersistentVolumeClaim By Object Cost Time:", end.Sub(begin))
	}

	getByName()
	getByObj()

	// Output:

	//2022/07/10 20:40:08 ===== Get PersistentVolumeClaim By Name
	//2022/07/10 20:40:08 GetStatus success: Bound
	//2022/07/10 20:40:08 GetVolume success: mypv
	//2022/07/10 20:40:08 GetCapacity success: 8000000000
	//2022/07/10 20:40:08 GetAccessModes success: [ReadWriteOnce ReadWriteMany ReadOnlyMany]
	//2022/07/10 20:40:08 GetStorageClass success: nfs-sc
	//2022/07/10 20:40:08 GetAge success: 37m11.057548s
	//2022/07/10 20:40:08 GetVolumeMode success: Filesystem
	//2022/07/10 20:40:08 ===== Get PersistentVolumeClaim By Name Cost Time: 19.274162ms
	//2022/07/10 20:40:08
	//2022/07/10 20:40:08 ===== Get PersistentVolumeClaim By Object
	//2022/07/10 20:40:08 GetStatus success: Bound
	//2022/07/10 20:40:08 GetVolume success: mypv
	//2022/07/10 20:40:08 GetCapacity success: 8000000000
	//2022/07/10 20:40:08 GetAccessModes success: [ReadWriteOnce ReadWriteMany ReadOnlyMany]
	//2022/07/10 20:40:08 GetStorageClass success: nfs-sc
	//2022/07/10 20:40:08 GetAge success: 37m11.059772s
	//2022/07/10 20:40:08 GetVolumeMode success: Filesystem
	//2022/07/10 20:40:08 ===== Get PersistentVolumeClaim By Object Cost Time: 15.368µs
}
