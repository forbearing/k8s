package namespace

import (
	"fmt"
	"reflect"

	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

// Apply applies namespace from type string, []byte, *corev1.Namespace,
// corev1.Namespace, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.
func (h *Handler) Apply(obj interface{}) (*corev1.Namespace, error) {
	switch val := obj.(type) {
	case string:
		return h.ApplyFromFile(val)
	case []byte:
		return h.ApplyFromBytes(val)
	case *corev1.Namespace:
		return h.ApplyFromObject(val)
	case corev1.Namespace:
		return h.ApplyFromObject(&val)
	case runtime.Object:
		if reflect.TypeOf(val).String() == "*unstructured.Unstructured" {
			return h.ApplyFromUnstructured(val.(*unstructured.Unstructured))
		}
		return h.ApplyFromObject(val)
	case *unstructured.Unstructured:
		return h.ApplyFromUnstructured(val)
	case unstructured.Unstructured:
		return h.ApplyFromUnstructured(&val)
	case map[string]interface{}:
		return h.ApplyFromMap(val)
	default:
		return nil, ErrInvalidApplyType
	}
}

// ApplyFromFile applies namespace from yaml file.
func (h *Handler) ApplyFromFile(filename string) (ns *corev1.Namespace, err error) {
	ns, err = h.CreateFromFile(filename)
	if k8serrors.IsAlreadyExists(err) { // if namespace already exist, update it.
		ns, err = h.UpdateFromFile(filename)
	}
	return
}

// ApplyFromBytes pply namespace from bytes.
func (h *Handler) ApplyFromBytes(data []byte) (ns *corev1.Namespace, err error) {
	ns, err = h.CreateFromBytes(data)
	if k8serrors.IsAlreadyExists(err) {
		ns, err = h.UpdateFromBytes(data)
	}
	return
}

// ApplyFromObject applies namespace from runtime.Object.
func (h *Handler) ApplyFromObject(obj runtime.Object) (*corev1.Namespace, error) {
	ns, ok := obj.(*corev1.Namespace)
	if !ok {
		return nil, fmt.Errorf("object type is not *corev1.Namespace")
	}
	return h.applyNamespace(ns)
}

// ApplyFromUnstructured applies namespace from *unstructured.Unstructured.
func (h *Handler) ApplyFromUnstructured(u *unstructured.Unstructured) (*corev1.Namespace, error) {
	ns := &corev1.Namespace{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), ns)
	if err != nil {
		return nil, err
	}
	return h.applyNamespace(ns)
}

// ApplyFromMap applies namespace from map[string]interface{}.
func (h *Handler) ApplyFromMap(u map[string]interface{}) (*corev1.Namespace, error) {
	ns := &corev1.Namespace{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, ns)
	if err != nil {
		return nil, err
	}
	return h.applyNamespace(ns)
}

// applyNamespace
func (h *Handler) applyNamespace(ns *corev1.Namespace) (*corev1.Namespace, error) {
	_, err := h.createNamespace(ns)
	if k8serrors.IsAlreadyExists(err) {
		return h.updateNamespace(ns)
	}
	return ns, err
}
