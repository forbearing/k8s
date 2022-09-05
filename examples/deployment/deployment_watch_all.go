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
			log.Printf(`added deployment: %s/%s.`, deploy.Namespace, deploy.Name)
		}
		modifyFunc = func(obj interface{}) {
			deploy := obj.(*appsv1.Deployment)
			log.Printf(`modified deployment: %s/%s.`, deploy.Namespace, deploy.Name)
		}
		deleteFunc = func(obj interface{}) {
			deploy := obj.(*appsv1.Deployment)
			log.Printf(`deleted deployment: %s/%s.`, deploy.Namespace, deploy.Name)
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
	//2022/09/05 11:28:01 added deployment: kube-system/coredns.
	//2022/09/05 11:28:01 added deployment: local-path-storage/local-path-provisioner.
	//2022/09/05 11:28:01 added deployment: test/mydep.
	//2022/09/05 11:28:01 added deployment: test/mydep2.
	//2022/09/05 11:28:01 added deployment: default/nginx.
	//2022/09/05 11:28:03 deleted deployment: test/mydep2.
	//2022/09/05 11:28:06 deleted deployment: test/mydep2.
	//2022/09/05 11:28:09 deleted deployment: test/mydep.
	//2022/09/05 11:28:12 deleted deployment: test/mydep2.
	//2022/09/05 11:28:16 modified deployment: test/mydep.
	//2022/09/05 11:28:16 modified deployment: test/mydep2.
	//2022/09/05 11:28:16 added deployment: test/mydep.
	//2022/09/05 11:28:16 added deployment: test/mydep2.
	//2022/09/05 11:28:16 deleted deployment: test/mydep.
	//2022/09/05 11:28:16 deleted deployment: test/mydep2.
	//2022/09/05 11:28:16 deleted deployment: test/mydep.
	//2022/09/05 11:28:16 deleted deployment: test/mydep2.
	//2022/09/05 11:28:16 deleted deployment: test/mydep.
	//2022/09/05 11:28:16 deleted deployment: test/mydep.
	//2022/09/05 11:28:16 deleted deployment: test/mydep2.
	//2022/09/05 11:28:16 deleted deployment: test/mydep2.
	//2022/09/05 11:28:25 deleted deployment: test/mydep.
	//2022/09/05 11:28:27 deleted deployment: test/mydep2.
	//2022/09/05 11:28:36 deleted deployment: test/mydep.
	//2022/09/05 11:28:36 modified deployment: test/mydep.
	//2022/09/05 11:28:36 modified deployment: test/mydep2.
	//2022/09/05 11:28:36 added deployment: test/mydep.
	//2022/09/05 11:28:36 added deployment: test/mydep2.
	//2022/09/05 11:28:36 deleted deployment: test/mydep.
	//2022/09/05 11:28:36 deleted deployment: test/mydep2.
	//2022/09/05 11:28:36 deleted deployment: test/mydep.
	//2022/09/05 11:28:36 deleted deployment: test/mydep2.
	//2022/09/05 11:28:36 deleted deployment: test/mydep2.
	//2022/09/05 11:28:36 deleted deployment: test/mydep.
	//2022/09/05 11:28:44 deleted deployment: test/mydep2.
	//2022/09/05 11:28:46 deleted deployment: test/mydep2.
	//2022/09/05 11:28:49 deleted deployment: test/mydep.
	//2022/09/05 11:28:56 modified deployment: test/mydep.
	//2022/09/05 11:28:56 modified deployment: test/mydep2.
}
