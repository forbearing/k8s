package deployment

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// UpdateStatus updates deployment from type string, []byte, *appsv1.Deployment,
// appsv1.Deployment, runtime.Object, *unstructured.Unstructured,
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
	//case runtime.Object:
	//    if reflect.TypeOf(val).String() == "*unstructured.Unstructured" {
	//        return h.UpdateStatusFromUnstructured(val.(*unstructured.Unstructured))
	//    }
	//    return h.UpdateStatusFromObject(val)
	case *unstructured.Unstructured:
		return h.UpdateStatusFromUnstructured(val)
	case unstructured.Unstructured:
		return h.UpdateStatusFromUnstructured(&val)
	case map[string]interface{}:
		return h.UpdateStatusFromMap(val)
	default:
		return nil, ErrInvalidUpdateType
	}
}

// UpdateStatusFromFile updates deployment from yaml file.
func (h *Handler) UpdateStatusFromFile(filename string) (*appsv1.Deployment, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.UpdateStatusFromBytes(data)
}

// UpdateStatusFromBytes updates deployment from bytes.
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

// UpdateStatusFromObject updates deployment from runtime.Object.
func (h *Handler) UpdateStatusFromObject(obj runtime.Object) (*appsv1.Deployment, error) {
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
	var namespace string
	if len(deploy.Namespace) != 0 {
		namespace = deploy.Namespace
	} else {
		namespace = h.namespace
	}
	// resourceVersion cann't be set, the resourceVersion field is empty.
	deploy.UID = ""
	deploy.ResourceVersion = ""
	return h.clientset.AppsV1().Deployments(namespace).UpdateStatus(h.ctx, deploy, h.Options.UpdateOptions)
}
