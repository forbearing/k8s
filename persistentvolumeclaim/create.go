package persistentvolumeclaim

import (
	"encoding/json"
	"io/ioutil"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// CreateFromRaw create persistentvolumeclaim from map[string]interface{}.
func (h *Handler) CreateFromRaw(raw map[string]interface{}) (*corev1.PersistentVolumeClaim, error) {
	pvc := &corev1.PersistentVolumeClaim{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(raw, pvc)
	if err != nil {
		return nil, err
	}

	var namespace string
	if len(pvc.Namespace) != 0 {
		namespace = pvc.Namespace
	} else {
		namespace = h.namespace
	}

	return h.clientset.CoreV1().PersistentVolumeClaims(namespace).Create(h.ctx, pvc, h.Options.CreateOptions)
}

// CreateFromBytes create persistentvolumeclaim from bytes.
func (h *Handler) CreateFromBytes(data []byte) (*corev1.PersistentVolumeClaim, error) {
	pvcJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	pvc := &corev1.PersistentVolumeClaim{}
	err = json.Unmarshal(pvcJson, pvc)
	if err != nil {
		return nil, err
	}

	var namespace string
	if len(pvc.Namespace) != 0 {
		namespace = pvc.Namespace
	} else {
		namespace = h.namespace
	}

	return h.clientset.CoreV1().PersistentVolumeClaims(namespace).Create(h.ctx, pvc, h.Options.CreateOptions)
}

// CreateFromFile create persistentvolumeclaim from yaml file.
func (h *Handler) CreateFromFile(filename string) (*corev1.PersistentVolumeClaim, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.CreateFromBytes(data)
}

// Create create persistentvolumeclaim from yaml file, alias to "CreateFromFile".
func (h *Handler) Create(filename string) (*corev1.PersistentVolumeClaim, error) {
	return h.CreateFromFile(filename)
}
