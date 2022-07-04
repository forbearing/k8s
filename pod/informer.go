package pod

import (
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
)

func (h *Handler) TestInformer(stopCh chan struct{}) {
	h.informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			myObj := obj.(metav1.Object)
			log.Infof("New Pod Added to Store: %s", myObj.GetName())
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			oObj := oldObj.(metav1.Object)
			nObj := newObj.(metav1.Object)
			log.Infof("%s Pod Updated to %s", oObj.GetName(), nObj.GetName())
		},
		DeleteFunc: func(obj interface{}) {
			myObj := obj.(metav1.Object)
			log.Infof("Pod Deleted from Store: %s", myObj.GetName())
		},
	})
	h.informer.Run(stopCh)
}

// addFunc, updateFunc, stopChan
// informer 的三个回调函数 addFunc, updateFunc, deleteFunc
// 这个管道用来存放回调函数处理的 k8s 资源对象
// RunInformer
func (h *Handler) RunInformer(
	addFunc func(obj interface{}),
	updateFunc func(oldObj, newObj interface{}),
	deleteFunc func(obj interface{}),
	stopCh chan struct{}) {
	h.informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    addFunc,
		UpdateFunc: updateFunc,
		DeleteFunc: deleteFunc,
	})
	h.informer.Run(stopCh)
}
