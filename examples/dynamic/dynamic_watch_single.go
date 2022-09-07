package main

import (
	"context"
	"log"
	"time"

	"github.com/forbearing/k8s/deployment"
	"github.com/forbearing/k8s/dynamic"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func Dynamic_Watch_Single() {
	filename := "../../testdata/examples/deployment.yaml"
	name := "mydep"
	gvk := deployment.GVK()

	var (
		addFunc = func(obj interface{}) {
			unstructObj := obj.(*unstructured.Unstructured)
			log.Printf(`added %s: "%s/%s".`, gvk.Kind, unstructObj.GetNamespace(), unstructObj.GetName())
		}
		modifyFunc = func(obj interface{}) {
			unstructObj := obj.(*unstructured.Unstructured)
			log.Printf(`modified %s: "%s/%s".`, gvk.Kind, unstructObj.GetNamespace(), unstructObj.GetName())
		}
		deleteFunc = func(obj interface{}) {
			unstructObj := obj.(*unstructured.Unstructured)
			log.Printf(`deleted %s: "%s/%s".`, gvk.Kind, unstructObj.GetNamespace(), unstructObj.GetName())
		}
	)

	handler := dynamic.NewOrDie(context.TODO(), "", namespace)
	ctx, cancel := context.WithCancel(context.TODO())

	go func(ctx context.Context) {
		time.Sleep(time.Second * 5)
		handler.WithGVK(gvk).WatchByName(name, addFunc, modifyFunc, deleteFunc)
	}(ctx)
	go func(ctx context.Context) {
		for {
			handler.Apply(filename)
			time.Sleep(time.Second * 20)
			handler.WithGVK(gvk).Delete(name)
		}
	}(ctx)

	timer := time.NewTimer(time.Second * 60)
	<-timer.C
	cancel()
	handler.WithGVK(gvk).Delete(name)

	// Output

	//2022/09/07 21:41:32 added Deployment: "test/mydep".
	//2022/09/07 21:41:44 modified Deployment: "test/mydep".
	//2022/09/07 21:41:47 modified Deployment: "test/mydep".
	//2022/09/07 21:41:47 deleted Deployment: "test/mydep".
	//2022/09/07 21:41:47 added Deployment: "test/mydep".
	//2022/09/07 21:41:47 modified Deployment: "test/mydep".
	//2022/09/07 21:41:47 modified Deployment: "test/mydep".
	//2022/09/07 21:41:47 modified Deployment: "test/mydep".
	//2022/09/07 21:41:47 modified Deployment: "test/mydep".
	//2022/09/07 21:41:54 modified Deployment: "test/mydep".
	//2022/09/07 21:42:04 modified Deployment: "test/mydep".
	//2022/09/07 21:42:07 deleted Deployment: "test/mydep".
	//2022/09/07 21:42:07 added Deployment: "test/mydep".
	//2022/09/07 21:42:07 modified Deployment: "test/mydep".
	//2022/09/07 21:42:07 modified Deployment: "test/mydep".
	//2022/09/07 21:42:07 modified Deployment: "test/mydep".
	//2022/09/07 21:42:13 modified Deployment: "test/mydep".
	//2022/09/07 21:42:17 modified Deployment: "test/mydep".
}
