package job

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

// JobInformer returns underlying JobInformer which provides access to a shared
// informer and lister for job.
func (h *Handler) JobInformer() informersbatch.JobInformer {
	return h.informerFactory.Batch().V1().Jobs()
}

// Informer returns underlying SharedIndexInformer which provides add and Indexers
// ability based on SharedInformer.
func (h *Handler) Informer() cache.SharedIndexInformer {
	return h.informerFactory.Batch().V1().Jobs().Informer()
}

// Lister returns underlying JobLister which helps list jobs.
func (h *Handler) Lister() listersbatch.JobLister {
	return h.informerFactory.Batch().V1().Jobs().Lister()
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
