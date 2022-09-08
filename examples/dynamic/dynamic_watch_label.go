package main

import (
	"context"
	"log"
	"time"

	"github.com/forbearing/k8s/deployment"
	"github.com/forbearing/k8s/dynamic"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func Dynamic_Watch_Label() {
	filename1 := "../../testdata/examples/deployment.yaml"
	filename2 := "../../testdata/examples/deployment-2.yaml"
	filename3 := "../../testdata/examples/deployment-nolabel.yaml"
	name1 := "mydep"
	name2 := "mydep2"
	name3 := "mydep-nolabel"
	label := "type=deployment"
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
		handler.WithGVK(gvk).WatchByLabel(label, addFunc, modifyFunc, deleteFunc)
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

	//2022/09/08 11:34:03 added Deployment: "test/mydep".
	//2022/09/08 11:34:03 added Deployment: "test/mydep2".
	//2022/09/08 11:34:03 modified Deployment: "test/mydep2".
	//2022/09/08 11:34:08 modified Deployment: "test/mydep".
	//2022/09/08 11:34:12 modified Deployment: "test/mydep".
	//2022/09/08 11:34:18 deleted Deployment: "test/mydep".
	//2022/09/08 11:34:18 deleted Deployment: "test/mydep2".
	//2022/09/08 11:34:18 added Deployment: "test/mydep".
	//2022/09/08 11:34:18 added Deployment: "test/mydep2".
	//2022/09/08 11:34:18 modified Deployment: "test/mydep".
	//2022/09/08 11:34:18 modified Deployment: "test/mydep2".
	//2022/09/08 11:34:18 modified Deployment: "test/mydep".
	//2022/09/08 11:34:18 modified Deployment: "test/mydep2".
	//2022/09/08 11:34:18 modified Deployment: "test/mydep".
	//2022/09/08 11:34:18 modified Deployment: "test/mydep2".
	//2022/09/08 11:34:31 modified Deployment: "test/mydep2".
	//2022/09/08 11:34:33 modified Deployment: "test/mydep".
	//2022/09/08 11:34:38 deleted Deployment: "test/mydep".
	//2022/09/08 11:34:38 deleted Deployment: "test/mydep2".
	//2022/09/08 11:34:38 added Deployment: "test/mydep".
	//2022/09/08 11:34:38 added Deployment: "test/mydep2".
	//2022/09/08 11:34:38 modified Deployment: "test/mydep2".
	//2022/09/08 11:34:38 modified Deployment: "test/mydep".
	//2022/09/08 11:34:38 modified Deployment: "test/mydep2".
	//2022/09/08 11:34:38 modified Deployment: "test/mydep".
	//2022/09/08 11:34:38 modified Deployment: "test/mydep2".
	//2022/09/08 11:34:38 modified Deployment: "test/mydep".
	//2022/09/08 11:34:38 modified Deployment: "test/mydep".
	//2022/09/08 11:34:38 modified Deployment: "test/mydep2".
	//2022/09/08 11:34:52 modified Deployment: "test/mydep2".
	//2022/09/08 11:34:58 deleted Deployment: "test/mydep".
	//2022/09/08 11:34:58 deleted Deployment: "test/mydep2".
}
