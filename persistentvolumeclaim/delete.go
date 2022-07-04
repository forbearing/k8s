package persistentvolumeclaim

import (
	"encoding/json"
	"io/ioutil"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// DeleteFromBytes delete persistentvolumeclaim from bytes.
func (h *Handler) DeleteFromBytes(data []byte) error {
	pvcJson, err := yaml.ToJSON(data)
	if err != nil {
		return err
	}

	pvc := &corev1.PersistentVolumeClaim{}
	if err = json.Unmarshal(pvcJson, pvc); err != nil {
		return err
	}

	var namespace string
	if len(pvc.Namespace) != 0 {
		namespace = pvc.Namespace
	} else {
		namespace = h.namespace
	}

	return h.WithNamespace(namespace).DeleteByName(pvc.Name)
}

// DeleteFromFile delete persistentvolumeclaim from yaml file.
func (h *Handler) DeleteFromFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return h.DeleteFromBytes(data)
}

// DeleteByName delete persistentvolumeclaim by name.
func (h *Handler) DeleteByName(name string) error {
	return h.clientset.CoreV1().PersistentVolumeClaims(h.namespace).Delete(h.ctx, name, h.Options.DeleteOptions)
}

// Delete delete persistentvolumeclaim by name, alias to "DeleteByName".
func (h *Handler) Delete(name string) error {
	return h.DeleteByName(name)
}
