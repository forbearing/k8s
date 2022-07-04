package pod

import (
	"encoding/json"
	"io/ioutil"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// DeleteFromBytes delete pod from bytes.
func (h *Handler) DeleteFromBytes(data []byte) error {
	podJson, err := yaml.ToJSON(data)
	if err != nil {
		return err
	}

	pod := &corev1.Pod{}
	err = json.Unmarshal(podJson, pod)
	if err != nil {
		return err
	}

	var namespace string
	if len(pod.Namespace) != 0 {
		namespace = pod.Namespace
	} else {
		namespace = h.namespace
	}

	return h.WithNamespace(namespace).DeleteByName(pod.Name)
}

// DeleteFromFile delete pod from yaml file.
func (h *Handler) DeleteFromFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return h.DeleteFromBytes(data)
}

// DeleteByName delete pod by name.
func (h *Handler) DeleteByName(name string) error {
	return h.clientset.CoreV1().Pods(h.namespace).Delete(h.ctx, name, h.Options.DeleteOptions)
}

// Delete delete pod by name, alias to "DeleteByName".
func (h *Handler) Delete(name string) error {
	return h.DeleteByName(name)
}
