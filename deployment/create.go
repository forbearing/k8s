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

// Create creates deployment from type string, []byte, *appsv1.Deployment,
// appsv1.Deployment, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.
func (h *Handler) Create(obj interface{}) (*appsv1.Deployment, error) {
	switch val := obj.(type) {
	case string:
		return h.CreateFromFile(val)
	case []byte:
		return h.CreateFromBytes(val)
	case *appsv1.Deployment:
		return h.CreateFromObject(val)
	case appsv1.Deployment:
		return h.CreateFromObject(&val)
	case runtime.Object:
		return h.CreateFromObject(val)
	case *unstructured.Unstructured:
		// - 如果传入的类型是 *unstructured.Unstructured 做类型断言时,它会自动转换成
		//   runtime.Object 类型, 而不是 *unstructured.Unstructured
		// - 所以不支持从 *unstructured.Unstructured 来创建 pod
		//   只支持从 unstructured.Unstructured 来创建 pod
		return h.CreateFromUnstructured(val)
	case unstructured.Unstructured:
		return h.CreateFromUnstructured(&val)
	case map[string]interface{}:
		return h.CreateFromMap(val)
	default:
		return nil, ERR_TYPE_CREATE
	}
}

// CreateFromFile creates deployment from yaml file.
func (h *Handler) CreateFromFile(filename string) (*appsv1.Deployment, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.CreateFromBytes(data)
}

// CreateFromBytes creates deployment from bytes.
func (h *Handler) CreateFromBytes(data []byte) (*appsv1.Deployment, error) {
	deployJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	deploy := &appsv1.Deployment{}
	if err = json.Unmarshal(deployJson, deploy); err != nil {
		return nil, err
	}
	return h.createDeployment(deploy)
}

// CreateFromObject creates deployment from runtime.Object.
func (h *Handler) CreateFromObject(obj runtime.Object) (*appsv1.Deployment, error) {
	deploy, ok := obj.(*appsv1.Deployment)
	if !ok {
		return nil, fmt.Errorf("object type is not *appsv1.Deployment")
	}
	return h.createDeployment(deploy)
}

// CreateFromUnstructured creates deployment from *unstructured.Unstructured.
func (h *Handler) CreateFromUnstructured(u *unstructured.Unstructured) (*appsv1.Deployment, error) {
	deploy := &appsv1.Deployment{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), deploy)
	if err != nil {
		return nil, err
	}
	return h.createDeployment(deploy)
}

// CreateFromMap creates deployment from map[string]interface{}.
func (h *Handler) CreateFromMap(u map[string]interface{}) (*appsv1.Deployment, error) {
	deploy := &appsv1.Deployment{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, deploy)
	if err != nil {
		return nil, err
	}
	return h.createDeployment(deploy)
}

// createDeployment
func (h *Handler) createDeployment(deploy *appsv1.Deployment) (*appsv1.Deployment, error) {
	// TODO: Check if the *appsv1.deployment resource always has Namespace field
	// to explicitly specify in which namespace the current deployment resource runs.
	// If deployment resource always has a Namespace field, and the Namespace field
	// always not empty, then additionally setting namespace is not nedded.
	var namespace string
	if len(deploy.Namespace) != 0 {
		namespace = deploy.Namespace
	} else {
		namespace = h.namespace
	}
	// resourceVersion must be empty, otherwise the error
	// "resourceVersion should not be set on objects to be created" will be returned.
	deploy.ResourceVersion = ""
	deploy.UID = ""
	return h.clientset.AppsV1().Deployments(namespace).Create(h.ctx, deploy, h.Options.CreateOptions)
}
