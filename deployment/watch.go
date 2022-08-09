package deployment

import (
	log "github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

// WatchByName watch deployment by name.
func (h *Handler) WatchByName(name string,
	addFunc, modifyFunc, deleteFunc func(x interface{}), x interface{}) (err error) {
	var (
		watcher watch.Interface
		timeout = int64(0)
		isExist bool
	)

	// if event channel is closed, it means the server has closed the connection,
	// reconnect to kubernetes.
	for {
		//watcher, err := clientset.AppsV1().Deployments(namespace).Watch(ctx,
		//    metav1.SingleObject(metav1.ObjectMeta{Name: "dep", Namespace: namespace}))
		listOptions := metav1.SingleObject(metav1.ObjectMeta{Name: name, Namespace: h.namespace})
		listOptions.TimeoutSeconds = &timeout
		if watcher, err = h.clientset.AppsV1().Deployments(h.namespace).Watch(h.ctx, listOptions); err != nil {
			return
		}
		if _, err = h.Get(name); err != nil {
			isExist = false // deployment not exist
		} else {
			isExist = true // deployment exist
		}
		//for {
		//    // kubernetes retains the resource event history, which includes this
		//    // initial event, so that when our program first start, we are automatically
		//    // notified of the deployment existence and current state.
		//    event, isOpen := <-watcher.ResultChan()

		//    if isOpen {
		//        switch event.Type {
		//        case watch.Added:
		//            // if deployment exist, skip deployment history add event.
		//            if !isExist {
		//                addFunc()
		//            }
		//            isExist = true
		//        case watch.Modified:
		//            modifyFunc()
		//            isExist = true
		//        case watch.Deleted:
		//            deleteFunc()
		//            isExist = false
		//        //case watch.Bookmark:
		//        //    log.Debug("bookmark")
		//        //case watch.Error:
		//        //    log.Error("error")
		//        default: // do nothing
		//        }
		//    } else {
		//        // If event channel is closed, it means the server has closed the connection
		//        log.Debug("reconnect to kubernetes.")
		//        break
		//    }
		//}
		for event := range watcher.ResultChan() {
			switch event.Type {
			case watch.Added:
				if !isExist {
					addFunc(x)
				}
				isExist = true
			case watch.Modified:
				modifyFunc(x)
				isExist = true
			case watch.Deleted:
				deleteFunc(x)
				isExist = false
			case watch.Bookmark:
				log.Debug("watch deployment: bookmark.")
			case watch.Error:
				log.Debug("watch deployment: error")
			}
		}
		// If event channel is closed, it means the server has closed the connection
		log.Debug("watch deployment: reconnect to kubernetes")
		watcher.Stop()
	}
}

// WatchByLabel watch deployment by labels.
func (h *Handler) WatchByLabel(labels string,
	addFunc, modifyFunc, deleteFunc func(x interface{}), x interface{}) (err error) {
	var (
		watcher    watch.Interface
		timeout    = int64(0)
		isExist    bool
		deployList []*appsv1.Deployment
	)
	// if event channel is closed, it means the server has closed the connection,
	// reconnect to kubernetes.
	for {
		//watcher, err := clientset.AppsV1().Deployments(namespace).Watch(ctx,
		//    metav1.SingleObject(metav1.ObjectMeta{Name: "dep", Namespace: namespace}))
		// 这个 timeout 一定要设置为 0, 否则 watcher 就会中断
		if watcher, err = h.clientset.AppsV1().Deployments(h.namespace).Watch(h.ctx,
			metav1.ListOptions{LabelSelector: labels, TimeoutSeconds: &timeout}); err != nil {
			return
		}
		if deployList, err = h.ListByLabel(labels); err != nil {
			return
		}
		if len(deployList) == 0 {
			isExist = false // deployment not exist
		} else {
			isExist = true // deployment exist
		}
		//for {
		//    // kubernetes retains the resource event history, which includes this
		//    // initial event, so that when our program first start, we are automatically
		//    // notified of the deployment existence and current state.
		//    event, isOpen := <-watcher.ResultChan()

		//    if isOpen {
		//        switch event.Type {
		//        case watch.Added:
		//            // if deployment exist, skip deployment history add event.
		//            if !isExist {
		//                addFunc(x)
		//            }
		//            isExist = true
		//        case watch.Modified:
		//            modifyFunc(x)
		//            isExist = true
		//        case watch.Deleted:
		//            deleteFunc(x)
		//            isExist = false
		//        //case watch.Bookmark:
		//        //    log.Debug("bookmark")
		//        //case watch.Error:
		//        //    log.Error("error")
		//        default: // do nothing
		//        }
		//    } else {
		//        // If event channel is closed, it means the server has closed the connection
		//        log.Debug("reconnect to kubernetes.")
		//        break
		//    }
		//}
		for event := range watcher.ResultChan() {
			switch event.Type {
			case watch.Added:
				if !isExist {
					addFunc(x)
				}
				isExist = true
			case watch.Modified:
				modifyFunc(x)
				isExist = true
			case watch.Deleted:
				deleteFunc(x)
				isExist = false
			case watch.Bookmark:
				log.Debug("watch deployment: bookmark.")
			case watch.Error:
				log.Debug("watch deployment: error")
			}
		}
		// If event channel is closed, it means the server has closed the connection
		log.Debug("watch deployment: reconnect to kubernetes")
		watcher.Stop()
	}
}

// Watch watch deployment by label, alias to "WatchByLabel".
func (h *Handler) Watch(name string,
	addFunc, modifyFunc, deleteFunc func(x interface{}), x interface{}) (err error) {
	return h.WatchByName(name, addFunc, modifyFunc, deleteFunc, x)
}
