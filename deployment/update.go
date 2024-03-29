package deployment

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// Update updates deployment from type string, []byte, *appsv1.Deployment,
// appsv1.Deployment, metav1.Object, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.
func (h *Handler) Update(obj interface{}) (*appsv1.Deployment, error) {
	switch val := obj.(type) {
	case string:
		return h.UpdateFromFile(val)
	case []byte:
		return h.UpdateFromBytes(val)
	case *appsv1.Deployment:
		return h.UpdateFromObject(val)
	case appsv1.Deployment:
		return h.UpdateFromObject(&val)
	case *unstructured.Unstructured:
		return h.UpdateFromUnstructured(val)
	case unstructured.Unstructured:
		return h.UpdateFromUnstructured(&val)
	case map[string]interface{}:
		return h.UpdateFromMap(val)
	case metav1.Object, runtime.Object:
		//if reflect.TypeOf(val).String() == "*unstructured.Unstructured" {
		//    return h.UpdateFromUnstructured(val.(*unstructured.Unstructured))
		//}
		return h.UpdateFromObject(val)
	default:
		return nil, ErrInvalidUpdateType
	}
}

// UpdateFromFile updates deployment from yaml or json file.
func (h *Handler) UpdateFromFile(filename string) (*appsv1.Deployment, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.UpdateFromBytes(data)
}

// UpdateFromBytes updates deployment from bytes data.
func (h *Handler) UpdateFromBytes(data []byte) (*appsv1.Deployment, error) {
	deployJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	deploy := &appsv1.Deployment{}
	if err = json.Unmarshal(deployJson, deploy); err != nil {
		return nil, err
	}
	return h.updateDeployment(deploy)
}

// UpdateFromObject updates deployment from metav1.Object or runtime.Object.
func (h *Handler) UpdateFromObject(obj interface{}) (*appsv1.Deployment, error) {
	deploy, ok := obj.(*appsv1.Deployment)
	if !ok {
		return nil, fmt.Errorf("object type is not *appsv1.Deployment")
	}
	return h.updateDeployment(deploy)
}

// UpdateFromUnstructured updates deployment from *unstructured.Unstructured.
func (h *Handler) UpdateFromUnstructured(u *unstructured.Unstructured) (*appsv1.Deployment, error) {
	deploy := &appsv1.Deployment{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), deploy)
	if err != nil {
		return nil, err
	}
	return h.updateDeployment(deploy)
}

// UpdateFromMap updates deployment from map[string]interface{}.
func (h *Handler) UpdateFromMap(u map[string]interface{}) (*appsv1.Deployment, error) {
	deploy := &appsv1.Deployment{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, deploy)
	if err != nil {
		return nil, err
	}
	return h.updateDeployment(deploy)
}

// updateDeployment
func (h *Handler) updateDeployment(deploy *appsv1.Deployment) (*appsv1.Deployment, error) {
	namespace := deploy.GetNamespace()
	if len(namespace) == 0 {
		namespace = h.namespace
	}
	// resourceVersion cann't be set, the resourceVersion field is empty.
	deploy.ResourceVersion = ""
	deploy.UID = ""
	return h.clientset.AppsV1().Deployments(namespace).Update(h.ctx, deploy, h.Options.UpdateOptions)
}
