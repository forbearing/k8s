package main

import (
	"context"
	"log"
	"time"

	"github.com/forbearing/k8s/deployment"
	appsv1 "k8s.io/api/apps/v1"
)

func Deployment_Watch_Namespace() {
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
		handler.WatchByNamespace(namespace, addFunc, modifyFunc, deleteFunc)
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

	// Output:

	//2022/09/07 21:02:01 added deployment: "test/mydep".
	//2022/09/07 21:02:01 added deployment: "test/mydep2".
	//2022/09/07 21:02:01 deleted deployment: "test/mydep".
	//2022/09/07 21:02:08 deleted deployment: "test/mydep2".
	//2022/09/07 21:02:11 deleted deployment: "test/mydep2".
	//2022/09/07 21:02:14 deleted deployment: "test/mydep".
	//2022/09/07 21:02:16 modified deployment: "test/mydep".
	//2022/09/07 21:02:16 modified deployment: "test/mydep2".
	//2022/09/07 21:02:16 added deployment: "test/mydep".
	//2022/09/07 21:02:16 added deployment: "test/mydep2".
	//2022/09/07 21:02:16 deleted deployment: "test/mydep2".
	//2022/09/07 21:02:16 deleted deployment: "test/mydep".
	//2022/09/07 21:02:16 deleted deployment: "test/mydep2".
	//2022/09/07 21:02:16 deleted deployment: "test/mydep".
	//2022/09/07 21:02:16 deleted deployment: "test/mydep2".
	//2022/09/07 21:02:16 deleted deployment: "test/mydep".
	//2022/09/07 21:02:23 deleted deployment: "test/mydep2".
	//2022/09/07 21:02:25 deleted deployment: "test/mydep".
	//2022/09/07 21:02:28 deleted deployment: "test/mydep2".
	//2022/09/07 21:02:32 deleted deployment: "test/mydep".
	//2022/09/07 21:02:36 modified deployment: "test/mydep".
	//2022/09/07 21:02:36 modified deployment: "test/mydep2".
	//2022/09/07 21:02:36 added deployment: "test/mydep".
	//2022/09/07 21:02:36 added deployment: "test/mydep2".
	//2022/09/07 21:02:36 deleted deployment: "test/mydep2".
	//2022/09/07 21:02:36 deleted deployment: "test/mydep".
	//2022/09/07 21:02:36 deleted deployment: "test/mydep".
	//2022/09/07 21:02:36 deleted deployment: "test/mydep2".
	//2022/09/07 21:02:36 deleted deployment: "test/mydep2".
	//2022/09/07 21:02:36 deleted deployment: "test/mydep".
	//2022/09/07 21:02:36 deleted deployment: "test/mydep2".
	//2022/09/07 21:02:36 deleted deployment: "test/mydep".
	//2022/09/07 21:02:44 deleted deployment: "test/mydep2".
	//2022/09/07 21:02:47 deleted deployment: "test/mydep2".
	//2022/09/07 21:02:50 deleted deployment: "test/mydep".
	//2022/09/07 21:02:53 deleted deployment: "test/mydep".
	//2022/09/07 21:02:56 modified deployment: "test/mydep".
}
