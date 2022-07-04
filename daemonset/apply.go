package daemonset

import (
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
)

// ApplyFromRaw apply daemonset from map[string]interface{}.
func (h *Handler) ApplyFromRaw(raw map[string]interface{}) (*appsv1.DaemonSet, error) {
	pod := &appsv1.DaemonSet{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(raw, pod)
	if err != nil {
		return nil, err
	}

	var namespace string
	if len(pod.Namespace) != 0 {
		namespace = pod.Namespace
	} else {
		namespace = h.namespace
	}

	pod, err = h.clientset.AppsV1().DaemonSets(namespace).Create(h.ctx, pod, h.Options.CreateOptions)
	if k8serrors.IsAlreadyExists(err) {
		pod, err = h.clientset.AppsV1().DaemonSets(namespace).Update(h.ctx, pod, h.Options.UpdateOptions)
	}
	return pod, err
}

// ApplyFromBytes apply daemonset from bytes.
func (h *Handler) ApplyFromBytes(data []byte) (daemonset *appsv1.DaemonSet, err error) {
	daemonset, err = h.CreateFromBytes(data)
	if errors.IsAlreadyExists(err) {
		daemonset, err = h.UpdateFromBytes(data)
	}
	return
}

// ApplyFromFile apply daemonset from yaml file.
func (h *Handler) ApplyFromFile(filename string) (daemonset *appsv1.DaemonSet, err error) {
	daemonset, err = h.CreateFromFile(filename)
	if errors.IsAlreadyExists(err) {
		daemonset, err = h.UpdateFromFile(filename)
	}
	return
}

// Apply apply daemonset from yaml file, alias to "ApplyFromFile".
func (h *Handler) Apply(filename string) (*appsv1.DaemonSet, error) {
	return h.ApplyFromFile(filename)
}
