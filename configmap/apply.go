package configmap

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
)

// ApplyFromRaw apply configmap from map[string]interface{}.
func (h *Handler) ApplyFromRaw(raw map[string]interface{}) (*corev1.ConfigMap, error) {
	configmap := &corev1.ConfigMap{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(raw, configmap)
	if err != nil {
		return nil, err
	}

	var namespace string
	if len(configmap.Namespace) != 0 {
		namespace = configmap.Namespace
	} else {
		namespace = h.namespace
	}

	configmap, err = h.clientset.CoreV1().ConfigMaps(namespace).Create(h.ctx, configmap, h.Options.CreateOptions)
	if k8serrors.IsAlreadyExists(err) {
		configmap, err = h.clientset.CoreV1().ConfigMaps(namespace).Update(h.ctx, configmap, h.Options.UpdateOptions)
	}
	return configmap, err
}

// ApplyFromBytes apply configmap from bytes.
func (h *Handler) ApplyFromBytes(data []byte) (configmap *corev1.ConfigMap, err error) {
	configmap, err = h.CreateFromBytes(data)
	if errors.IsAlreadyExists(err) {
		configmap, err = h.UpdateFromBytes(data)
	}
	return
}

// ApplyFromFile apply configmap from yaml file.
func (h *Handler) ApplyFromFile(filename string) (configmap *corev1.ConfigMap, err error) {
	configmap, err = h.CreateFromFile(filename)
	if errors.IsAlreadyExists(err) {
		configmap, err = h.UpdateFromFile(filename)
	}
	return
}

// Apply apply configmap from yaml file, alias to "ApplyFromFile".
func (h *Handler) Apply(filename string) (*corev1.ConfigMap, error) {
	return h.ApplyFromFile(filename)
}
