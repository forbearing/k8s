package pod

import (
	"encoding/json"
	"io/ioutil"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// GetFromBytes get pod from bytes.
func (h *Handler) GetFromBytes(data []byte) (*corev1.Pod, error) {
	podJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	pod := &corev1.Pod{}
	if err = json.Unmarshal(podJson, pod); err != nil {
		return nil, err
	}

	var namespace string
	if len(pod.Namespace) != 0 {
		namespace = pod.Namespace
	} else {
		namespace = h.namespace
	}

	return h.WithNamespace(namespace).GetByName(pod.Name)
}

// GetFromFile get pod from yaml file.
func (h *Handler) GetFromFile(filename string) (*corev1.Pod, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.GetFromBytes(data)
}

// GetByName get pod by name.
func (h *Handler) GetByName(name string) (*corev1.Pod, error) {
	return h.clientset.CoreV1().Pods(h.namespace).Get(h.ctx, name, h.Options.GetOptions)
}

// Get get pod by name, alias to "GetByName".
func (h *Handler) Get(name string) (pod *corev1.Pod, err error) {
	return h.GetByName(name)
}
