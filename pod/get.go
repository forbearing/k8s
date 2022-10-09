package pod

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// Get gets pod from type string, []byte, *corev1.Pod,
// corev1.Pod, metav1.Object, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.
//
// If passed parameter type is string, it will simply call GetByName instead of GetFromFile.
// You should always explicitly call GetFromFile to get a pod from file path.
func (h *Handler) Get(obj interface{}) (*corev1.Pod, error) {
	switch val := obj.(type) {
	case string:
		return h.GetByName(val)
	case []byte:
		return h.GetFromBytes(val)
	case *corev1.Pod:
		return h.GetFromObject(val)
	case corev1.Pod:
		return h.GetFromObject(&val)
	case *unstructured.Unstructured:
		return h.GetFromUnstructured(val)
	case unstructured.Unstructured:
		return h.GetFromUnstructured(&val)
	case map[string]interface{}:
		return h.GetFromMap(val)
	case metav1.Object, runtime.Object:
		return h.GetFromObject(val)
	default:
		return nil, ErrInvalidGetType
	}
}

// GetByName gets pod by name.
func (h *Handler) GetByName(name string) (*corev1.Pod, error) {
	return h.clientset.CoreV1().Pods(h.namespace).Get(h.ctx, name, h.Options.GetOptions)
}

// GetFromFile gets pod from yaml or json file.
func (h *Handler) GetFromFile(filename string) (*corev1.Pod, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.GetFromBytes(data)
}

// GetFromBytes gets pod from bytes data.
func (h *Handler) GetFromBytes(data []byte) (*corev1.Pod, error) {
	podJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	pod := &corev1.Pod{}
	if err = json.Unmarshal(podJson, pod); err != nil {
		return nil, err
	}
	return h.getPod(pod)
}

// GetFromObject gets pod from metav1.Object or runtime.Object.
func (h *Handler) GetFromObject(obj interface{}) (*corev1.Pod, error) {
	pod, ok := obj.(*corev1.Pod)
	if !ok {
		return nil, fmt.Errorf("object type is not *corev1.Pod")
	}
	return h.getPod(pod)
}

// GetFromUnstructured gets pod from *unstructured.Unstructured.
func (h *Handler) GetFromUnstructured(u *unstructured.Unstructured) (*corev1.Pod, error) {
	pod := &corev1.Pod{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), pod)
	if err != nil {
		return nil, err
	}
	return h.getPod(pod)
}

// GetFromMap gets pod from map[string]interface{}.
func (h *Handler) GetFromMap(u map[string]interface{}) (*corev1.Pod, error) {
	pod := &corev1.Pod{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, pod)
	if err != nil {
		return nil, err
	}
	return h.getPod(pod)
}

// getPod
// It's necessary to gets a new pod resource from a old pod resource,
// because old pod usually don't have pod.Status field.
func (h *Handler) getPod(pod *corev1.Pod) (*corev1.Pod, error) {
	namespace := pod.GetNamespace()
	if len(namespace) == 0 {
		namespace = h.namespace
	}
	return h.clientset.CoreV1().Pods(namespace).Get(h.ctx, pod.Name, h.Options.GetOptions)
}
