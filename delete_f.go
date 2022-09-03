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

// DeleteF work like "kubectl delete -f filename.yaml -n test",
// The namespace defined in yaml have higher precedence than namespace specified here.
func DeleteF(ctx context.Context, kubeconfig, filename string, namespace string, opts ...Options) error {
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
		err = handler.Delete(item)
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

		// If the err returned by dynamic handler is "NotFound", just output the
		// error message and continue process the next items.
		// You can call DeleteF() with IgnoreNotFound option to ignore the "NotFound" error.
		// A "NotFound" error will occurrs when you delete k8s resource that no longer exist in cluster.
		if err != nil && apierrors.IsNotFound(err) {
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
