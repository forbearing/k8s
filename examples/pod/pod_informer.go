package main

import (
	"reflect"

	"github.com/forbearing/k8s/pod"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type updateObject struct {
	oldObj interface{}
	newObj interface{}
}

func Pod_Informer() {
	handler, err := pod.New(ctx, namespace, kubeconfig)
	if err != nil {
		panic(err)
	}
	defer cleanup(handler)

	addQueue := make(chan interface{}, 100)
	updateQueue := make(chan updateObject, 100)
	deleteQueue := make(chan interface{}, 100)

	addFunc := func(obj interface{}) { addQueue <- obj }
	updateFunc := func(oldObj, newObj interface{}) {
		uo := updateObject{}
		uo.oldObj = oldObj
		uo.newObj = newObj
		updateQueue <- uo
	}
	deleteFunc := func(obj interface{}) { deleteQueue <- obj }
	stopCh := make(chan struct{}, 1)

	// RunInformer 必须开启一个新的 goroutine 来执行
	go func() {
		handler.RunInformer(addFunc, updateFunc, deleteFunc, stopCh)
	}()

	for {
		select {
		case obj := <-addQueue:
			myObj := obj.(metav1.Object)
			log.Printf("New Pod Added to Store: %s\n", myObj.GetName())
		case uo := <-updateQueue:
			oObj := uo.oldObj.(metav1.Object)
			nObj := uo.newObj.(metav1.Object)
			_, _ = oObj, nObj
			// OnUpdate is called with all existing objects on the specific resync interval
			// OnUpdate is called with all existing objects when a watch connection
			// is dropped by the server and a full relist is required
			if !reflect.DeepEqual(oObj, nObj) {
				log.Printf("Pod Updated to Store: %s\n", nObj.GetName())
			}
		case obj := <-deleteQueue:
			myObj := obj.(metav1.Object)
			log.Printf("Pod Deleted from Store: %s\n", myObj.GetName())
		case <-stopCh:
			log.Println("informer stopped")
			return
		}
	}
}
