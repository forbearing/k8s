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

	//2022/09/05 17:08:45 added deployment: "test/mydep".
	//2022/09/05 17:08:45 added deployment: "test/mydep2".
	//2022/09/05 17:08:46 modified deployment: "test/mydep2".
	//2022/09/05 17:08:50 modified deployment: "test/mydep2".
	//2022/09/05 17:08:53 modified deployment: "test/mydep2".
	//2022/09/05 17:08:57 modified deployment: "test/mydep".
	//2022/09/05 17:09:00 deleted deployment: "test/mydep".
	//2022/09/05 17:09:00 deleted deployment: "test/mydep2".
	//2022/09/05 17:09:00 added deployment: "test/mydep".
	//2022/09/05 17:09:00 added deployment: "test/mydep2".
	//2022/09/05 17:09:00 modified deployment: "test/mydep2".
	//2022/09/05 17:09:00 modified deployment: "test/mydep".
	//2022/09/05 17:09:00 modified deployment: "test/mydep".
	//2022/09/05 17:09:00 modified deployment: "test/mydep2".
	//2022/09/05 17:09:00 modified deployment: "test/mydep".
	//2022/09/05 17:09:00 modified deployment: "test/mydep2".
	//2022/09/05 17:09:00 modified deployment: "test/mydep".
	//2022/09/05 17:09:00 modified deployment: "test/mydep2".
	//2022/09/05 17:09:06 modified deployment: "test/mydep2".
	//2022/09/05 17:09:09 modified deployment: "test/mydep".
	//2022/09/05 17:09:14 modified deployment: "test/mydep".
	//2022/09/05 17:09:20 deleted deployment: "test/mydep".
	//2022/09/05 17:09:20 deleted deployment: "test/mydep2".
	//2022/09/05 17:09:20 added deployment: "test/mydep".
	//2022/09/05 17:09:20 added deployment: "test/mydep2".
	//2022/09/05 17:09:20 modified deployment: "test/mydep".
	//2022/09/05 17:09:20 modified deployment: "test/mydep2".
	//2022/09/05 17:09:20 modified deployment: "test/mydep".
	//2022/09/05 17:09:20 modified deployment: "test/mydep2".
	//2022/09/05 17:09:20 modified deployment: "test/mydep2".
	//2022/09/05 17:09:20 modified deployment: "test/mydep".
	//2022/09/05 17:09:27 modified deployment: "test/mydep2".
	//2022/09/05 17:09:40 deleted deployment: "test/mydep".
}
