package clusterrolebinding

import (
	"time"

	"k8s.io/client-go/informers"
	informersrbac "k8s.io/client-go/informers/rbac/v1"
	listersrbac "k8s.io/client-go/listers/rbac/v1"
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

// ClusterRoleBindingInformer returns underlying ClusterRoleBindingInformer
// which provides access to a shared informer and lister for clusterrolebinding.
func (h *Handler) ClusterRoleBindingInformer() informersrbac.ClusterRoleBindingInformer {
	return h.informerFactory.Rbac().V1().ClusterRoleBindings()
}

// Informer returns underlying SharedIndexInformer which provides add and Indexers
// ability based on SharedInformer.
func (h *Handler) Informer() cache.SharedIndexInformer {
	return h.informerFactory.Rbac().V1().ClusterRoleBindings().Informer()
}

// Lister returns underlying ClusterRoleBindingLister which helps list clusterrolebindings.
func (h *Handler) Lister() listersrbac.ClusterRoleBindingLister {
	return h.informerFactory.Rbac().V1().ClusterRoleBindings().Lister()
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
