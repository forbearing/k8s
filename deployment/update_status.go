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

// UpdateStatus updates deployment from type string, []byte, *appsv1.Deployment,
// appsv1.Deployment, metav1.Object, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.
func (h *Handler) UpdateStatus(obj interface{}) (*appsv1.Deployment, error) {
	switch val := obj.(type) {
	case string:
		return h.UpdateStatusFromFile(val)
	case []byte:
		return h.UpdateStatusFromBytes(val)
	case *appsv1.Deployment:
		return h.UpdateStatusFromObject(val)
	case appsv1.Deployment:
		return h.UpdateStatusFromObject(&val)
	case *unstructured.Unstructured:
		return h.UpdateStatusFromUnstructured(val)
	case unstructured.Unstructured:
		return h.UpdateStatusFromUnstructured(&val)
	case map[string]interface{}:
		return h.UpdateStatusFromMap(val)
	case metav1.Object, runtime.Object:
		return h.UpdateStatusFromObject(val)
	default:
		return nil, ErrInvalidUpdateType
	}
}

// UpdateStatusFromFile updates deployment from yaml or json file.
func (h *Handler) UpdateStatusFromFile(filename string) (*appsv1.Deployment, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.UpdateStatusFromBytes(data)
}

// UpdateStatusFromBytes updates deployment from bytes data.
func (h *Handler) UpdateStatusFromBytes(data []byte) (*appsv1.Deployment, error) {
	deployJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	deploy := &appsv1.Deployment{}
	if err = json.Unmarshal(deployJson, deploy); err != nil {
		return nil, err
	}
	return h.updateDeploymentStatus(deploy)
}

// UpdateStatusFromObject updates deployment from metav1.Object or runtime.Object.
func (h *Handler) UpdateStatusFromObject(obj interface{}) (*appsv1.Deployment, error) {
	deploy, ok := obj.(*appsv1.Deployment)
	if !ok {
		return nil, fmt.Errorf("object type is not *appsv1.Deployment")
	}
	return h.updateDeploymentStatus(deploy)
}

// UpdateStatusFromUnstructured updates deployment from *unstructured.Unstructured.
func (h *Handler) UpdateStatusFromUnstructured(u *unstructured.Unstructured) (*appsv1.Deployment, error) {
	deploy := &appsv1.Deployment{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), deploy)
	if err != nil {
		return nil, err
	}
	return h.updateDeploymentStatus(deploy)
}

// UpdateStatusFromMap updates deployment from map[string]interface{}.
func (h *Handler) UpdateStatusFromMap(u map[string]interface{}) (*appsv1.Deployment, error) {
	deploy := &appsv1.Deployment{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, deploy)
	if err != nil {
		return nil, err
	}
	return h.updateDeploymentStatus(deploy)
}

// updateDeploymentStatus
func (h *Handler) updateDeploymentStatus(deploy *appsv1.Deployment) (*appsv1.Deployment, error) {
	namespace := deploy.GetNamespace()
	if len(namespace) == 0 {
		namespace = h.namespace
	}
	// resourceVersion cann't be set, the resourceVersion field is empty.
	deploy.UID = ""
	deploy.ResourceVersion = ""
	return h.clientset.AppsV1().Deployments(namespace).UpdateStatus(h.ctx, deploy, h.Options.UpdateOptions)
}
