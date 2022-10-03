package main

import (
	"context"
	"log"
	"time"

	"github.com/forbearing/k8s/deployment"
	"github.com/forbearing/k8s/dynamic"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func Dynamic_Watch_All() {
	filename1 := "../../testdata/examples/deployment.yaml"
	filename2 := "../../testdata/examples/deployment-2.yaml"
	filename3 := "../../testdata/examples/deployment-nolabel.yaml"
	name1 := "mydep"
	name2 := "mydep2"
	name3 := "mydep-nolabel"
	gvk := deployment.GVK

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
		handler.WithGVK(gvk).Watch(addFunc, modifyFunc, deleteFunc)
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

	//2022/09/08 11:47:14 added Deployment: "test/mydep2".
	//2022/09/08 11:47:14 added Deployment: "test/mydep-nolabel".
	//2022/09/08 11:47:14 added Deployment: "default/nginx".
	//2022/09/08 11:47:14 added Deployment: "kube-system/coredns".
	//2022/09/08 11:47:14 added Deployment: "local-path-storage/local-path-provisioner".
	//2022/09/08 11:47:14 added Deployment: "test/mydep".
	//2022/09/08 11:47:15 modified Deployment: "test/mydep2".
	//2022/09/08 11:47:18 modified Deployment: "test/mydep2".
	//2022/09/08 11:47:25 modified Deployment: "test/mydep-nolabel".
	//2022/09/08 11:47:28 modified Deployment: "test/mydep".
	//2022/09/08 11:47:29 deleted Deployment: "test/mydep".
	//2022/09/08 11:47:29 deleted Deployment: "test/mydep2".
	//2022/09/08 11:47:29 deleted Deployment: "test/mydep-nolabel".
	//2022/09/08 11:47:29 added Deployment: "test/mydep".
	//2022/09/08 11:47:29 added Deployment: "test/mydep2".
	//2022/09/08 11:47:29 modified Deployment: "test/mydep".
	//2022/09/08 11:47:29 added Deployment: "test/mydep-nolabel".
	//2022/09/08 11:47:29 modified Deployment: "test/mydep".
	//2022/09/08 11:47:29 modified Deployment: "test/mydep2".
	//2022/09/08 11:47:29 modified Deployment: "test/mydep".
	//2022/09/08 11:47:29 modified Deployment: "test/mydep-nolabel".
	//2022/09/08 11:47:29 modified Deployment: "test/mydep2".
	//2022/09/08 11:47:29 modified Deployment: "test/mydep-nolabel".
	//2022/09/08 11:47:29 modified Deployment: "test/mydep".
	//2022/09/08 11:47:29 modified Deployment: "test/mydep-nolabel".
	//2022/09/08 11:47:29 modified Deployment: "test/mydep2".
	//2022/09/08 11:47:40 modified Deployment: "test/mydep-nolabel".
	//2022/09/08 11:47:41 modified Deployment: "test/mydep2".
	//2022/09/08 11:47:43 modified Deployment: "test/mydep2".
	//2022/09/08 11:47:47 modified Deployment: "test/mydep-nolabel".
	//2022/09/08 11:47:49 deleted Deployment: "test/mydep".
	//2022/09/08 11:47:49 deleted Deployment: "test/mydep2".
	//2022/09/08 11:47:49 deleted Deployment: "test/mydep-nolabel".
	//2022/09/08 11:47:49 added Deployment: "test/mydep".
	//2022/09/08 11:47:49 added Deployment: "test/mydep2".
	//2022/09/08 11:47:49 modified Deployment: "test/mydep".
	//2022/09/08 11:47:49 added Deployment: "test/mydep-nolabel".
	//2022/09/08 11:47:49 modified Deployment: "test/mydep2".
	//2022/09/08 11:47:49 modified Deployment: "test/mydep-nolabel".
	//2022/09/08 11:47:49 modified Deployment: "test/mydep".
	//2022/09/08 11:47:49 modified Deployment: "test/mydep2".
	//2022/09/08 11:47:49 modified Deployment: "test/mydep-nolabel".
	//2022/09/08 11:47:49 modified Deployment: "test/mydep".
	//2022/09/08 11:47:49 modified Deployment: "test/mydep-nolabel".
	//2022/09/08 11:47:49 modified Deployment: "test/mydep2".
	//2022/09/08 11:48:02 modified Deployment: "test/mydep".
	//2022/09/08 11:48:05 modified Deployment: "test/mydep-nolabel".
	//2022/09/08 11:48:09 modified Deployment: "test/mydep-nolabel".
	//2022/09/08 11:48:09 deleted Deployment: "test/mydep".
	//2022/09/08 11:48:09 deleted Deployment: "test/mydep2".
	//2022/09/08 11:48:09 deleted Deployment: "test/mydep-nolabel".
}
