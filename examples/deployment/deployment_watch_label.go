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

	// Output:
	//2022/09/05 10:44:00 added deployment: test/mydep.
	//2022/09/05 10:44:00 modified deployment: test/mydep.
	//2022/09/05 10:44:00 modified deployment: test/mydep2.
	//2022/09/05 10:44:00 modified deployment: test/mydep.
	//2022/09/05 10:44:00 modified deployment: test/mydep2.
	//2022/09/05 10:44:00 modified deployment: test/mydep2.
	//2022/09/05 10:44:00 modified deployment: test/mydep.
	//2022/09/05 10:44:08 modified deployment: test/mydep.
	//2022/09/05 10:44:13 modified deployment: test/mydep2.
	//2022/09/05 10:44:15 deleted deployment: test/mydep.
	//2022/09/05 10:44:15 deleted deployment: test/mydep2.
	//2022/09/05 10:44:20 added deployment: test/mydep.
	//2022/09/05 10:44:20 modified deployment: test/mydep2.
	//2022/09/05 10:44:20 modified deployment: test/mydep.
	//2022/09/05 10:44:20 modified deployment: test/mydep.
	//2022/09/05 10:44:20 modified deployment: test/mydep2.
	//2022/09/05 10:44:20 modified deployment: test/mydep2.
	//2022/09/05 10:44:20 modified deployment: test/mydep.
	//2022/09/05 10:44:26 modified deployment: test/mydep.
	//2022/09/05 10:44:29 modified deployment: test/mydep.
	//2022/09/05 10:44:35 deleted deployment: test/mydep.
	//2022/09/05 10:44:35 deleted deployment: test/mydep2.
	//2022/09/05 10:44:40 added deployment: test/mydep.
	//2022/09/05 10:44:40 modified deployment: test/mydep2.
	//2022/09/05 10:44:40 modified deployment: test/mydep.
	//2022/09/05 10:44:40 modified deployment: test/mydep2.
	//2022/09/05 10:44:40 modified deployment: test/mydep.
	//2022/09/05 10:44:40 modified deployment: test/mydep2.
	//2022/09/05 10:44:41 modified deployment: test/mydep.
	//2022/09/05 10:44:46 modified deployment: test/mydep2.
	//2022/09/05 10:44:49 modified deployment: test/mydep.
	//2022/09/05 10:44:53 modified deployment: test/mydep.
}
