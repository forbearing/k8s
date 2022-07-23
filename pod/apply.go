package pod

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
)

// Apply applies pod from type string, []byte, *corev1.pod, corev1.pod,
// runtime.Object or map[string]interface{}.
func (h *Handler) Apply(obj interface{}) (*corev1.Pod, error) {
	switch val := obj.(type) {
	case string:
		return h.ApplyFromFile(val)
	case []byte:
		return h.ApplyFromBytes(val)
	case *corev1.Pod:
		return h.CreateFromObject(val)
	case corev1.Pod:
		return h.CreateFromObject(&val)
	case map[string]interface{}:
		return h.CreateFromUnstructured(val)
	default:
		return nil, ERR_TYPE_APPLY
	}
}

// ApplyFromFile applies pod from yaml file.
func (h *Handler) ApplyFromFile(filename string) (pod *corev1.Pod, err error) {
	pod, err = h.CreateFromFile(filename)
	if k8serrors.IsAlreadyExists(err) {
		pod, err = h.UpdateFromFile(filename)
	}
	return
}

// ApplyFromBytes applies pod from bytes.
func (h *Handler) ApplyFromBytes(data []byte) (pod *corev1.Pod, err error) {
	pod, err = h.CreateFromBytes(data)
	if k8serrors.IsAlreadyExists(err) {
		log.Debug(err)
		pod, err = h.UpdateFromBytes(data)
	}
	return
}

// ApplyFromObject applies deployment from runtime.Object.
func (h *Handler) ApplyFromObject(obj runtime.Object) (*corev1.Pod, error) {
	pod, ok := obj.(*corev1.Pod)
	if !ok {
		return nil, fmt.Errorf("object is not *corev1.Pod")
	}
	return h.applyPod(pod)
}

// ApplyFromUnstructured applies pod from map[string]interface{}.
func (h *Handler) ApplyFromUnstructured(u map[string]interface{}) (*corev1.Pod, error) {
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
