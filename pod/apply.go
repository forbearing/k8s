package pod

import (
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
)

// ApplyFromRaw apply pod from map[string]interface{}.
func (h *Handler) ApplyFromRaw(raw map[string]interface{}) (*corev1.Pod, error) {
	pod := &corev1.Pod{}
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

	pod, err = h.clientset.CoreV1().Pods(namespace).Create(h.ctx, pod, h.Options.CreateOptions)
	if k8serrors.IsAlreadyExists(err) {
		pod, err = h.clientset.CoreV1().Pods(namespace).Update(h.ctx, pod, h.Options.UpdateOptions)
	}
	return pod, err
}

// ApplyFromBytes apply pod from bytes.
func (h *Handler) ApplyFromBytes(data []byte) (pod *corev1.Pod, err error) {
	pod, err = h.CreateFromBytes(data)
	if k8serrors.IsAlreadyExists(err) {
		log.Debug(err)
		pod, err = h.UpdateFromBytes(data)
	}
	return
}

// ApplyFromFile apply pod from yaml file.
func (h *Handler) ApplyFromFile(filename string) (pod *corev1.Pod, err error) {
	pod, err = h.CreateFromFile(filename)
	if k8serrors.IsAlreadyExists(err) {
		pod, err = h.UpdateFromFile(filename)
	}
	return
}

// Apply apply pod from file, alias to "ApplyFromFile".
func (h *Handler) Apply(filename string) (*corev1.Pod, error) {
	return h.ApplyFromFile(filename)
}
