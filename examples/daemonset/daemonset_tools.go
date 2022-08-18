package main

import (
	"log"
	"time"

	"github.com/forbearing/k8s"
	"github.com/forbearing/k8s/daemonset"
	corev1 "k8s.io/api/core/v1"
)

func DaemonSet_Tools() {
	handler, err := daemonset.New(ctx, kubeconfig, namespace)
	if err != nil {
		panic(err)
	}
	defer cleanup(handler)

	handler.Apply(filename)
	k8s.ApplyF(ctx, kubeconfig, filename2)
	log.Println(handler.IsReady(name))  // true
	log.Println(handler.IsReady(name2)) // true
	handler.WaitReady(name)
	handler.WaitReady(name2)
	log.Println(handler.IsReady(name))  // true
	log.Println(handler.IsReady(name2)) // true

	ds, err := handler.Get(name)
	if err != nil {
		panic(err)
	}
	ds2, err := handler.Get(name2)
	if err != nil {
		panic(err)
	}

	pods1, err := handler.GetPods(name)
	checkErr("GetPods", printPods(pods1), err)
	pods2, err := handler.GetPods(ds)
	checkErr("GetPods", printPods(pods2), err)

	pvList1, err := handler.GetPV(name2)
	checkErr("GetPV", pvList1, err)
	pvList2, err := handler.GetPV(ds2)
	checkErr("GetPV", pvList2, err)
	log.Println()

	getByName := func() {
		log.Println("===== Get DaemonSet By Name")
		begin := time.Now()

		pvcList, err := handler.GetPVC(name)
		checkErr("GetPVC", pvcList, err)
		numDesired, err := handler.NumDesired(name)
		checkErr("NumDesired", numDesired, err)
		numCurrent, err := handler.NumCurrent(name)
		checkErr("NumCurrent", numCurrent, err)
		numReady, err := handler.NumReady(name)
		checkErr("NumReady", numReady, err)
		numAvailable, err := handler.NumAvailable(name)
		checkErr("NumAvailable", numAvailable, err)
		age, err := handler.GetAge(name)
		checkErr("GetAge", age, err)
		containers, err := handler.GetContainers(name)
		checkErr("GetContainers", containers, err)
		images, err := handler.GetImages(name)
		checkErr("GetImages", images, err)

		end := time.Now()
		log.Println("===== Get DaemonSet By Name Cost Time:", end.Sub(begin))
		log.Println()
	}

	getByObj := func() {
		log.Println("===== Get DaemonSet By Object")
		begin := time.Now()

		pvcList, err := handler.GetPVC(ds)
		checkErr("GetPVC", pvcList, err)
		numDesired, err := handler.NumDesired(ds)
		checkErr("NumDesired", numDesired, err)
		numCurrent, err := handler.NumCurrent(ds)
		checkErr("NumCurrent", numCurrent, err)
		numReady, err := handler.NumReady(ds)
		checkErr("NumReady", numReady, err)
		numAvailable, err := handler.NumAvailable(ds)
		checkErr("NumAvailable", numAvailable, err)
		age, err := handler.GetAge(ds)
		checkErr("GetAge", age, err)
		containers, err := handler.GetContainers(ds)
		checkErr("GetContainers", containers, err)
		images, err := handler.GetImages(ds)
		checkErr("GetImages", images, err)

		end := time.Now()
		log.Println("===== Get DaemonSet By Object Cost Time:", end.Sub(begin))
	}

	getByName()
	getByObj()

	// Output:

	//2022/07/11 17:44:04 false
	//2022/07/11 17:44:04 false
	//2022/07/11 17:44:53 true
	//2022/07/11 17:44:53 true
	//2022/07/11 17:44:53 GetPods success: [myds-7dl8h myds-dbf48 myds-hcl4n myds-k8xq4 myds-mcfch myds-qttkb]
	//2022/07/11 17:44:53 GetPods success: [myds-7dl8h myds-dbf48 myds-hcl4n myds-k8xq4 myds-mcfch myds-qttkb]
	//2022/07/11 17:44:53 GetPV success: [pvc-639c2572-3a65-4a50-99ea-44854f04d911 pvc-d0125122-d8bc-4e57-b5c0-95ee0edde841]
	//2022/07/11 17:44:53 GetPV success: [pvc-639c2572-3a65-4a50-99ea-44854f04d911 pvc-d0125122-d8bc-4e57-b5c0-95ee0edde841]
	//2022/07/11 17:44:53
	//2022/07/11 17:44:53 ===== Get DaemonSet By Name
	//2022/07/11 17:44:53 GetPVC success: []
	//2022/07/11 17:44:53 NumDesired success: 6
	//2022/07/11 17:44:53 NumCurrent success: 6
	//2022/07/11 17:44:53 NumReady success: 6
	//2022/07/11 17:44:53 NumAvailable success: 6
	//2022/07/11 17:44:53 GetAge success: 50.860946s
	//2022/07/11 17:44:54 GetContainers success: [nginx]
	//2022/07/11 17:44:54 GetImages success: [nginx]
	//2022/07/11 17:44:54 ===== Get DaemonSet By Name Cost Time: 866.449419ms
	//2022/07/11 17:44:54
	//2022/07/11 17:44:54 ===== Get DaemonSet By Object
	//2022/07/11 17:44:54 GetPVC success: []
	//2022/07/11 17:44:54 NumDesired success: 6
	//2022/07/11 17:44:54 NumCurrent success: 6
	//2022/07/11 17:44:54 NumReady success: 6
	//2022/07/11 17:44:54 NumAvailable success: 6
	//2022/07/11 17:44:54 GetAge success: 51.260919s
	//2022/07/11 17:44:54 GetContainers success: [nginx]
	//2022/07/11 17:44:54 GetImages success: [nginx]
	//2022/07/11 17:44:54 ===== Get DaemonSet By Object Cost Time: 27.214µs

}

func printPods(podList []*corev1.Pod) []string {
	var pl []string
	for _, pod := range podList {
		pl = append(pl, pod.Name)
	}
	return pl
}
