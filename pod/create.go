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

// Create creates pod from type string, []byte, *corev1.pod, corev1.pod,
// metav1.Object, runtime.Object, *unstructured.Unstructured, unstructured.Unstructured
// or map[string]interface{}.
func (h *Handler) Create(obj interface{}) (*corev1.Pod, error) {
	switch val := obj.(type) {
	case string:
		return h.CreateFromFile(val)
	case []byte:
		return h.CreateFromBytes(val)
	case *corev1.Pod:
		return h.CreateFromObject(val)
	case corev1.Pod:
		return h.CreateFromObject(&val)
	case *unstructured.Unstructured:
		return h.CreateFromUnstructured(val)
	case unstructured.Unstructured:
		return h.CreateFromUnstructured(&val)
	case map[string]interface{}:
		return h.CreateFromMap(val)
	case metav1.Object, runtime.Object:
		return h.CreateFromObject(val)
	default:
		return nil, ErrInvalidCreateType
	}
}

// CreateFromFile creates pod from yaml or json file.
func (h *Handler) CreateFromFile(filename string) (*corev1.Pod, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.CreateFromBytes(data)
}

// CreateFromBytes creates pod from bytes data.
func (h *Handler) CreateFromBytes(data []byte) (*corev1.Pod, error) {
	podJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	pod := &corev1.Pod{}
	if err = json.Unmarshal(podJson, pod); err != nil {
		return nil, err
	}
	return h.createPod(pod)
}

// CreateFromObject creates pod from metav1.Object or runtime.Object.
func (h *Handler) CreateFromObject(obj interface{}) (*corev1.Pod, error) {
	pod, ok := obj.(*corev1.Pod)
	if !ok {
		return nil, fmt.Errorf("object type is not *corev1.Pod")
	}
	return h.createPod(pod)
}

// CreateFromUnstructured creates pod from *unstructured.Unstructured.
func (h *Handler) CreateFromUnstructured(u *unstructured.Unstructured) (*corev1.Pod, error) {
	pod := &corev1.Pod{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), pod)
	if err != nil {
		return nil, err
	}
	return h.createPod(pod)
}

// CreateFromMap creates pod from map[string]interface{}.
func (h *Handler) CreateFromMap(u map[string]interface{}) (*corev1.Pod, error) {
	pod := &corev1.Pod{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, pod)
	if err != nil {
		return nil, err
	}
	return h.createPod(pod)
}

// createPod
func (h *Handler) createPod(pod *corev1.Pod) (*corev1.Pod, error) {
	namespace := pod.GetNamespace()
	if len(namespace) == 0 {
		namespace = h.namespace
	}
	pod.UID = ""
	pod.ResourceVersion = ""
	return h.clientset.CoreV1().Pods(namespace).Create(h.ctx, pod, h.Options.CreateOptions)
}
