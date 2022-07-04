package namespace

import (
	"encoding/json"
	"io/ioutil"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// UpdateFromRaw update namespace from map[string]interface{}.
func (h *Handler) UpdateFromRaw(raw map[string]interface{}) (*corev1.Namespace, error) {
	namespace := &corev1.Namespace{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(raw, namespace)
	if err != nil {
		return nil, err
	}

	return h.clientset.CoreV1().Namespaces().Update(h.ctx, namespace, h.Options.UpdateOptions)
}

// UpdateFromBytes update namespace from bytes.
func (h *Handler) UpdateFromBytes(data []byte) (*corev1.Namespace, error) {
	nsJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	ns := &corev1.Namespace{}
	err = json.Unmarshal(nsJson, ns)
	if err != nil {
		return nil, err
	}

	return h.clientset.CoreV1().Namespaces().Update(h.ctx, ns, h.Options.UpdateOptions)
}

// UpdateFromFile update namespace from yaml file.
func (h *Handler) UpdateFromFile(filename string) (*corev1.Namespace, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.UpdateFromBytes(data)
}

// Update update namespace from yaml file, alias to "UpdateFromFile".
func (h *Handler) Update(filename string) (*corev1.Namespace, error) {
	return h.UpdateFromFile(filename)
}
