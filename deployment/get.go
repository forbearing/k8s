package deployment

import (
	"encoding/json"
	"io/ioutil"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// GetFromBytes get deployment from bytes.
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

	var namespace string
	if len(deploy.Namespace) != 0 {
		namespace = deploy.Namespace
	} else {
		namespace = h.namespace
	}

	return h.WithNamespace(namespace).GetByName(deploy.Name)
}

// GetFromFile get deployment from yaml file.
func (h *Handler) GetFromFile(filename string) (*appsv1.Deployment, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.GetFromBytes(data)
}

// GetByName get deployment by name.
func (h *Handler) GetByName(name string) (*appsv1.Deployment, error) {
	return h.clientset.AppsV1().Deployments(h.namespace).Get(h.ctx, name, h.Options.GetOptions)
}

// Get get deployment by name, alias to "GetByName".
func (h *Handler) Get(name string) (*appsv1.Deployment, error) {
	return h.clientset.AppsV1().Deployments(h.namespace).Get(h.ctx, name, h.Options.GetOptions)
}
