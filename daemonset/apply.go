package daemonset

import (
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
)

// ApplyFromRaw apply daemonset from map[string]interface{}.
func (h *Handler) ApplyFromRaw(raw map[string]interface{}) (*appsv1.DaemonSet, error) {
	daemonset := &appsv1.DaemonSet{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(raw, daemonset)
	if err != nil {
		return nil, err
	}

	var namespace string
	if len(daemonset.Namespace) != 0 {
		namespace = daemonset.Namespace
	} else {
		namespace = h.namespace
	}

	_, err = h.clientset.AppsV1().DaemonSets(namespace).Create(h.ctx, daemonset, h.Options.CreateOptions)
	if k8serrors.IsAlreadyExists(err) {
		daemonset, err = h.clientset.AppsV1().DaemonSets(namespace).Update(h.ctx, daemonset, h.Options.UpdateOptions)
	}
	return daemonset, err
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
