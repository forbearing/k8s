package service

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
)

// ApplyFromRaw apply service from map[string]interface{}.
func (h *Handler) ApplyFromRaw(raw map[string]interface{}) (*corev1.Service, error) {
	service := &corev1.Service{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(raw, service)
	if err != nil {
		return nil, err
	}

	var namespace string
	if len(service.Namespace) != 0 {
		namespace = service.Namespace
	} else {
		namespace = h.namespace
	}

	service, err = h.clientset.CoreV1().Services(namespace).Create(h.ctx, service, h.Options.CreateOptions)
	if k8serrors.IsAlreadyExists(err) {
		service, err = h.clientset.CoreV1().Services(namespace).Update(h.ctx, service, h.Options.UpdateOptions)
	}
	return service, err
}

// ApplyFromBytes apply service from bytes.
func (h *Handler) ApplyFromBytes(data []byte) (service *corev1.Service, err error) {
	service, err = h.CreateFromBytes(data)
	if errors.IsAlreadyExists(err) {
		service, err = h.UpdateFromBytes(data)
	}
	return
}

// ApplyFromFile apply service from yaml file.
func (h *Handler) ApplyFromFile(filename string) (service *corev1.Service, err error) {
	service, err = h.CreateFromFile(filename)
	if errors.IsAlreadyExists(err) {
		service, err = h.UpdateFromFile(filename)
	}
	return
}

// Apply apply service from yaml file, alias to "ApplyFromFile".
func (h *Handler) Apply(filename string) (*corev1.Service, error) {
	return h.ApplyFromFile(filename)
}
