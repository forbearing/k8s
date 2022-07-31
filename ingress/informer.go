package ingress

import (
	"time"

	"k8s.io/client-go/informers"
	informersnetworking "k8s.io/client-go/informers/networking/v1"
	listersnetworking "k8s.io/client-go/listers/networking/v1"
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

// IngressInformer returns underlying IngressInformer which provides access to
// a shared informer and lister for ingress.
func (h *Handler) IngressInformer() informersnetworking.IngressInformer {
	return h.informerFactory.Networking().V1().Ingresses()
}

// Informer returns underlying SharedIndexInformer which provides add and Indexers
// ability based on SharedInformer.
func (h *Handler) Informer() cache.SharedIndexInformer {
	return h.informerFactory.Networking().V1().Ingresses().Informer()
}

// Lister returns underlying IngressLister which helps list ingresses.
func (h *Handler) Lister() listersnetworking.IngressLister {
	return h.informerFactory.Networking().V1().Ingresses().Lister()
}

// RunInformer.
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
