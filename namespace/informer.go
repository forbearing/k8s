package namespace

import (
	"time"

	"k8s.io/client-go/informers"
	informerscore "k8s.io/client-go/informers/core/v1"
	listerscore "k8s.io/client-go/listers/core/v1"
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

// NamespaceInformer returns underlying NamespaceInformer which provides access
// to a shared informer and lister for namespace.
func (h *Handler) NamespaceInformer() informerscore.NamespaceInformer {
	return h.informerFactory.Core().V1().Namespaces()
}

// Informer returns underlying SharedIndexInformer which provides add and Indexers
// ability based on SharedInformer.
func (h *Handler) Informer() cache.SharedIndexInformer {
	return h.informerFactory.Core().V1().Namespaces().Informer()
}

// Lister returns underlying NamespaceLister which helps list namespaces.
func (h *Handler) Lister() listerscore.NamespaceLister {
	return h.informerFactory.Core().V1().Namespaces().Lister()
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
