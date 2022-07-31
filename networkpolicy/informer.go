package networkpolicy

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

// NetworkPolicyInformer returns underlying NetworkPolicyInformer which provides
// access to a shared informer and lister for networkpolicy.
func (h *Handler) NetworkPolicyInformer() informersnetworking.NetworkPolicyInformer {
	return h.informerFactory.Networking().V1().NetworkPolicies()
}

// Informer returns underlying SharedIndexInformer which provides add and Indexers
// ability based on SharedInformer.
func (h *Handler) Informer() cache.SharedIndexInformer {
	return h.informerFactory.Networking().V1().NetworkPolicies().Informer()
}

// Lister returns underlying NetworkPolicyLister which helps list networkpolicies.
func (h *Handler) Lister() listersnetworking.NetworkPolicyLister {
	return h.informerFactory.Networking().V1().NetworkPolicies().Lister()
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
