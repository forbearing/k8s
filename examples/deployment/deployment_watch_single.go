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
	//2022/09/05 10:22:21 added deployment: test/mydep.
	//2022/09/05 10:22:21 modified deployment: test/mydep.
	//2022/09/05 10:22:22 modified deployment: test/mydep.
	//2022/09/05 10:22:22 modified deployment: test/mydep.
	//2022/09/05 10:22:27 modified deployment: test/mydep.
	//2022/09/05 10:22:29 modified deployment: test/mydep.
	//2022/09/05 10:22:32 modified deployment: test/mydep.
	//2022/09/05 10:22:36 deleted deployment: test/mydep.
	//2022/09/05 10:22:41 added deployment: test/mydep.
	//2022/09/05 10:22:42 modified deployment: test/mydep.
	//2022/09/05 10:22:42 modified deployment: test/mydep.
	//2022/09/05 10:22:42 modified deployment: test/mydep.
	//2022/09/05 10:22:46 modified deployment: test/mydep.
	//2022/09/05 10:22:49 modified deployment: test/mydep.
	//2022/09/05 10:22:52 modified deployment: test/mydep.
	//2022/09/05 10:22:57 deleted deployment: test/mydep.
	//2022/09/05 10:23:02 added deployment: test/mydep.
	//2022/09/05 10:23:02 modified deployment: test/mydep.
	//2022/09/05 10:23:02 modified deployment: test/mydep.
	//2022/09/05 10:23:02 modified deployment: test/mydep.
	//2022/09/05 10:23:06 modified deployment: test/mydep.
	//2022/09/05 10:23:09 modified deployment: test/mydep.
	//2022/09/05 10:23:12 modified deployment: test/mydep.
}
