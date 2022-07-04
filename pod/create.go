package pod

import (
	"encoding/json"
	"io/ioutil"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// CreateFromRaw create pod from map[string]interface{}.
func (h *Handler) CreateFromRaw(raw map[string]interface{}) (*corev1.Pod, error) {
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

	return h.clientset.CoreV1().Pods(namespace).Create(h.ctx, pod, h.Options.CreateOptions)
}

// CreateFromBytes create pod from bytes.
func (h *Handler) CreateFromBytes(data []byte) (*corev1.Pod, error) {
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

	//unstructMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(pod)
	//if err != nil {
	//    log.Error("ToUnstructured error")
	//    log.Error(err)
	//} else {
	//    log.Infof("%#v", unstructMap)
	//}
	return h.clientset.CoreV1().Pods(namespace).Create(h.ctx, pod, h.Options.CreateOptions)
}

// CreateFromFile create pod from yaml file.
func (h *Handler) CreateFromFile(filename string) (*corev1.Pod, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.CreateFromBytes(data)
}

// Create create pod from file, alias to "CreateFromFile".
func (h *Handler) Create(filename string) (*corev1.Pod, error) {
	return h.CreateFromFile(filename)
}
