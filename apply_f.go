package k8s

import (
	"bytes"
	"context"
	"io/ioutil"
	"regexp"

	utilerrors "github.com/forbearing/k8s/util/errors"
	"github.com/sirupsen/logrus"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
)

// ApplyF work like "kubectl apply -f filename.yaml -n test",
// The namespace defined in yaml have higher precedence than namespace specified here.
func ApplyF(ctx context.Context, kubeconfig, filename string, namespace string, opts ...Options) error {
	handler, err := New(ctx, kubeconfig, namespace)
	if err != nil {
		return err
	}

	yamlData, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	// Remove all comments from yaml documents.
	removeComments := regexp.MustCompile(`#.*`)
	yamlData = removeComments.ReplaceAll(yamlData, []byte(""))
	// Split yaml documents into multiple single yaml document base on the delimiter("---")
	yamlList := bytes.Split(yamlData, []byte("---"))

	for _, item := range yamlList {
		// If the yaml document is empty, skip create it.
		if len(bytes.TrimSpace(item)) == 0 {
			continue
		}
		// If the k8s resource is cluster scope, the namespace specified in dynamic.New() will be ignored.
		// If the k8s resource is namespace scope and no namespace is defined in yaml file, then
		// dynaimc hanler will create the k8s resource is the namespace specified in dynamic.New().
		// (namespace defined in yaml file have higher precedence than specified in dynamic.New())
		_, err = handler.Apply(item)
		for _, opt := range opts {
			switch opt {
			case IgnoreAlreadyExists:
				err = utilerrors.IgnoreAlreadyExists(err)
			case IgnoreNotFound:
				err = utilerrors.IgnoreNotFound(err)
			case IgnoreInvalid:
				err = utilerrors.IgnoreInvalid(err)
			}
		}

		// If the error returned by dynamic handler is "AlreadyExists" or "Invalid",
		// just output the error message continue handle the next itmes.
		// You can call ApplyF() with IgnoreInvalid or/and IgnoreInvalid options to
		// ignore these errors.
		// A "Invalid" error will occurrs when you update the pod/job/persistentvolume resource.
		if err != nil && (apierrors.IsAlreadyExists(err) || apierrors.IsInvalid(err)) {
			logrus.Error(err)
			continue
		}
		// Unexpected error, return it.
		if err != nil {
			return err
		}
	}

	return nil
}
