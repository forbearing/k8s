package deployment

import (
	"encoding/json"
	"io/ioutil"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// CreateFromRaw create deployment from map[string]interface{}.
func (h *Handler) CreateFromRaw(raw map[string]interface{}) (*appsv1.Deployment, error) {
	deploy := &appsv1.Deployment{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(raw, deploy)
	if err != nil {
		return nil, err
	}

	var namespace string
	if len(deploy.Namespace) != 0 {
		namespace = deploy.Namespace
	} else {
		namespace = h.namespace
	}

	return h.clientset.AppsV1().Deployments(namespace).Create(h.ctx, deploy, h.Options.CreateOptions)
}

// CreateFromBytes create deployment from bytes.
func (h *Handler) CreateFromBytes(data []byte) (*appsv1.Deployment, error) {
	deployJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	deploy := &appsv1.Deployment{}
	err = json.Unmarshal(deployJson, deploy)
	if err != nil {
		return nil, err
	}

	var namespace string
	if len(deploy.Namespace) != 0 {
		namespace = deploy.Namespace
	} else {
		namespace = h.namespace
	}

	return h.clientset.AppsV1().Deployments(namespace).Create(h.ctx, deploy, h.Options.CreateOptions)
}

// CreateFromFile create deployment from yaml file.
func (h *Handler) CreateFromFile(filename string) (*appsv1.Deployment, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.CreateFromBytes(data)
}

// Create create deployment from yaml file, alias to "CreateFromBytes".
func (h *Handler) Create(filename string) (*appsv1.Deployment, error) {
	return h.CreateFromFile(filename)
}
