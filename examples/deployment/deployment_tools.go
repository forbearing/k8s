package main

import (
	"log"
	"time"

	"github.com/forbearing/k8s"
	"github.com/forbearing/k8s/deployment"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

func Deployment_Tools() {
	handler, err := deployment.New(ctx, kubeconfig, namespace)
	if err != nil {
		panic(err)
	}
	defer cleanup(handler)

	// kubectl apply -f ../../testdata/nginx/nginx-deploy.yaml
	k8s.ApplyF(ctx, kubeconfig, filename2)

	log.Println(handler.IsReady(name2)) // false
	handler.WaitReady(name2)            // block until the deployment is ready and available.
	log.Println(handler.IsReady(name2)) // true

	deploy, err := handler.Get(name2)
	if err != nil {
		panic(err)
	}

	getByName := func() {
		log.Println("===== Get Deployment By Name")
		begin := time.Now()

		// GetRS get all replicaset that created by the deployment.
		rsList, err := handler.GetRS(name2)
		checkErr("GetRS", printRS(rsList), err)
		// GetPods get all pods that created by the deployment.
		podList, err := handler.GetPods(name2)
		checkErr("GetRS", printPods(podList), err)
		// GetPV get all persistentvolume that attached by the deployment.
		pvList, err := handler.GetPV(name2)
		checkErr("GetPV", pvList, err)
		// GetPVC get all persistentvolumeclaim that attached by the deployment.
		pvcList, err := handler.GetPVC(name2)
		checkErr("GetPVC", pvcList, err)

		end := time.Now()
		log.Println("===== Get Deployment By Name Cost Time:", end.Sub(begin))
		log.Println()
	}

	getByObj := func() {
		log.Println("===== Get Deployment By Object")
		begin := time.Now()

		// GetRS get all replicaset that created by the deployment.
		rsList, err := handler.GetRS(deploy)
		checkErr("GetRS", printRS(rsList), err)
		// GetPods get all pods that created by the deployment.
		podList, err := handler.GetPods(deploy)
		checkErr("GetRS", printPods(podList), err)
		// GetPV get all persistentvolume that attached by the deployment.
		pvList, err := handler.GetPV(deploy)
		checkErr("GetPV", pvList, err)
		// GetPVC get all persistentvolumeclaim that attached by the deployment.
		pvcList, err := handler.GetPVC(deploy)
		checkErr("GetPVC", pvcList, err)

		end := time.Now()
		log.Println("===== Get Deployment By Object Cost Time:", end.Sub(begin))
	}

	getByName()
	getByObj()

	// Output:

	//2022/07/12 09:20:45 false
	//2022/07/12 09:22:16 true
	//2022/07/12 09:22:16 ===== Get Deployment By Name
	//2022/07/12 09:22:16 GetRS success: [nginx-deploy-79979d95dd]
	//2022/07/12 09:22:16 GetRS success: [nginx-deploy-79979d95dd-4lm8h nginx-deploy-79979d95dd-5l9rk nginx-deploy-79979d95dd-scjw9]
	//2022/07/12 09:22:16 GetPV success: [pvc-93ebe9a0-c464-439b-a252-51afb4d87069 pvc-c048ccf9-4d0c-4312-bb36-15e4fa7a1746 pvc-dc16fea0-f511-42d7-b78e-6fcac96fcc9b]
	//2022/07/12 09:22:16 GetPVC success: [deploy-k8s-tools-data deploy-nginx-data deploy-nginx-html]
	//2022/07/12 09:22:16 ===== Get Deployment By Name Cost Time: 82.467272ms
	//2022/07/12 09:22:16
	//2022/07/12 09:22:16 ===== Get Deployment By Object
	//2022/07/12 09:22:16 GetRS success: [nginx-deploy-79979d95dd]
	//2022/07/12 09:22:16 GetRS success: [nginx-deploy-79979d95dd-4lm8h nginx-deploy-79979d95dd-5l9rk nginx-deploy-79979d95dd-scjw9]
	//2022/07/12 09:22:17 GetPV success: [pvc-93ebe9a0-c464-439b-a252-51afb4d87069 pvc-c048ccf9-4d0c-4312-bb36-15e4fa7a1746 pvc-dc16fea0-f511-42d7-b78e-6fcac96fcc9b]
	//2022/07/12 09:22:17 GetPVC success: [deploy-k8s-tools-data deploy-nginx-data deploy-nginx-html]
	//2022/07/12 09:22:17 ===== Get Deployment By Object Cost Time: 134.639944ms
}

func printPods(podList []corev1.Pod) []string {
	var pl []string
	for _, pod := range podList {
		pl = append(pl, pod.Name)
	}
	return pl
}
func printRS(rsList []appsv1.ReplicaSet) []string {
	var rl []string
	for _, rs := range rsList {
		rl = append(rl, rs.Name)
	}
	return rl
}
