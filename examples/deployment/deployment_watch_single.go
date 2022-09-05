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

	//2022/09/05 17:05:44 added deployment: "test/mydep".
	//2022/09/05 17:05:45 modified deployment: "test/mydep".
	//2022/09/05 17:05:59 deleted deployment: "test/mydep".
	//2022/09/05 17:05:59 added deployment: "test/mydep".
	//2022/09/05 17:05:59 modified deployment: "test/mydep".
	//2022/09/05 17:05:59 modified deployment: "test/mydep".
	//2022/09/05 17:05:59 modified deployment: "test/mydep".
	//2022/09/05 17:05:59 modified deployment: "test/mydep".
	//2022/09/05 17:06:17 modified deployment: "test/mydep".
	//2022/09/05 17:06:19 deleted deployment: "test/mydep".
	//2022/09/05 17:06:19 added deployment: "test/mydep".
	//2022/09/05 17:06:19 modified deployment: "test/mydep".
	//2022/09/05 17:06:19 modified deployment: "test/mydep".
	//2022/09/05 17:06:19 modified deployment: "test/mydep".
	//2022/09/05 17:06:25 modified deployment: "test/mydep".
	//2022/09/05 17:06:28 modified deployment: "test/mydep".
	//2022/09/05 17:06:39 deleted deployment: "test/mydep".
}
