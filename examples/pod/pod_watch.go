package main

import (
	"context"
	"log"
	"time"

	"github.com/forbearing/k8s/deployment"
	"github.com/forbearing/k8s/pod"
	corev1 "k8s.io/api/core/v1"
)

func Pod_Watch() {
	var (
		addFunc = func(obj interface{}) {
			pod := obj.(*corev1.Pod)
			log.Printf(`added pod: "%s/%s".`, pod.Namespace, pod.Name)
		}
		modifyFunc = func(obj interface{}) {
			pod := obj.(*corev1.Pod)
			log.Printf(`modified pod: "%s/%s".`, pod.Namespace, pod.Name)
		}
		deleteFunc = func(obj interface{}) {
			pod := obj.(*corev1.Pod)
			log.Printf(`deleted pod: "%s/%s".`, pod.Namespace, pod.Name)
		}
	)

	filename := "../../testdata/examples/deployment.yaml"
	name := "mydep"

	podHandler := pod.NewOrDie(ctx, "", namespace)
	deployHandler := deployment.NewOrDie(ctx, "", namespace)
	ctx, cancel := context.WithCancel(ctx)

	go func(ctx context.Context) {
		time.Sleep(time.Second * 5)
		podHandler.Watch(addFunc, modifyFunc, deleteFunc)
	}(ctx)
	go func(ctx context.Context) {
		for {
			deployHandler.Apply(filename)
			time.Sleep(time.Second * 20)
			deployHandler.Delete(name)
		}
	}(ctx)

	timer := time.NewTimer(time.Second * 60)
	<-timer.C
	cancel()
	deployHandler.Delete(name)

	// Output:

	//2022/09/08 22:47:01 added pod: "kube-system/etcd-operator-control-plane".
	//2022/09/08 22:47:01 added pod: "kube-system/kube-apiserver-operator-control-plane".
	//2022/09/08 22:47:01 added pod: "kube-system/kube-controller-manager-operator-control-plane".
	//2022/09/08 22:47:01 added pod: "kube-system/kube-scheduler-operator-control-plane".
	//2022/09/08 22:47:01 added pod: "test/mydep-859646944d-72nqr".
	//2022/09/08 22:47:01 added pod: "local-path-storage/local-path-provisioner-66b445c94-k8frj".
	//2022/09/08 22:47:01 added pod: "test/mydep-859646944d-ls4mk".
	//2022/09/08 22:47:01 added pod: "test/mydep-859646944d-scr5v".
	//2022/09/08 22:47:01 added pod: "default/nginx-85b98978db-tf4w7".
	//2022/09/08 22:47:01 added pod: "kube-system/coredns-64897985d-467q7".
	//2022/09/08 22:47:01 added pod: "kube-system/coredns-64897985d-tgjph".
	//2022/09/08 22:47:01 added pod: "kube-system/kindnet-pg7xh".
	//2022/09/08 22:47:01 added pod: "kube-system/kube-proxy-p6k42".
	//2022/09/08 22:47:04 deleted pod: "test/mydep-859646944d-ls4mk".
	//2022/09/08 22:47:09 deleted pod: "test/mydep-859646944d-72nqr".
	//2022/09/08 22:47:11 deleted pod: "test/mydep-859646944d-scr5v".
	//2022/09/08 22:47:16 added pod: "test/mydep-67fcc784fc-hfjhx".
	//2022/09/08 22:47:16 added pod: "test/mydep-67fcc784fc-r6b5g".
	//2022/09/08 22:47:16 deleted pod: "test/mydep-67fcc784fc-hfjhx".
	//2022/09/08 22:47:16 added pod: "test/mydep-67fcc784fc-jlccm".
	//2022/09/08 22:47:16 deleted pod: "test/mydep-859646944d-72nqr".
	//2022/09/08 22:47:16 deleted pod: "test/mydep-859646944d-scr5v".
	//2022/09/08 22:47:16 deleted pod: "test/mydep-859646944d-ls4mk".
	//2022/09/08 22:47:16 deleted pod: "test/mydep-67fcc784fc-r6b5g".
	//2022/09/08 22:47:16 deleted pod: "test/mydep-67fcc784fc-jlccm".
	//2022/09/08 22:47:16 deleted pod: "test/mydep-67fcc784fc-hfjhx".
	//2022/09/08 22:47:16 deleted pod: "test/mydep-859646944d-scr5v".
	//2022/09/08 22:47:16 deleted pod: "test/mydep-67fcc784fc-r6b5g".
	//2022/09/08 22:47:16 deleted pod: "test/mydep-67fcc784fc-jlccm".
	//2022/09/08 22:47:17 deleted pod: "test/mydep-859646944d-ls4mk".
	//2022/09/08 22:47:17 deleted pod: "test/mydep-859646944d-ls4mk".
	//2022/09/08 22:47:17 modified pod: "test/mydep-859646944d-ls4mk".
	//2022/09/08 22:47:18 deleted pod: "test/mydep-859646944d-72nqr".
	//2022/09/08 22:47:18 deleted pod: "test/mydep-859646944d-72nqr".
	//2022/09/08 22:47:18 modified pod: "test/mydep-859646944d-72nqr".
	//2022/09/08 22:47:19 deleted pod: "test/mydep-859646944d-scr5v".
	//2022/09/08 22:47:19 modified pod: "test/mydep-859646944d-scr5v".
	//2022/09/08 22:47:22 deleted pod: "test/mydep-67fcc784fc-hfjhx".
	//2022/09/08 22:47:27 deleted pod: "test/mydep-67fcc784fc-jlccm".
	//2022/09/08 22:47:30 deleted pod: "test/mydep-67fcc784fc-r6b5g".
	//2022/09/08 22:47:36 deleted pod: "test/mydep-67fcc784fc-hfjhx".
	//2022/09/08 22:47:36 deleted pod: "test/mydep-67fcc784fc-r6b5g".
	//2022/09/08 22:47:36 added pod: "test/mydep-859646944d-xl7rb".
	//2022/09/08 22:47:36 deleted pod: "test/mydep-67fcc784fc-jlccm".
	//2022/09/08 22:47:36 added pod: "test/mydep-859646944d-glwh6".
	//2022/09/08 22:47:36 deleted pod: "test/mydep-859646944d-xl7rb".
	//2022/09/08 22:47:36 added pod: "test/mydep-859646944d-z55cl".
	//2022/09/08 22:47:36 deleted pod: "test/mydep-859646944d-glwh6".
	//2022/09/08 22:47:36 deleted pod: "test/mydep-859646944d-z55cl".
	//2022/09/08 22:47:36 deleted pod: "test/mydep-859646944d-xl7rb".
	//2022/09/08 22:47:36 deleted pod: "test/mydep-859646944d-z55cl".
	//2022/09/08 22:47:36 deleted pod: "test/mydep-859646944d-glwh6".
	//2022/09/08 22:47:37 deleted pod: "test/mydep-67fcc784fc-jlccm".
	//2022/09/08 22:47:37 deleted pod: "test/mydep-67fcc784fc-jlccm".
	//2022/09/08 22:47:37 modified pod: "test/mydep-67fcc784fc-jlccm".
	//2022/09/08 22:47:38 deleted pod: "test/mydep-67fcc784fc-hfjhx".
	//2022/09/08 22:47:38 deleted pod: "test/mydep-67fcc784fc-hfjhx".
	//2022/09/08 22:47:38 modified pod: "test/mydep-67fcc784fc-hfjhx".
	//2022/09/08 22:47:39 deleted pod: "test/mydep-67fcc784fc-r6b5g".
	//2022/09/08 22:47:39 deleted pod: "test/mydep-67fcc784fc-r6b5g".
	//2022/09/08 22:47:39 modified pod: "test/mydep-67fcc784fc-r6b5g".
	//2022/09/08 22:47:41 deleted pod: "test/mydep-859646944d-xl7rb".
	//2022/09/08 22:47:45 deleted pod: "test/mydep-859646944d-glwh6".
	//2022/09/08 22:47:48 deleted pod: "test/mydep-859646944d-z55cl".
}
