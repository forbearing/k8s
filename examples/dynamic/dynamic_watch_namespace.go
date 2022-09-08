package main

import (
	"context"
	"log"
	"time"

	"github.com/forbearing/k8s/deployment"
	"github.com/forbearing/k8s/dynamic"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func Dynamic_Watch_Namespace() {
	filename1 := "../../testdata/examples/deployment.yaml"
	filename2 := "../../testdata/examples/deployment-2.yaml"
	filename3 := "../../testdata/examples/deployment-nolabel.yaml"
	name1 := "mydep"
	name2 := "mydep2"
	name3 := "mydep-nolabel"
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
		handler.WithGVK(gvk).WatchByNamespace(namespace, addFunc, modifyFunc, deleteFunc)
	}(ctx)
	go func(ctx context.Context) {
		for {
			handler.Apply(filename1)
			handler.Apply(filename2)
			handler.Apply(filename3)
			time.Sleep(time.Second * 20)
			handler.WithGVK(gvk).Delete(name1)
			handler.WithGVK(gvk).Delete(name2)
			handler.WithGVK(gvk).Delete(name3)
		}
	}(ctx)

	timer := time.NewTimer(time.Second * 60)
	<-timer.C
	cancel()
	handler.WithGVK(gvk).Delete(name1)
	handler.WithGVK(gvk).Delete(name2)
	handler.WithGVK(gvk).Delete(name3)

	// Output:

	//2022/09/08 11:41:12 added Deployment: "test/mydep".
	//2022/09/08 11:41:12 added Deployment: "test/mydep2".
	//2022/09/08 11:41:12 added Deployment: "test/mydep-nolabel".
	//2022/09/08 11:41:15 modified Deployment: "test/mydep-nolabel".
	//2022/09/08 11:41:18 modified Deployment: "test/mydep".
	//2022/09/08 11:41:21 modified Deployment: "test/mydep2".
	//2022/09/08 11:41:27 deleted Deployment: "test/mydep".
	//2022/09/08 11:41:27 deleted Deployment: "test/mydep2".
	//2022/09/08 11:41:27 deleted Deployment: "test/mydep-nolabel".
	//2022/09/08 11:41:27 added Deployment: "test/mydep".
	//2022/09/08 11:41:27 added Deployment: "test/mydep2".
	//2022/09/08 11:41:27 added Deployment: "test/mydep-nolabel".
	//2022/09/08 11:41:27 modified Deployment: "test/mydep".
	//2022/09/08 11:41:27 modified Deployment: "test/mydep2".
	//2022/09/08 11:41:27 modified Deployment: "test/mydep-nolabel".
	//2022/09/08 11:41:27 modified Deployment: "test/mydep".
	//2022/09/08 11:41:27 modified Deployment: "test/mydep2".
	//2022/09/08 11:41:27 modified Deployment: "test/mydep-nolabel".
	//2022/09/08 11:41:27 modified Deployment: "test/mydep".
	//2022/09/08 11:41:27 modified Deployment: "test/mydep2".
	//2022/09/08 11:41:27 modified Deployment: "test/mydep-nolabel".
	//2022/09/08 11:41:27 modified Deployment: "test/mydep".
	//2022/09/08 11:41:27 modified Deployment: "test/mydep2".
	//2022/09/08 11:41:39 modified Deployment: "test/mydep".
	//2022/09/08 11:41:39 modified Deployment: "test/mydep".
	//2022/09/08 11:41:40 modified Deployment: "test/mydep2".
	//2022/09/08 11:41:45 modified Deployment: "test/mydep-nolabel".
	//2022/09/08 11:41:47 deleted Deployment: "test/mydep".
	//2022/09/08 11:41:47 deleted Deployment: "test/mydep2".
	//2022/09/08 11:41:47 deleted Deployment: "test/mydep-nolabel".
	//2022/09/08 11:41:47 added Deployment: "test/mydep".
	//2022/09/08 11:41:47 added Deployment: "test/mydep2".
	//2022/09/08 11:41:47 modified Deployment: "test/mydep".
	//2022/09/08 11:41:47 added Deployment: "test/mydep-nolabel".
	//2022/09/08 11:41:47 modified Deployment: "test/mydep2".
	//2022/09/08 11:41:47 modified Deployment: "test/mydep".
	//2022/09/08 11:41:47 modified Deployment: "test/mydep2".
	//2022/09/08 11:41:47 modified Deployment: "test/mydep-nolabel".
	//2022/09/08 11:41:48 modified Deployment: "test/mydep-nolabel".
	//2022/09/08 11:41:48 modified Deployment: "test/mydep".
	//2022/09/08 11:41:48 modified Deployment: "test/mydep2".
	//2022/09/08 11:41:48 modified Deployment: "test/mydep-nolabel".
	//2022/09/08 11:41:59 modified Deployment: "test/mydep".
	//2022/09/08 11:41:59 modified Deployment: "test/mydep".
	//2022/09/08 11:42:00 modified Deployment: "test/mydep2".
	//2022/09/08 11:42:04 modified Deployment: "test/mydep-nolabel".
	//2022/09/08 11:42:07 deleted Deployment: "test/mydep".
	//2022/09/08 11:42:07 deleted Deployment: "test/mydep2".
	//2022/09/08 11:42:07 deleted Deployment: "test/mydep-nolabel".
}
