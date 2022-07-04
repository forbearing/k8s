package pod

import (
	"encoding/json"
	"io/ioutil"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// UpdateFromRaw update pod from map[string]interface{}.
func (h *Handler) UpdateFromRaw(raw map[string]interface{}) (*corev1.Pod, error) {
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

	return h.clientset.CoreV1().Pods(namespace).Update(h.ctx, pod, h.Options.UpdateOptions)
}

// UpdateFromBytes update pod from bytes.
func (h *Handler) UpdateFromBytes(data []byte) (*corev1.Pod, error) {
	podJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	pod := &corev1.Pod{}
	err = json.Unmarshal(podJson, pod)
	if err != nil {
		return nil, err
	}

	var namespace string
	if len(pod.Namespace) != 0 {
		namespace = pod.Namespace
	} else {
		namespace = h.namespace
	}

	return h.clientset.CoreV1().Pods(namespace).Update(h.ctx, pod, h.Options.UpdateOptions)
}

// UpdateFromFile update pod from yaml file.
func (h *Handler) UpdateFromFile(filename string) (*corev1.Pod, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.UpdateFromBytes(data)
}

// Update update pod from file, alias to "UpdateFromFile".
func (h *Handler) Update(filename string) (*corev1.Pod, error) {
	return h.UpdateFromFile(filename)
}
