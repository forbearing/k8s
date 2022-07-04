package deployment

import (
	"encoding/json"
	"io/ioutil"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// UpdateFromRaw update deployment from map[string]interface{}.
func (h *Handler) UpdateFromRaw(raw map[string]interface{}) (*appsv1.Deployment, error) {
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

	return h.clientset.AppsV1().Deployments(namespace).Update(h.ctx, deploy, h.Options.UpdateOptions)
}

// UpdateFromBytes update deployment from bytes.
func (h *Handler) UpdateFromBytes(data []byte) (*appsv1.Deployment, error) {
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

	return h.clientset.AppsV1().Deployments(namespace).Update(h.ctx, deploy, h.Options.UpdateOptions)
}

// UpdateFromFile update deployment from yaml file.
func (h *Handler) UpdateFromFile(filename string) (*appsv1.Deployment, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.UpdateFromBytes(data)
}

// Update update deployment from yaml file, alias to "UpdateFromFile".
func (h *Handler) Update(filename string) (deploy *appsv1.Deployment, err error) {
	return h.UpdateFromFile(filename)
}
