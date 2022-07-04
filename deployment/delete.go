package deployment

import (
	"encoding/json"
	"io/ioutil"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// DeleteFromBytes delete deploy from bytes.
func (h *Handler) DeleteFromBytes(data []byte) error {
	deployJson, err := yaml.ToJSON(data)
	if err != nil {
		return err
	}

	deploy := &appsv1.Deployment{}
	err = json.Unmarshal(deployJson, deploy)
	if err != nil {
		return err
	}

	var namespace string
	if len(deploy.Namespace) != 0 {
		namespace = deploy.Namespace
	} else {
		namespace = h.namespace
	}

	return h.WithNamespace(namespace).DeleteByName(deploy.Name)
}

// DeleteFromFile delete deployment from yaml file.
func (h *Handler) DeleteFromFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return h.DeleteFromBytes(data)
}

// DeleteByName delete deployment by name
func (h *Handler) DeleteByName(name string) (err error) {
	return h.clientset.AppsV1().Deployments(h.namespace).Delete(h.ctx, name, h.Options.DeleteOptions)
}

// Delete delete deployment by name, alias to "DeleteByName".
func (h *Handler) Delete(name string) error {
	return h.DeleteByName(name)
}
