package cronjob

import (
	"time"

	"k8s.io/client-go/informers"
	informersbatch "k8s.io/client-go/informers/batch/v1"
	listersbatch "k8s.io/client-go/listers/batch/v1"
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

// CronJobInformer returns underlying CronJobInformer which provides access to
// a shared informer and lister for cronjob.
func (h *Handler) CronJobInformer() informersbatch.CronJobInformer {
	return h.informerFactory.Batch().V1().CronJobs()
}

// Informer returns underlying SharedIndexInformer which provides add and Indexers
// ability based on SharedInformer.
func (h *Handler) Informer() cache.SharedIndexInformer {
	return h.informerFactory.Batch().V1().CronJobs().Informer()
}

// Lister returns underlying CronJobLister which helps list cronjobs.
func (h *Handler) Lister() listersbatch.CronJobLister {
	return h.informerFactory.Batch().V1().CronJobs().Lister()
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
