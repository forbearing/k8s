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
	field := "metadata.namespace=test"
	name := "mydep"
	name2 := "mydep2"

	handler := deployment.NewOrDie(ctx, "", namespace)
	ctx, cancel := context.WithCancel(ctx)

	go func(ctx context.Context) {
		handler.WatchByField(field, addFunc, modifyFunc, deleteFunc)
	}(ctx)
	go func(ctx context.Context) {
		for {
			time.Sleep(time.Second * 5)
			handler.Apply(filename)
			handler.Apply(filename2)
			time.Sleep(time.Second * 15)
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
	//2022/09/05 10:58:26 added deployment: test/mydep.
	//2022/09/05 10:58:26 modified deployment: test/mydep2.
	//2022/09/05 10:58:26 modified deployment: test/mydep.
	//2022/09/05 10:58:26 modified deployment: test/mydep2.
	//2022/09/05 10:58:26 modified deployment: test/mydep.
	//2022/09/05 10:58:26 modified deployment: test/mydep2.
	//2022/09/05 10:58:26 modified deployment: test/mydep.
	//2022/09/05 10:58:32 modified deployment: test/mydep2.
	//2022/09/05 10:58:35 modified deployment: test/mydep2.
	//2022/09/05 10:58:38 modified deployment: test/mydep.
	//2022/09/05 10:58:41 deleted deployment: test/mydep.
	//2022/09/05 10:58:41 deleted deployment: test/mydep2.
	//2022/09/05 10:58:46 added deployment: test/mydep.
	//2022/09/05 10:58:46 modified deployment: test/mydep.
	//2022/09/05 10:58:46 modified deployment: test/mydep2.
	//2022/09/05 10:58:46 modified deployment: test/mydep.
	//2022/09/05 10:58:46 modified deployment: test/mydep2.
	//2022/09/05 10:58:46 modified deployment: test/mydep.
	//2022/09/05 10:58:46 modified deployment: test/mydep2.
	//2022/09/05 10:58:50 modified deployment: test/mydep2.
	//2022/09/05 10:58:56 modified deployment: test/mydep.
	//2022/09/05 10:59:01 deleted deployment: test/mydep.
	//2022/09/05 10:59:01 deleted deployment: test/mydep2.
	//2022/09/05 10:59:06 added deployment: test/mydep.
	//2022/09/05 10:59:06 modified deployment: test/mydep2.
	//2022/09/05 10:59:06 modified deployment: test/mydep.
	//2022/09/05 10:59:06 modified deployment: test/mydep2.
	//2022/09/05 10:59:06 modified deployment: test/mydep.
	//2022/09/05 10:59:06 modified deployment: test/mydep.
	//2022/09/05 10:59:06 modified deployment: test/mydep2.
	//2022/09/05 10:59:11 modified deployment: test/mydep2.
	//2022/09/05 10:59:14 modified deployment: test/mydep.
	//2022/09/05 10:59:17 modified deployment: test/mydep2.
	//2022/09/05 10:59:20 modified deployment: test/mydep.
	//2022/09/05 10:59:21 deleted deployment: test/mydep.
}
