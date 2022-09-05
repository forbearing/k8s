package main

import (
	"context"
	"log"
	"time"

	"github.com/forbearing/k8s/deployment"
	appsv1 "k8s.io/api/apps/v1"
)

func Deployment_Watch_Field() {
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
	field := "metadata.namespace=test"
	name := "mydep"
	name2 := "mydep2"

	handler := deployment.NewOrDie(ctx, "", namespace)
	ctx, cancel := context.WithCancel(ctx)

	go func(ctx context.Context) {
		time.Sleep(time.Second * 5)
		handler.WatchByField(field, addFunc, modifyFunc, deleteFunc)
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

	//2022/09/05 17:11:14 added deployment: "test/mydep".
	//2022/09/05 17:11:14 added deployment: "test/mydep2".
	//2022/09/05 17:11:15 modified deployment: "test/mydep".
	//2022/09/05 17:11:19 modified deployment: "test/mydep".
	//2022/09/05 17:11:22 modified deployment: "test/mydep2".
	//2022/09/05 17:11:25 modified deployment: "test/mydep2".
	//2022/09/05 17:11:29 deleted deployment: "test/mydep".
	//2022/09/05 17:11:29 deleted deployment: "test/mydep2".
	//2022/09/05 17:11:29 added deployment: "test/mydep".
	//2022/09/05 17:11:29 added deployment: "test/mydep2".
	//2022/09/05 17:11:29 modified deployment: "test/mydep".
	//2022/09/05 17:11:29 modified deployment: "test/mydep2".
	//2022/09/05 17:11:29 modified deployment: "test/mydep".
	//2022/09/05 17:11:29 modified deployment: "test/mydep2".
	//2022/09/05 17:11:29 modified deployment: "test/mydep".
	//2022/09/05 17:11:29 modified deployment: "test/mydep2".
	//2022/09/05 17:11:29 modified deployment: "test/mydep".
	//2022/09/05 17:11:29 modified deployment: "test/mydep2".
	//2022/09/05 17:11:37 modified deployment: "test/mydep".
	//2022/09/05 17:11:37 modified deployment: "test/mydep".
	//2022/09/05 17:11:40 modified deployment: "test/mydep".
	//2022/09/05 17:11:43 modified deployment: "test/mydep2".
	//2022/09/05 17:11:47 modified deployment: "test/mydep2".
	//2022/09/05 17:11:49 deleted deployment: "test/mydep".
	//2022/09/05 17:11:49 deleted deployment: "test/mydep2".
	//2022/09/05 17:11:49 added deployment: "test/mydep".
	//2022/09/05 17:11:49 added deployment: "test/mydep2".
	//2022/09/05 17:11:49 modified deployment: "test/mydep".
	//2022/09/05 17:11:49 modified deployment: "test/mydep2".
	//2022/09/05 17:11:49 modified deployment: "test/mydep".
	//2022/09/05 17:11:49 modified deployment: "test/mydep2".
	//2022/09/05 17:11:49 modified deployment: "test/mydep2".
	//2022/09/05 17:11:49 modified deployment: "test/mydep".
	//2022/09/05 17:11:57 modified deployment: "test/mydep".
	//2022/09/05 17:11:58 modified deployment: "test/mydep".
	//2022/09/05 17:12:01 modified deployment: "test/mydep".
	//2022/09/05 17:12:09 deleted deployment: "test/mydep".
	//2022/09/05 17:12:09 deleted deployment: "test/mydep2".
}
