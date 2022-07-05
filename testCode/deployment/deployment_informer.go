package main

import (
	"reflect"

	"github.com/forbearing/k8s/deployment"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type updateObj struct {
	oldObj interface{}
	newObj interface{}
}

func Deployment_Informer() {
	handler, err := deployment.New(ctx, namespace, kubeconfig)
	if err != nil {
		panic(err)
	}
	defer cleanup(handler)
	addQueue := make(chan interface{}, 100)
	updateQueue := make(chan updateObj, 100)
	deleteQueue := make(chan interface{}, 100)

	addFunc := func(obj interface{}) { addQueue <- obj }
	updateFunc := func(oldObj interface{}, newObj interface{}) {
		uo := updateObj{}
		uo.oldObj = oldObj
		uo.newObj = newObj
		updateQueue <- uo
	}
	deleteFunc := func(obj interface{}) { deleteQueue <- obj }
	stopCh := make(chan struct{}, 1)

	go func() {
		handler.RunInformer(addFunc, updateFunc, deleteFunc, stopCh)
	}()

	for {
		select {
		case obj := <-addQueue:
			myObj := obj.(metav1.Object)
			log.Println("Add: ", myObj.GetName())
		case ou := <-updateQueue:
			oldObj := ou.oldObj.(metav1.Object)
			newObj := ou.newObj.(metav1.Object)
			if !reflect.DeepEqual(oldObj, newObj) {
				log.Println("Update:", newObj.GetName())
			}
		case obj := <-deleteQueue:
			myObj := obj.(metav1.Object)
			log.Println("Delete", myObj.GetName())
		case <-stopCh:
			log.Println("Informer stopped.")
		}

	}
}
