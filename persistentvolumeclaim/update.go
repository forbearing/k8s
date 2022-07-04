package persistentvolumeclaim

import (
	"encoding/json"
	"io/ioutil"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// UpdateFromRaw update persistentvolumeclaim from map[string]interface{}.
func (h *Handler) UpdateFromRaw(raw map[string]interface{}) (*corev1.PersistentVolumeClaim, error) {
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

	return h.clientset.CoreV1().PersistentVolumeClaims(namespace).Update(h.ctx, pvc, h.Options.UpdateOptions)
}

// UpdateFromBytes update persistentvolumeclaim from bytes.
func (h *Handler) UpdateFromBytes(data []byte) (*corev1.PersistentVolumeClaim, error) {
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

	return h.clientset.CoreV1().PersistentVolumeClaims(namespace).Update(h.ctx, pvc, h.Options.UpdateOptions)
}

// UpdateFromFile update persistentvolumeclaim from yaml file.
func (h *Handler) UpdateFromFile(filename string) (*corev1.PersistentVolumeClaim, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.UpdateFromBytes(data)
}

// Update update persistentvolumeclaim from file, alias to "UpdateFromFile".
func (h *Handler) Update(filename string) (*corev1.PersistentVolumeClaim, error) {
	return h.UpdateFromFile(filename)
}
