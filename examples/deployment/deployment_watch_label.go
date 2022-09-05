package main

import (
	"context"
	"log"
	"time"

	"github.com/forbearing/k8s/deployment"
	appsv1 "k8s.io/api/apps/v1"
)

func Deployment_Watch_Label() {
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
	label := "type=deployment"
	name := "mydep"
	name2 := "mydep2"

	handler := deployment.NewOrDie(ctx, "", namespace)
	ctx, cancel := context.WithCancel(ctx)

	go func(ctx context.Context) {
		time.Sleep(time.Second * 5)
		handler.WatchByLabel(label, addFunc, modifyFunc, deleteFunc)
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
	//2022/09/05 11:21:13 added deployment: test/mydep.
	//2022/09/05 11:21:13 added deployment: test/mydep2.
	//2022/09/05 11:21:14 modified deployment: test/mydep.
	//2022/09/05 11:21:17 modified deployment: test/mydep.
	//2022/09/05 11:21:20 modified deployment: test/mydep2.
	//2022/09/05 11:21:23 modified deployment: test/mydep.
	//2022/09/05 11:21:28 deleted deployment: test/mydep.
	//2022/09/05 11:21:28 deleted deployment: test/mydep2.
	//2022/09/05 11:21:28 added deployment: test/mydep.
	//2022/09/05 11:21:28 added deployment: test/mydep2.
	//2022/09/05 11:21:28 modified deployment: test/mydep.
	//2022/09/05 11:21:28 modified deployment: test/mydep2.
	//2022/09/05 11:21:28 modified deployment: test/mydep2.
	//2022/09/05 11:21:28 modified deployment: test/mydep.
	//2022/09/05 11:21:28 modified deployment: test/mydep2.
	//2022/09/05 11:21:28 modified deployment: test/mydep.
	//2022/09/05 11:21:29 modified deployment: test/mydep2.
	//2022/09/05 11:21:29 modified deployment: test/mydep.
	//2022/09/05 11:21:36 modified deployment: test/mydep.
	//2022/09/05 11:21:38 modified deployment: test/mydep2.
	//2022/09/05 11:21:41 modified deployment: test/mydep2.
	//2022/09/05 11:21:48 deleted deployment: test/mydep.
	//2022/09/05 11:21:48 deleted deployment: test/mydep2.
	//2022/09/05 11:21:48 added deployment: test/mydep.
	//2022/09/05 11:21:48 added deployment: test/mydep2.
	//2022/09/05 11:21:48 modified deployment: test/mydep.
	//2022/09/05 11:21:48 modified deployment: test/mydep2.
	//2022/09/05 11:21:48 modified deployment: test/mydep.
	//2022/09/05 11:21:48 modified deployment: test/mydep2.
	//2022/09/05 11:21:49 modified deployment: test/mydep.
	//2022/09/05 11:21:49 modified deployment: test/mydep2.
	//2022/09/05 11:21:56 modified deployment: test/mydep2.
	//2022/09/05 11:21:57 modified deployment: test/mydep.
	//2022/09/05 11:22:00 modified deployment: test/mydep2.
	//2022/09/05 11:22:05 modified deployment: test/mydep.
	//2022/09/05 11:22:08 deleted deployment: test/mydep.
}
