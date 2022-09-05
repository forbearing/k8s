package main

import (
	"context"
	"log"
	"time"

	"github.com/forbearing/k8s/deployment"
	appsv1 "k8s.io/api/apps/v1"
)

func Deployment_Watch_All() {
	var (
		addFunc = func(obj interface{}) {
			deploy := obj.(*appsv1.Deployment)
			log.Printf(`added deployment: "%s/%s".`, deploy.Namespace, deploy.Name)
		}
		modifyFunc = func(obj interface{}) {
			deploy := obj.(*appsv1.Deployment)
			log.Printf(`modified deployment: "%s/%s".`, deploy.Namespace, deploy.Name)
		}
		deleteFunc = func(obj interface{}) {
			deploy := obj.(*appsv1.Deployment)
			log.Printf(`deleted deployment: "%s/%s".`, deploy.Namespace, deploy.Name)
		}
	)

	filename := "../../testdata/examples/deployment.yaml"
	filename2 := "../../testdata/examples/deployment-2.yaml"
	name := "mydep"
	name2 := "mydep2"

	handler := deployment.NewOrDie(ctx, "", namespace)
	ctx, cancel := context.WithCancel(ctx)

	go func(ctx context.Context) {
		time.Sleep(time.Second * 5)
		handler.Watch(addFunc, modifyFunc, deleteFunc)
	}(ctx)
	go func(ctx context.Context) {
		for {
			handler.Apply(filename)
			handler.Apply(filename2)
			time.Sleep(time.Second * 20)
			handler.Delete(name)
			handler.Delete(name2)
		}
	}(ctx)

	timer := time.NewTimer(time.Second * 60)
	<-timer.C
	cancel()
	handler.Delete(name)
	handler.Delete(name2)

	// Outputs:

	//2022/09/05 17:13:14 added deployment: "kube-system/coredns".
	//2022/09/05 17:13:14 added deployment: "local-path-storage/local-path-provisioner".
	//2022/09/05 17:13:14 added deployment: "test/mydep".
	//2022/09/05 17:13:14 added deployment: "test/mydep2".
	//2022/09/05 17:13:14 added deployment: "default/nginx".
	//2022/09/05 17:13:15 deleted deployment: "test/mydep2".
	//2022/09/05 17:13:18 deleted deployment: "test/mydep".
	//2022/09/05 17:13:21 deleted deployment: "test/mydep2".
	//2022/09/05 17:13:24 deleted deployment: "test/mydep2".
	//2022/09/05 17:13:28 deleted deployment: "test/mydep".
	//2022/09/05 17:13:29 modified deployment: "test/mydep".
	//2022/09/05 17:13:29 modified deployment: "test/mydep2".
	//2022/09/05 17:13:29 added deployment: "test/mydep".
	//2022/09/05 17:13:29 added deployment: "test/mydep2".
	//2022/09/05 17:13:29 deleted deployment: "test/mydep".
	//2022/09/05 17:13:29 deleted deployment: "test/mydep2".
	//2022/09/05 17:13:29 deleted deployment: "test/mydep".
	//2022/09/05 17:13:29 deleted deployment: "test/mydep2".
	//2022/09/05 17:13:29 deleted deployment: "test/mydep".
	//2022/09/05 17:13:29 deleted deployment: "test/mydep2".
	//2022/09/05 17:13:29 deleted deployment: "test/mydep2".
	//2022/09/05 17:13:29 deleted deployment: "test/mydep".
	//2022/09/05 17:13:37 deleted deployment: "test/mydep2".
	//2022/09/05 17:13:39 deleted deployment: "test/mydep".
	//2022/09/05 17:13:44 deleted deployment: "test/mydep".
	//2022/09/05 17:13:49 modified deployment: "test/mydep".
	//2022/09/05 17:13:49 modified deployment: "test/mydep2".
	//2022/09/05 17:13:49 added deployment: "test/mydep".
	//2022/09/05 17:13:49 added deployment: "test/mydep2".
	//2022/09/05 17:13:49 deleted deployment: "test/mydep".
	//2022/09/05 17:13:49 deleted deployment: "test/mydep2".
	//2022/09/05 17:13:49 deleted deployment: "test/mydep".
	//2022/09/05 17:13:49 deleted deployment: "test/mydep2".
	//2022/09/05 17:13:49 deleted deployment: "test/mydep".
	//2022/09/05 17:13:49 deleted deployment: "test/mydep2".
	//2022/09/05 17:13:57 deleted deployment: "test/mydep".
	//2022/09/05 17:13:58 deleted deployment: "test/mydep2".
	//2022/09/05 17:14:01 deleted deployment: "test/mydep".
	//2022/09/05 17:14:04 deleted deployment: "test/mydep".
	//2022/09/05 17:14:07 deleted deployment: "test/mydep2".
	//2022/09/05 17:14:09 modified deployment: "test/mydep".
	//2022/09/05 17:14:09 modified deployment: "test/mydep2".
}
