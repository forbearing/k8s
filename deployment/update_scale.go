package deployment

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// Scale set deployment replicas from type string, []byte, *appsv1.Deployment,
// appsv1.Deployment, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.
//
// If passed parameter type is string, it will simply call ScaleByName instead of ScaleFromFile.
// You should always explicitly call ScaleFromFile to set deployment replicas from file path.
func (h *Handler) Scale(obj interface{}, replicas int32) (*appsv1.Deployment, error) {
	switch val := obj.(type) {
	case string:
		return h.ScaleByName(val, replicas)
	case []byte:
		return h.ScaleFromBytes(val, replicas)
	case *appsv1.Deployment:
		return h.ScaleFromObject(val, replicas)
	case appsv1.Deployment:
		return h.ScaleFromObject(&val, replicas)
	case runtime.Object:
		if reflect.TypeOf(val).String() == "*unstructured.Unstructured" {
			return h.ScaleFromUnstructured(val.(*unstructured.Unstructured), replicas)
		}
		return h.ScaleFromObject(val, replicas)
	case *unstructured.Unstructured:
		return h.ScaleFromUnstructured(val, replicas)
	case unstructured.Unstructured:
		return h.ScaleFromUnstructured(&val, replicas)
	case map[string]interface{}:
		return h.ScaleFromMap(val, replicas)
	default:
		return nil, ErrInvalidScaleType
	}
}

// ScaleByName scale deployment by name.
func (h *Handler) ScaleByName(name string, replicas int32) (*appsv1.Deployment, error) {
	deploy, err := h.Get(name)
	if err != nil {
		return nil, err
	}
	copiedDeploy := deploy.DeepCopy()
	if copiedDeploy.Spec.Replicas != nil {
		copiedDeploy.Spec.Replicas = &replicas
	}
	return h.Update(copiedDeploy)

	//scale := &autoscalingv1.Scale{}
	//scale.Spec.Replicas = replicas
	//_, err := h.clientset.AppsV1().Deployments(h.namespace).UpdateScale(h.ctx, name, scale, h.Options.UpdateOptions)
	//return nil, err
}

// ScaleFromFile scale deployment from yaml file.
func (h *Handler) ScaleFromFile(filename string, replicas int32) (*appsv1.Deployment, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.ScaleFromBytes(data, replicas)
}

// ScaleFromBytes scale deployment from bytes.
func (h *Handler) ScaleFromBytes(data []byte, replicas int32) (*appsv1.Deployment, error) {
	deployJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	deploy := &appsv1.Deployment{}
	if err = json.Unmarshal(deployJson, deploy); err != nil {
		return nil, err
	}
	return h.ScaleByName(deploy.Name, replicas)
}

// ScaleFromObject scale deployment from runtime.Object.
func (h *Handler) ScaleFromObject(obj runtime.Object, replicas int32) (*appsv1.Deployment, error) {
	deploy, ok := obj.(*appsv1.Deployment)
	if !ok {
		return nil, fmt.Errorf("object type is not *appsv1.Deployment")
	}
	return h.ScaleByName(deploy.Name, replicas)
}

// ScaleFromUnstructured scale deployment from *unstructured.Unstructured.
func (h *Handler) ScaleFromUnstructured(u *unstructured.Unstructured, replicas int32) (*appsv1.Deployment, error) {
	deploy := &appsv1.Deployment{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), deploy)
	if err != nil {
		return nil, err
	}
	return h.ScaleByName(deploy.Name, replicas)
}

// ScaleFromMap scale deployment from map[string]interface{}.
func (h *Handler) ScaleFromMap(u map[string]interface{}, replicas int32) (*appsv1.Deployment, error) {
	deploy := &appsv1.Deployment{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, deploy)
	if err != nil {
		return nil, err
	}
	return h.ScaleByName(deploy.Name, replicas)
}
