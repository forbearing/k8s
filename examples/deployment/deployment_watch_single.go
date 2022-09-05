package main

import (
	"context"
	"log"
	"time"

	"github.com/forbearing/k8s/deployment"
	appsv1 "k8s.io/api/apps/v1"
)

func Deployment_Watch_Single() {
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

	handler := deployment.NewOrDie(ctx, "", namespace)
	ctx, cancel := context.WithCancel(ctx)

	go func(ctx context.Context) {
		time.Sleep(time.Second * 5)
		handler.WatchByName(name, addFunc, modifyFunc, deleteFunc)
	}(ctx)
	go func(ctx context.Context) {
		for {
			handler.Apply(filename)
			time.Sleep(time.Second * 20)
			handler.Delete(name)
		}
	}(ctx)

	timer := time.NewTimer(time.Second * 60)
	<-timer.C
	cancel()
	handler.Delete(name)

	// Output:
	//2022/09/05 11:19:19 added deployment: test/mydep.
	//2022/09/05 11:19:21 modified deployment: test/mydep.
	//2022/09/05 11:19:24 modified deployment: test/mydep.
	//2022/09/05 11:19:34 deleted deployment: test/mydep.
	//2022/09/05 11:19:34 added deployment: test/mydep.
	//2022/09/05 11:19:34 modified deployment: test/mydep.
	//2022/09/05 11:19:34 modified deployment: test/mydep.
	//2022/09/05 11:19:34 modified deployment: test/mydep.
	//2022/09/05 11:19:34 modified deployment: test/mydep.
	//2022/09/05 11:19:39 modified deployment: test/mydep.
	//2022/09/05 11:19:42 modified deployment: test/mydep.
	//2022/09/05 11:19:54 deleted deployment: test/mydep.
	//2022/09/05 11:19:54 added deployment: test/mydep.
	//2022/09/05 11:19:54 modified deployment: test/mydep.
	//2022/09/05 11:19:54 modified deployment: test/mydep.
	//2022/09/05 11:19:54 modified deployment: test/mydep.
	//2022/09/05 11:19:59 modified deployment: test/mydep.
	//2022/09/05 11:20:14 deleted deployment: test/mydep.
}
