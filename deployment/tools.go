package deployment

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

// IsReady check if the deployment is ready.
func (h *Handler) IsReady(name string) bool {
	deploy, err := h.Get(name)
	if err != nil {
		return false
	}
	// 必须 Type=Available 和 Status=True 才能算 Deployment 就绪了
	for _, cond := range deploy.Status.Conditions {
		if cond.Type == appsv1.DeploymentAvailable && cond.Status == corev1.ConditionTrue {
			return true
		}
	}
	return false
}

// WaitReady wait for the deployment to be in the ready state.
func (h *Handler) WaitReady(name string) (err error) {
	var (
		watcher watch.Interface
		timeout = int64(0)
	)
	// 在 watch deployment 之前先判断 deployment 是否就绪, 如果 deployment 已经
	// 就绪了,就没必要 watch 了.
	if h.IsReady(name) {
		return nil
	}
	// 检查 deploymen 是否存在,如果存在则不用 watch
	if _, err = h.Get(name); err != nil {
		return err
	}

	for {
		// 开始监听 deployment 事件,
		// 1. 如果监听到了 modified 事件, 就检查下 deployment 的状态.
		//    如果 conditions.Type == appsv1.DeploymentAvailable 并且 conditions.Status == corev1.ConditionTrue
		//    说明 deployment 已经准备好了.
		// 2. 如果监听到 watch.Deleted 事件, 说明 deployment 已经删除了, 不需要再监听了
		// 3. 注意这个 watcher 的创建必须在第一层循环之内,如果 watcher 在循环外面
		//    当 apiserver timeout, 这个 watcher 会被 close, 即使再循环也是在使用
		//    一个被 close 掉的 watcher, 就无法再监听 deployment.
		listOptions := metav1.SingleObject(metav1.ObjectMeta{Name: name, Namespace: h.namespace})
		listOptions.TimeoutSeconds = &timeout
		watcher, err = h.clientset.AppsV1().Deployments(h.namespace).Watch(h.ctx, listOptions)
		if err != nil {
			return err
		}
		// 连接 kubernetes 是有通过 http/https 方式连接的, 有一个 keepalived 的时间
		// 时间一到, 就会断开 kubernetes 的连接, 此时  watch.ResultChan 通道就会关闭.
		// 所以说, 这个方法 WaitReady 等待 deployment 处于就绪状态的最长等待时间就是
		// 连接 kubernetes 的 keepalive 时间. 好像是10分钟
		for event := range watcher.ResultChan() {
			switch event.Type {
			case watch.Modified:
				if h.IsReady(name) {
					watcher.Stop()
					return nil
				}
			case watch.Deleted: // deployment 已经删除了, 退出监听
				watcher.Stop()
				return fmt.Errorf("%s deleted", name)
			case watch.Bookmark:
				log.Debug("watch deployment: bookmark.")
			case watch.Error:
				log.Debug("watch deployment: error")
			}
		}
		// If event channel is closed, it means the server has closed the connection
		log.Debug("watch deployment: reconnect to kubernetes.")
		watcher.Stop()
	}
}

// GetRS get all replicaset created by the deployment.
func (h *Handler) GetRS(name string) ([]appsv1.ReplicaSet, error) {
	_, err := h.Get(name)
	if err != nil {
		return nil, err
	}

	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.LabelSelector = ""
	rsList, err := h.clientset.AppsV1().ReplicaSets(h.namespace).List(h.ctx, *listOptions)
	if err != nil {
		return nil, err
	}

	var rl []appsv1.ReplicaSet
	for _, r := range rsList.Items {
		for _, or := range r.OwnerReferences {
			if or.Name == name {
				rl = append(rl, r)
			}
		}
	}
	return rl, nil
}

// GetPods get all pods created by the deployment.
func (h *Handler) GetPods(name string) ([]corev1.Pod, error) {
	if _, err := h.Get(name); err != nil {
		return nil, err
	}

	rsList, err := h.GetRS(name)
	if err != nil {
		return nil, err
	}
	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.LabelSelector = ""
	podList, err := h.clientset.CoreV1().Pods(h.namespace).List(h.ctx, *listOptions)
	if err != nil {
		return nil, err
	}

	var pl []corev1.Pod
	for _, rs := range rsList {
		for _, pod := range podList.Items {
			for _, or := range pod.OwnerReferences {
				if or.Name == rs.Name {
					pl = append(pl, pod)
				}
			}
		}
	}
	return pl, nil
}

// GetPVC get deployment pvc list by name.
func (h *Handler) GetPVC(name string) ([]string, error) {
	deploy, err := h.Get(name)
	if err != nil {
		return nil, err
	}

	var pl []string
	for _, volume := range deploy.Spec.Template.Spec.Volumes {
		// 有些 volume.PersistentVolumeClaim 是不存在的, 其值默认是 nil 如果不加以判断就直接获取
		// volume.PersistentVolumeClaim.ClaimName, 就操作了非法地址, 程序会直接 panic.
		if volume.PersistentVolumeClaim != nil {
			pl = append(pl, volume.PersistentVolumeClaim.ClaimName)
		}
	}
	return pl, nil
}

// GetPV get deployment pv list by name.
func (h *Handler) GetPV(name string) ([]string, error) {
	if _, err := h.Get(name); err != nil {
		return nil, err
	}

	pvcList, err := h.GetPVC(name)
	if err != nil {
		return nil, err
	}

	// pvc.spec.volumeName 的值就是 pvc 中绑定的 pv 的名字
	var pl []string
	for _, pvc := range pvcList {
		pvcObj, err := h.clientset.CoreV1().
			PersistentVolumeClaims(h.namespace).Get(h.ctx, pvc, h.Options.GetOptions)
		if err == nil {
			pl = append(pl, pvcObj.Spec.VolumeName)
		}
	}
	return pl, nil
}

// GetAge get deployment age.
func (h *Handler) GetAge(name string) (time.Duration, error) {
	deploy, err := h.Get(name)
	if err != nil {
		return time.Duration(int64(0)), err
	}

	ctime := deploy.CreationTimestamp.Time
	return time.Now().Sub(ctime), nil
}
