package pod

import (
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	listerscorev1 "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
)

func (h *Handler) TestInformer(stopCh chan struct{}) {
	h.informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			myObj := obj.(metav1.Object)
			log.Printf("New Pod Added to Store: %s", myObj.GetName())
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			newPod := newObj.(*corev1.Pod)
			oldPod := oldObj.(*corev1.Pod)
			if newPod.ResourceVersion != oldPod.ResourceVersion {
				log.Printf("Pod Updated to Store: %s\n", newPod.Name)
			}

			//if !reflect.DeepEqual(newObj, oldObj) {
			//    log.Printf("Pod Updated to Store: %s\n", newObj.(metav1.Object).GetName())
			//}
		},
		DeleteFunc: func(obj interface{}) {
			myObj := obj.(metav1.Object)
			log.Printf("Pod Deleted from Store: %s", myObj.GetName())
		},
	})
	h.informerFactory.WaitForCacheSync(stopCh)
	//cache.WaitForCacheSync(stopCh, h.informer.HasSynced)
	h.informer.Run(stopCh)
}

// addFunc, updateFunc, stopChan
// informer 的三个回调函数 addFunc, updateFunc, deleteFunc
// 这个管道用来存放回调函数处理的 k8s 资源对象
// RunInformer
func (h *Handler) RunInformer(
	addFunc func(obj interface{}),
	updateFunc func(oldObj, newObj interface{}),
	deleteFunc func(obj interface{}),
	stopCh chan struct{}) {
	h.informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    addFunc,
		UpdateFunc: updateFunc,
		DeleteFunc: deleteFunc,
	})
	h.informerFactory.WaitForCacheSync(stopCh)
	//cache.WaitForCacheSync(stopCh, h.informer.HasSynced)
	h.informer.Run(stopCh)

	//h.informer.AddIndexers(cache.Indexers{
	//    controllerUIDIndex: func(obj interface{}) ([]string, error) {
	//        pod, ok := obj.(*corev1.Pod)
	//        if ok {
	//            return []string{}, nil
	//        }
	//    },
	//})
}

//func (h *Handler) Informer(resync time.Duration) informerscorev1.PodInformer {
//    return informers.NewSharedInformerFactory(h.clientset, resync).Core().V1().Pods()
//}

func (h *Handler) Informer() cache.SharedIndexInformer {
	return h.informer
}
func (h *Handler) Lister() listerscorev1.PodLister {
	return h.lister
}

// 1.PodInformer 继承了 SharedIndexInformer 和 PodLister
// 2.SharedIndexInformer 继承了 SharedInformer
// 3.SharedInformer 拥有 AddEventHandler, AddEventHandlerWithResyncPeriod  等各种方法
// 4.最终关系是: PodInformer -> SharedIndexInformer -> SharedInformer
//   SharedIndexInformer 和 PodLister 是同级别的关系

//// PodInformer provides access to a shared informer and lister for Pods.
// informerFactory = informers.NewSharedInformerFactory(clientset, time.Minute)
// PodInformer = informerFactory.Core().V1().Pods().Informer()
//type PodInformer interface {
//    Informer() cache.SharedIndexInformer
//    Lister() v1.PodLister
//}

//// SharedIndexInformer provides add and get Indexers ability based on SharedInformer.
//type SharedIndexInformer interface {
//    SharedInformer
//    // AddIndexers add indexers to the informer before it starts.
//    AddIndexers(indexers Indexers) error
//    GetIndexer() Indexer
//}

//type SharedInformer interface {
//    // AddEventHandler adds an event handler to the shared informer using the shared informer's resync
//    // period.  Events to a single handler are delivered sequentially, but there is no coordination
//    // between different handlers.
//    AddEventHandler(handler ResourceEventHandler)
//    // AddEventHandlerWithResyncPeriod adds an event handler to the
//    // shared informer with the requested resync period; zero means
//    // this handler does not care about resyncs.  The resync operation
//    // consists of delivering to the handler an update notification
//    // for every object in the informer's local cache; it does not add
//    // any interactions with the authoritative storage.  Some
//    // informers do no resyncs at all, not even for handlers added
//    // with a non-zero resyncPeriod.  For an informer that does
//    // resyncs, and for each handler that requests resyncs, that
//    // informer develops a nominal resync period that is no shorter
//    // than the requested period but may be longer.  The actual time
//    // between any two resyncs may be longer than the nominal period
//    // because the implementation takes time to do work and there may
//    // be competing load and scheduling noise.
//    AddEventHandlerWithResyncPeriod(handler ResourceEventHandler, resyncPeriod time.Duration)
//    // GetStore returns the informer's local cache as a Store.
//    GetStore() Store
//    // GetController is deprecated, it does nothing useful
//    GetController() Controller
//    // Run starts and runs the shared informer, returning after it stops.
//    // The informer will be stopped when stopCh is closed.
//    Run(stopCh <-chan struct{})
//    // HasSynced returns true if the shared informer's store has been
//    // informed by at least one full LIST of the authoritative state
//    // of the informer's object collection.  This is unrelated to "resync".
//    HasSynced() bool
//    // LastSyncResourceVersion is the resource version observed when last synced with the underlying
//    // store. The value returned is not synchronized with access to the underlying store and is not
//    // thread-safe.
//    LastSyncResourceVersion() string

//    // The WatchErrorHandler is called whenever ListAndWatch drops the
//    // connection with an error. After calling this handler, the informer
//    // will backoff and retry.
//    //
//    // The default implementation looks at the error type and tries to log
//    // the error message at an appropriate level.
//    //
//    // There's only one handler, so if you call this multiple times, last one
//    // wins; calling after the informer has been started returns an error.
//    //
//    // The handler is intended for visibility, not to e.g. pause the consumers.
//    // The handler should return quickly - any expensive processing should be
//    // offloaded.
//    SetWatchErrorHandler(handler WatchErrorHandler) error

//    // The TransformFunc is called for each object which is about to be stored.
//    //
//    // This function is intended for you to take the opportunity to
//    // remove, transform, or normalize fields. One use case is to strip unused
//    // metadata fields out of objects to save on RAM cost.
//    //
//    // Must be set before starting the informer.
//    //
//    // Note: Since the object given to the handler may be already shared with
//    //	other goroutines, it is advisable to copy the object being
//    //  transform before mutating it at all and returning the copy to prevent
//    //	data races.
//    SetTransform(handler TransformFunc) error
//}

//type PodLister interface {
//    // List lists all Pods in the indexer.
//    // Objects returned here must be treated as read-only.
//    List(selector labels.Selector) (ret []*v1.Pod, err error)
//    // Pods returns an object that can list and get Pods.
//    Pods(namespace string) PodNamespaceLister
//    PodListerExpansion
//}
