package storageclass

import (
	"time"

	"k8s.io/client-go/informers"
	informersstorage "k8s.io/client-go/informers/storage/v1"
	listersstorage "k8s.io/client-go/listers/storage/v1"
	"k8s.io/client-go/tools/cache"
)

// SetInformerResyncPeriod will set informer resync period.
func (h *Handler) SetInformerResyncPeriod(resyncPeriod time.Duration) {
	h.informerFactory = informers.NewSharedInformerFactory(h.clientset, resyncPeriod)
}

// InformerFactory returns underlying SharedInformerFactory which provides
// shared informer for resources in all known API group version.
func (h *Handler) InformerFactory() informers.SharedInformerFactory {
	return h.informerFactory
}

// StorageClassInformer returns underlying StorageClassInformer which provides
// access to a shared informer and lister for storageclass.
func (h *Handler) StorageClassInformer() informersstorage.StorageClassInformer {
	return h.informerFactory.Storage().V1().StorageClasses()
}

// Informer returns underlying SharedIndexInformer which provides add and Indexers
// ability based on SharedInformer.
func (h *Handler) Informer() cache.SharedIndexInformer {
	return h.informerFactory.Storage().V1().StorageClasses().Informer()
}

// Lister returns underlying StorageClassLister which helps list storageclasss.
func (h *Handler) Lister() listersstorage.StorageClassLister {
	return h.informerFactory.Storage().V1().StorageClasses().Lister()
}

// RunInformer
func (h *Handler) RunInformer(
	addFunc func(obj interface{}),
	updateFunc func(oldObj, newObj interface{}),
	deleteFunc func(obj interface{}),
	stopCh chan struct{}) {
	h.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    addFunc,
		UpdateFunc: updateFunc,
		DeleteFunc: deleteFunc,
	})
	h.Informer().Run(stopCh)
}
