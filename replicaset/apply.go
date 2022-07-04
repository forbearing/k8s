package replicaset

import (
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
)

// ApplyFromRaw apply replicaset from map[string]interface{}.
func (h *Handler) ApplyFromRaw(raw map[string]interface{}) (*appsv1.ReplicaSet, error) {
	replicaset := &appsv1.ReplicaSet{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(raw, replicaset)
	if err != nil {
		return nil, err
	}

	var namespace string
	if len(replicaset.Namespace) != 0 {
		namespace = replicaset.Namespace
	} else {
		namespace = h.namespace
	}

	_, err = h.clientset.AppsV1().ReplicaSets(namespace).Create(h.ctx, replicaset, h.Options.CreateOptions)
	if k8serrors.IsAlreadyExists(err) {
		replicaset, err = h.clientset.AppsV1().ReplicaSets(namespace).Update(h.ctx, replicaset, h.Options.UpdateOptions)
	}
	return replicaset, err
}

// ApplyFromBytes apply replicaset from bytes.
func (h *Handler) ApplyFromBytes(data []byte) (replicaset *appsv1.ReplicaSet, err error) {
	replicaset, err = h.CreateFromBytes(data)
	if errors.IsAlreadyExists(err) {
		replicaset, err = h.UpdateFromBytes(data)
	}
	return
}

// ApplyFromFile apply replicaset from yaml file.
func (h *Handler) ApplyFromFile(filename string) (replicaset *appsv1.ReplicaSet, err error) {
	replicaset, err = h.CreateFromFile(filename)
	if errors.IsAlreadyExists(err) {
		replicaset, err = h.UpdateFromFile(filename)
	}
	return
}

// Apply apply replicaset from yaml file, alias to "ApplyFromFile".
func (h *Handler) Apply(filename string) (*appsv1.ReplicaSet, error) {
	return h.ApplyFromFile(filename)
}
