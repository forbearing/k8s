package persistentvolumeclaim

import (
	"encoding/json"
	"io/ioutil"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// GetFromBytes get persistentvolumeclaim from bytes.
func (h *Handler) GetFromBytes(data []byte) (*corev1.PersistentVolumeClaim, error) {
	pvcJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	pvc := &corev1.PersistentVolumeClaim{}
	if err = json.Unmarshal(pvcJson, pvc); err != nil {
		return nil, err
	}

	var namespace string
	if len(pvc.Namespace) != 0 {
		namespace = pvc.Namespace
	} else {
		namespace = h.namespace
	}

	return h.WithNamespace(namespace).GetByName(pvc.Name)
}

// GetFromFile get persistentvolumeclaim from yaml file.
func (h *Handler) GetFromFile(filename string) (*corev1.PersistentVolumeClaim, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.GetFromBytes(data)
}

// GetByName get persistentvolumeclaim by name.
func (h *Handler) GetByName(name string) (*corev1.PersistentVolumeClaim, error) {
	return h.clientset.CoreV1().PersistentVolumeClaims(h.namespace).Get(h.ctx, name, h.Options.GetOptions)
}

// Get get persistentvolumeclaim by name, alias to "GetByName".
func (h *Handler) Get(name string) (*corev1.PersistentVolumeClaim, error) {
	return h.GetByName(name)
}
