package deployment

import (
	"time"

	"github.com/sirupsen/logrus"
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

// RunInformer start and run the shared informer, returning after it stops.
// The informer will be stopped when stopCh is closed.
//
// AddFunc, updateFunc, and deleteFunc are used to handle add, update,
// and delete event of k8s deployment resource, respectively.
func (h *Handler) RunInformer(
	stopCh <-chan struct{},
	addFunc func(obj interface{}),
	updateFunc func(oldObj, newObj interface{}),
	deleteFunc func(obj interface{})) {

	h.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    addFunc,
		UpdateFunc: updateFunc,
		DeleteFunc: deleteFunc,
	})

	// method 1, recommended
	h.InformerFactory().Start(stopCh)
	logrus.Info("Waiting for informer caches to sync")
	if ok := cache.WaitForCacheSync(stopCh, h.Informer().HasSynced); !ok {
		logrus.Error("failed to wait for caches to sync")
	}

	//// method 2
	//h.InformerFactory().Start(stopCh)
	//logrus.Info("Waiting for informer caches to sync")
	//h.InformerFactory().WaitForCacheSync(stopCh)

	//// method 3
	//logrus.Info("Waiting for informer caches to sync")
	//h.informerFactory.WaitForCacheSync(stopCh)
	//h.Informer().Run(stopCh)
}

// StartInformer simply call RunInformer.
func (h *Handler) StartInformer(
	stopCh <-chan struct{},
	addFunc func(obj interface{}),
	updateFunc func(oldObj, newObj interface{}),
	deleteFunc func(obj interface{})) {

	h.RunInformer(stopCh, addFunc, updateFunc, deleteFunc)
}
