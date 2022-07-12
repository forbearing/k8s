package deployment

import (
	listersappsv1 "k8s.io/client-go/listers/apps/v1"
	"k8s.io/client-go/tools/cache"
)

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
	h.informerFactory.WaitForCacheSync(stopCh)
	h.informer.Run(stopCh)
}

//func (h *Handler) InformerFactory(resync time.Duration) informers.SharedInformerFactory {
//    return informers.NewSharedInformerFactory(h.clientset, resync)
//}

func (h *Handler) Informer() cache.SharedIndexInformer {
	return h.informer
}
func (h *Handler) Lister() listersappsv1.DeploymentLister {
	return h.lister
}
