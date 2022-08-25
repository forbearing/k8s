package pod

import (
	"time"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	informerscore "k8s.io/client-go/informers/core/v1"
	listerscore "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
)

/*
ref:
	https://mp.weixin.qq.com/s/_mWiqvKeq-Uvu6QxE0f-qQ
	https://github.com/kubernetes-sigs/controller-runtime/issues/521

1.Informer 只会调用 kubernetes List 和 Watch 两种类型的 API. Informer 在初始化
  的时候, 先调用 Kubernetes List API 获得某种 resource 的全部 Object, 缓存在
  内存中,然后调用 Watch API 去 watch 这种 resource, 去维护这份缓存, 最后 Informer
  就不再调用 kubernetes 的任何 API.
2.用 List/Watch 去维护缓存,保持一致性是非常典型的做法,但令人费解的是, Informer
  只在初始化调用一次 List API, 之后完全依赖 Watch API 去维护缓存,没有任何 resync 机制.
3.按照多数人思路,通过 resync 机制,重新 List 一遍 resource 下所有 Object, 可以更好
  的保证 Informer 缓存和 Kubernetes 中数据的一致性.
4.咨询过Google 内部 Kubernetes开发人员之后，得到的回复是:
  在 Informer 设计之初,确实存在一个 relist 无法去执行 resync 操作,但后来被取消了,
  原因是现在的这种 List/Watch 机制,完全能够保证永远不会漏掉任何事件,因为完全没有
  必要再添加 relist 方法去 resync informer 的缓存. 这种做法也说明了 kubernetes
  完全信任 etcd.
*/

// SetInformerResyncPeriod will set informer resync period.
func (h *Handler) SetInformerResyncPeriod(resyncPeriod time.Duration) {
	h.informerFactory = informers.NewSharedInformerFactory(h.clientset, resyncPeriod)
}

// InformerFactory returns underlying SharedInformerFactory which provides
// shared informer for resources in all known API group version.
func (h *Handler) InformerFactory() informers.SharedInformerFactory {
	return h.informerFactory
}

// PodInformer returns underlying PodInformer which provides access to a shared
// informer and lister for pod.
func (h *Handler) PodInformer() informerscore.PodInformer {
	return h.informerFactory.Core().V1().Pods()
}

// Informer returns underlying SharedIndexInformer which provides add and Indexers
// ability based on SharedInformer.
func (h *Handler) Informer() cache.SharedIndexInformer {
	return h.informerFactory.Core().V1().Pods().Informer()
}

// Lister returns underlying PodLister which helps list pods.
func (h *Handler) Lister() listerscore.PodLister {
	return h.informerFactory.Core().V1().Pods().Lister()
}

// TestInformer
func (h *Handler) TestInformer(stopCh chan struct{}) {
	h.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
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
	h.InformerFactory().Start(stopCh)
	logrus.Info("Waiting for informer caches to sync")
	if ok := cache.WaitForCacheSync(stopCh, h.Informer().HasSynced); !ok {
		logrus.Error("failed to wait for caches to sync")
	}
}

// RunInformer start and run the shared informer, returning after it stops.
// The informer will be stopped when stopCh is closed.
//
// AddFunc, updateFunc, and deleteFunc are used to handle add, update,
// and delete event of k8s pod resource, respectively.
func (h *Handler) RunInformer(
	stopCh chan struct{},
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
	addFunc func(obj interface{}),
	updateFunc func(oldObj, newObj interface{}),
	deleteFunc func(obj interface{}),
	stopCh chan struct{}) {

	h.RunInformer(stopCh, addFunc, updateFunc, deleteFunc)
}

// 1.PodInformer 继承了 SharedIndexInformer 和 PodLister
// 2.SharedIndexInformer 继承了 SharedInformer
// 3.SharedInformer 拥有 AddEventHandler, AddEventHandlerWithResyncPeriod, HasSynced 等各种方法
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
