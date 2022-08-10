package pod

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
)

/*
ref:
    https://hoelz.ro/blog/which-fields-can-you-use-with-kubernetes-field-selectors
    https://kubernetes.io/docs/concepts/overview/working-with-objects/field-selectors/
*/

// List list all pods in k8s cluster where the pod is running, it simply call `ListAll`.
func (h *Handler) List() ([]*corev1.Pod, error) {
	return h.ListAll()
}

// ListByLabel list pods by labels.
// Multiple labels separated by comma(",") eg: "name=myapp,role=devops",
// and there is an "And" relationship between multiple labels.
func (h *Handler) ListByLabel(labels string) ([]*corev1.Pod, error) {
	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.LabelSelector = labels
	//listOptions.ResourceVersion = ""
	podList, err := h.clientset.CoreV1().Pods(h.namespace).List(h.ctx, *listOptions)
	if err != nil {
		return nil, err
	}
	return extractList(podList), nil
}

// ListByField list pods by field, work like `kubectl get xxx --field-selector=xxx`.
func (h *Handler) ListByField(field string) ([]*corev1.Pod, error) {
	// ParseSelector takes a string representing a selector and returns an
	// object suitable for matching, or an error.
	fieldSelector, err := fields.ParseSelector(field)
	if err != nil {
		return nil, err
	}
	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.FieldSelector = fieldSelector.String()

	podList, err := h.clientset.CoreV1().Pods(h.namespace).List(h.ctx, *listOptions)
	if err != nil {
		return nil, err
	}
	return extractList(podList), nil
}

// ListByNamespace list all pods in the specified namespace.
func (h *Handler) ListByNamespace(namespace string) ([]*corev1.Pod, error) {
	return h.WithNamespace(namespace).ListByLabel("")
}

// ListAll list all pods in the k8s cluster where the pod is running.
func (h *Handler) ListAll() ([]*corev1.Pod, error) {
	return h.WithNamespace(metav1.NamespaceAll).ListByLabel("")
}

// ListByNode list all pods in the k8s node where the pod is running.
func (h *Handler) ListByNode(name string) ([]*corev1.Pod, error) {
	field := fmt.Sprintf("spec.nodeName=%s", name)
	return h.WithNamespace(metav1.NamespaceAll).ListByField(field)
}

// ListRunning list all pods whose .status.phase is "Running" in the k8s cluster.
// "Running" means the pod has been bound to a node and all of the containers have been started.
// At least one container is still running or is in the process of being restarted.
func (h *Handler) ListRunning() ([]*corev1.Pod, error) {
	field := "status.phase=Running"
	return h.WithNamespace(metav1.NamespaceAll).ListByField(field)
}

// ListSucceeded list all pods whose .status.phase is "Succeeded" in the k8s cluster.
// "Succeeded" means that all containers in the pod have voluntarily terminated
// with a container exit code of 0, and the system is not going to restart
// any of these containers.
func (h *Handler) ListSucceeded() ([]*corev1.Pod, error) {
	field := "status.phase=Succeeded"
	return h.WithNamespace(metav1.NamespaceAll).ListByField(field)
}

// ListFailed list all pods whose .status.phase is "Failed" in the k8s cluster.
// "Failed" means that all containers in the pod have terminated, and at least
// one container has terminated in a failure (exited with a non-zero exit code
// or was stopped by the system).
func (h *Handler) ListFailed() ([]*corev1.Pod, error) {
	field := "status.phase=Failed"
	return h.WithNamespace(metav1.NamespaceAll).ListByField(field)
}

// ListPending list all pods whose .status.phase is "Pending" in the k8s cluster.
// "Pending" means the pod has been accepted by the system, but one or more of
// the containers has not been started. This includes time before being bound to
// a node, as well as time spent pulling images onto the host.
func (h *Handler) ListPending() ([]*corev1.Pod, error) {
	field := "status.phase=Pending"
	return h.WithNamespace(metav1.NamespaceAll).ListByField(field)
}

// ListUnknow list all pods whose .status.phase is "Unknow" in the k8s cluster.
// "Unknown" means that for some reason the state of the pod could not be obtained,
// typically due to an error in communicating with the host of the pod.
// Deprecated: It isn't being set since 2015 (74da3b14b0c0f658b3bb8d2def5094686d0e9095)
func (h *Handler) ListUnknow() ([]*corev1.Pod, error) {
	field := "status.phase=Unknown"
	return h.WithNamespace(metav1.NamespaceAll).ListByField(field)
}

// extractList
func extractList(podList *corev1.PodList) []*corev1.Pod {
	//var pl []*corev1.Pod
	// 不能这种方式遍历
	// 参考: https://segmentfault.com/a/1190000038352530
	//for _, p := range podList.Items {
	//    pl = append(pl, &p)
	//}

	var pl []*corev1.Pod
	for i := range podList.Items {
		//val := ([]corev1.Pod(podList.Items))[i]
		//pl = append(pl, &val)
		pl = append(pl, &podList.Items[i])
	}
	return pl
}

// ListByStatus, reverse=true/false
// https://kubernetes.io/docs/concepts/overview/working-with-objects/field-selectors/

// ListByLabel(label... string)
// ListByNamespace(namespace... string)

// ListBy

// https://github.com/kubernetes-sigs/kubebuilder/blob/7cd3532662567e0a7568415e271f0b29cece202c/docs/book/src/cronjob-tutorial/testdata/project/controllers/cronjob_controller.go#L130
// kubebuilder
