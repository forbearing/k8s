package pod

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// Update updates pod from type string, []byte, *corev1.pod, corev1.pod,
// runtime.Object, *unstructured.Unstructured, unstructured.Unstructured
// or map[string]interface{}.
func (h *Handler) Update(obj interface{}) (*corev1.Pod, error) {
	switch val := obj.(type) {
	case string:
		return h.UpdateFromFile(val)
	case []byte:
		return h.UpdateFromBytes(val)
	case *corev1.Pod:
		return h.UpdateFromObject(val)
	case corev1.Pod:
		return h.UpdateFromObject(&val)
	case *unstructured.Unstructured:
		return h.UpdateFromUnstructured(val)
	case unstructured.Unstructured:
		return h.UpdateFromUnstructured(&val)
	case map[string]interface{}:
		return h.UpdateFromMap(val)
	default:
		return nil, ERR_TYPE_UPDATE
	}
}

// UpdateFromFile updates pod from yaml file.
func (h *Handler) UpdateFromFile(filename string) (*corev1.Pod, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.UpdateFromBytes(data)
}

// UpdateFromBytes updates pod from bytes.
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
	return h.updatePod(pod)
}

// UpdateFromObject updates pod from runtime.Object.
func (h *Handler) UpdateFromObject(obj runtime.Object) (*corev1.Pod, error) {
	pod, ok := obj.(*corev1.Pod)
	if !ok {
		return nil, fmt.Errorf("object type is not *corev1.Pod")
	}
	return h.updatePod(pod)
}

// UpdateFromUnstructured updates pod from *unstructured.Unstructured.
func (h *Handler) UpdateFromUnstructured(u *unstructured.Unstructured) (*corev1.Pod, error) {
	pod := &corev1.Pod{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), pod)
	if err != nil {
		return nil, err
	}
	return h.updatePod(pod)
}

// UpdateFromMap updates pod from map[string]interface{}.
func (h *Handler) UpdateFromMap(u map[string]interface{}) (*corev1.Pod, error) {
	pod := &corev1.Pod{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, pod)
	if err != nil {
		return nil, err
	}
	return h.updatePod(pod)
}

// updatePod
func (h *Handler) updatePod(pod *corev1.Pod) (*corev1.Pod, error) {
	var namespace string
	if len(pod.Namespace) != 0 {
		namespace = pod.Namespace
	} else {
		namespace = h.namespace
	}
	pod.UID = ""
	pod.ResourceVersion = ""
	return h.clientset.CoreV1().Pods(namespace).Update(h.ctx, pod, h.Options.UpdateOptions)
}
