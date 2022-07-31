package daemonset

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

// DaemonSetInformer returns underlying DaemonSetInformer which provides access
// to a shared informer and lister for daemonset.
func (h *Handler) DaemonSetInformer() informersapps.DaemonSetInformer {
	return h.informerFactory.Apps().V1().DaemonSets()
}

// Informer returns underlying SharedIndexInformer which provides add and Indexers
// ability based on SharedInformer.
func (h *Handler) Informer() cache.SharedIndexInformer {
	return h.informerFactory.Apps().V1().DaemonSets().Informer()
}

// Lister returns underlying DaemonSetLister which helps list daemonsets.
func (h *Handler) Lister() listersapps.DaemonSetLister {
	return h.informerFactory.Apps().V1().DaemonSets().Lister()
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
