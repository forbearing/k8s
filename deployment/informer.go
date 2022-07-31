package deployment

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

// DeploymentInformer returns underlying DeploymentInformer which provides
// access to a shared informer and lister for deployment.
func (h *Handler) DeploymentInformer() informersapps.DeploymentInformer {
	return h.informerFactory.Apps().V1().Deployments()
}

// Informer returns underlying SharedIndexInformer which provides add and Indexers
// ability based on SharedInformer.
func (h *Handler) Informer() cache.SharedIndexInformer {
	return h.informerFactory.Apps().V1().Deployments().Informer()
}

// Lister returns underlying DeploymentLister which helps list deployments.
func (h *Handler) Lister() listersapps.DeploymentLister {
	return h.informerFactory.Apps().V1().Deployments().Lister()
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
	h.informerFactory.WaitForCacheSync(stopCh)
	h.Informer().Run(stopCh)
}
