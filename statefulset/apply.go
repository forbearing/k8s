package statefulset

import (
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
)

// ApplyFromRaw apply statefulset from map[string]interface{}.
func (h *Handler) ApplyFromRaw(raw map[string]interface{}) (*appsv1.StatefulSet, error) {
	sts := &appsv1.StatefulSet{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(raw, sts)
	if err != nil {
		return nil, err
	}

	var namespace string
	if len(sts.Namespace) != 0 {
		namespace = sts.Namespace
	} else {
		namespace = h.namespace
	}

	sts, err = h.clientset.AppsV1().StatefulSets(namespace).Create(h.ctx, sts, h.Options.CreateOptions)
	if k8serrors.IsAlreadyExists(err) {
		sts, err = h.clientset.AppsV1().StatefulSets(namespace).Update(h.ctx, sts, h.Options.UpdateOptions)
	}
	return sts, err
}

// ApplyFromBytes apply statefulset from bytes.
func (h *Handler) ApplyFromBytes(data []byte) (statefulset *appsv1.StatefulSet, err error) {
	statefulset, err = h.CreateFromBytes(data)
	if errors.IsAlreadyExists(err) {
		statefulset, err = h.UpdateFromBytes(data)
	}
	return
}

// ApplyFromFile apply statefulset from yaml file.
func (h *Handler) ApplyFromFile(filename string) (statefulset *appsv1.StatefulSet, err error) {
	statefulset, err = h.CreateFromFile(filename)
	if errors.IsAlreadyExists(err) {
		statefulset, err = h.UpdateFromFile(filename)
	}
	return
}

// Apply apply statefulset from file, alias to "ApplyFromFile".
func (h *Handler) Apply(filename string) (*appsv1.StatefulSet, error) {
	return h.ApplyFromFile(filename)
}
