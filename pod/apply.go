package pod

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

// Apply applies pod from type string, []byte, *corev1.pod, corev1.pod,
// metav1.Object, runtime.Object, *unstructured.Unstructured, unstructured.Unstructured
// or map[string]interface{}.
func (h *Handler) Apply(obj interface{}) (*corev1.Pod, error) {
	switch val := obj.(type) {
	case string:
		return h.ApplyFromFile(val)
	case []byte:
		return h.ApplyFromBytes(val)
	case *corev1.Pod:
		return h.ApplyFromObject(val)
	case corev1.Pod:
		return h.ApplyFromObject(&val)
	case *unstructured.Unstructured:
		return h.ApplyFromUnstructured(val)
	case unstructured.Unstructured:
		return h.ApplyFromUnstructured(&val)
	case map[string]interface{}:
		return h.ApplyFromMap(val)
	case metav1.Object, runtime.Object:
		return h.ApplyFromObject(val)
	default:
		return nil, ErrInvalidApplyType
	}
}

// ApplyFromFile applies pod from yaml or json file.
func (h *Handler) ApplyFromFile(filename string) (pod *corev1.Pod, err error) {
	pod, err = h.CreateFromFile(filename)
	if k8serrors.IsAlreadyExists(err) {
		pod, err = h.UpdateFromFile(filename)
	}
	return
}

// ApplyFromBytes applies pod from bytes data.
func (h *Handler) ApplyFromBytes(data []byte) (pod *corev1.Pod, err error) {
	pod, err = h.CreateFromBytes(data)
	if k8serrors.IsAlreadyExists(err) {
		pod, err = h.UpdateFromBytes(data)
	}
	return
}

// ApplyFromObject applies deployment from metav1.Object or runtime.Object.
func (h *Handler) ApplyFromObject(obj interface{}) (*corev1.Pod, error) {
	pod, ok := obj.(*corev1.Pod)
	if !ok {
		return nil, fmt.Errorf("object type is not *corev1.Pod")
	}
	return h.applyPod(pod)
}

// ApplyFromUnstructured applies pod from *unstructured.Unstructured.
func (h *Handler) ApplyFromUnstructured(u *unstructured.Unstructured) (*corev1.Pod, error) {
	pod := &corev1.Pod{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), pod)
	if err != nil {
		return nil, err
	}
	return h.applyPod(pod)
}

// ApplyFromMap applies pod from map[string]interface{}.
func (h *Handler) ApplyFromMap(u map[string]interface{}) (*corev1.Pod, error) {
	pod := &corev1.Pod{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, pod)
	if err != nil {
		return nil, err
	}
	return h.applyPod(pod)
}

// applyPod
func (h *Handler) applyPod(pod *corev1.Pod) (*corev1.Pod, error) {
	_, err := h.createPod(pod)
	if k8serrors.IsAlreadyExists(err) {
		return h.updatePod(pod)
	}
	return pod, err
}
