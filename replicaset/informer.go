package replicaset

import (
	"time"

	"k8s.io/client-go/informers"
	informersapps "k8s.io/client-go/informers/apps/v1"
	listersapps "k8s.io/client-go/listers/apps/v1"
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

// ReplicaSetInformer returns underlying ReplicaSetInformer which provides access
// to a shared informer and lister for replicaset.
func (h *Handler) ReplicaSetInformer() informersapps.ReplicaSetInformer {
	return h.informerFactory.Apps().V1().ReplicaSets()
}

// Informer returns underlying SharedIndexInformer which provides add and Indexers
// ability based on SharedInformer.
func (h *Handler) Informer() cache.SharedIndexInformer {
	return h.informerFactory.Apps().V1().ReplicaSets().Informer()
}

// Lister returns underlying ReplicaSetLister which helps list replicasets.
func (h *Handler) Lister() listersapps.ReplicaSetLister {
	return h.informerFactory.Apps().V1().ReplicaSets().Lister()
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
