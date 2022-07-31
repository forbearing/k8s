package role

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

// RoleInformer returns underlying RoleInformer which provides access to a shared
// informer and lister for role.
func (h *Handler) RoleInformer() informersrbac.RoleInformer {
	return h.informerFactory.Rbac().V1().Roles()
}

// Informer returns underlying SharedIndexInformer which provides add and Indexers
// ability based on SharedInformer.
func (h *Handler) Informer() cache.SharedIndexInformer {
	return h.informerFactory.Rbac().V1().Roles().Informer()
}

// Lister returns underlying RoleLister which helps list roles.
func (h *Handler) Lister() listersrbac.RoleLister {
	return h.informerFactory.Rbac().V1().Roles().Lister()
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
