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

// Delete deletes deployment from type string, []byte, *appsv1.Deployment,
// appsv1.Deployment, metav1.Object, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.
//
// If passed parameter type is string, it will simply call DeleteByName instead of DeleteFromFile.
// You should always explicitly call DeleteFromFile to delete a deployment from file path.
func (h *Handler) Delete(obj interface{}) error {
	switch val := obj.(type) {
	case string:
		return h.DeleteByName(val)
	case []byte:
		return h.DeleteFromBytes(val)
	case *appsv1.Deployment:
		return h.DeleteFromObject(val)
	case appsv1.Deployment:
		return h.DeleteFromObject(&val)
	case *unstructured.Unstructured:
		return h.DeleteFromUnstructured(val)
	case unstructured.Unstructured:
		return h.DeleteFromUnstructured(&val)
	case map[string]interface{}:
		return h.DeleteFromMap(val)
	case metav1.Object, runtime.Object:
		return h.DeleteFromObject(val)
	default:
		return ErrInvalidDeleteType
	}
}

// DeleteByName deletes deployment by name.
func (h *Handler) DeleteByName(name string) error {
	return h.clientset.AppsV1().Deployments(h.namespace).Delete(h.ctx, name, h.Options.DeleteOptions)
}

// DeleteFromFile deletes deployment from yaml or json file.
func (h *Handler) DeleteFromFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return h.DeleteFromBytes(data)
}

// DeleteFromBytes deletes deployment from bytes data.
func (h *Handler) DeleteFromBytes(data []byte) error {
	deployJson, err := yaml.ToJSON(data)
	if err != nil {
		return err
	}

	deploy := &appsv1.Deployment{}
	if err = json.Unmarshal(deployJson, deploy); err != nil {
		return err
	}
	return h.deleteDeployment(deploy)
}

// DeleteFromObject deletes deployment from metav1.Object or runtime.Object.
func (h *Handler) DeleteFromObject(obj interface{}) error {
	deploy, ok := obj.(*appsv1.Deployment)
	if !ok {
		return fmt.Errorf("object type is not *appsv1.Deployment")
	}
	return h.deleteDeployment(deploy)
}

// DeleteFromUnstructured deletes deployment from *unstructured.Unstructured.
func (h *Handler) DeleteFromUnstructured(u *unstructured.Unstructured) error {
	deploy := &appsv1.Deployment{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), deploy)
	if err != nil {
		return err
	}
	return h.deleteDeployment(deploy)
}

// DeleteFromMap deletes deployment from map[string]interface{}.
func (h *Handler) DeleteFromMap(u map[string]interface{}) error {
	deploy := &appsv1.Deployment{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, deploy)
	if err != nil {
		return err
	}
	return h.deleteDeployment(deploy)
}

// deleteDeployment
func (h *Handler) deleteDeployment(deploy *appsv1.Deployment) error {
	namespace := deploy.GetNamespace()
	if len(namespace) == 0 {
		namespace = h.namespace
	}
	return h.clientset.AppsV1().Deployments(namespace).Delete(h.ctx, deploy.Name, h.Options.DeleteOptions)
}
