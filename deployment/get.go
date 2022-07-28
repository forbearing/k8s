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

// Get gets deployment from type string, []byte, *appsv1.Deployment,
// appsv1.Deployment, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.

// If passed parameter type is string, it will simply call GetByName instead of GetFromFile.
// You should always explicitly call GetFromFile to get a deployment from file path.
func (h *Handler) Get(obj interface{}) (*appsv1.Deployment, error) {
	switch val := obj.(type) {
	case string:
		return h.GetByName(val)
	case []byte:
		return h.GetFromBytes(val)
	case *appsv1.Deployment:
		return h.GetFromObject(val)
	case appsv1.Deployment:
		return h.GetFromObject(&val)
	case runtime.Object:
		return h.GetFromObject(val)
	case *unstructured.Unstructured:
		return h.GetFromUnstructured(val)
	case unstructured.Unstructured:
		return h.GetFromUnstructured(&val)
	case map[string]interface{}:
		return h.GetFromMap(val)
	default:
		return nil, ERR_TYPE_GET
	}
}

// GetByName gets deployment by name.
func (h *Handler) GetByName(name string) (*appsv1.Deployment, error) {
	return h.clientset.AppsV1().Deployments(h.namespace).Get(h.ctx, name, h.Options.GetOptions)
}

// GetFromFile gets deployment from yaml file.
func (h *Handler) GetFromFile(filename string) (*appsv1.Deployment, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.GetFromBytes(data)
}

// GetFromBytes gets deployment from bytes.
func (h *Handler) GetFromBytes(data []byte) (*appsv1.Deployment, error) {
	deployJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	deploy := &appsv1.Deployment{}
	err = json.Unmarshal(deployJson, deploy)
	if err != nil {
		return nil, err
	}
	return h.getDeployment(deploy)
}

// GetFromObject gets deployment from runtime.Object.
func (h *Handler) GetFromObject(obj runtime.Object) (*appsv1.Deployment, error) {
	deploy, ok := obj.(*appsv1.Deployment)
	if !ok {
		return nil, fmt.Errorf("object is not *appsv1.Deployment")
	}
	return h.getDeployment(deploy)
}

// GetFromUnstructured gets deployment from *unstructured.Unstructured.
func (h *Handler) GetFromUnstructured(u *unstructured.Unstructured) (*appsv1.Deployment, error) {
	deploy := &appsv1.Deployment{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), deploy)
	if err != nil {
		return nil, err
	}
	return h.getDeployment(deploy)
}

// GetFromMap gets deployment from map[string]interface{}.
func (h *Handler) GetFromMap(u map[string]interface{}) (*appsv1.Deployment, error) {
	deploy := &appsv1.Deployment{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, deploy)
	if err != nil {
		return nil, err
	}
	return h.getDeployment(deploy)
}

// getDeployment
// It's necessary to get a new deployment resource from a old deployment resource,
// because old deployment usually don't have deployment.Status field.
func (h *Handler) getDeployment(deploy *appsv1.Deployment) (*appsv1.Deployment, error) {
	var namespace string
	if len(deploy.Namespace) != 0 {
		namespace = deploy.Namespace
	} else {
		namespace = h.namespace
	}
	return h.clientset.AppsV1().Deployments(namespace).Get(h.ctx, deploy.Name, h.Options.GetOptions)
}
