package labels

import (
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation"
)

/*
reference:
https://github.com/kubernetes/kubectl/blob/master/pkg/cmd/label/label.go
https://github.com/fenggolang/client-go-example/blob/master/vendor/k8s.io/apimachinery/pkg/labels/labels.go
*/

// Has will return true if the k8s object has specified label, otherwise return false.
// Label key and value separated by "=".
//
// If the provided label only contains label key, then only to check whether
// the labels of the k8s object contains the label key.
// If the provided label contains label key and value, then check whether
// the labels of the k8s object contains the label key and value.
func Has(obj runtime.Object, label string) bool {
	key, val, err := parseLabel(label)
	if err != nil {
		return false
	}

	// meta.Accessor convert runtime.Object to metav1.Object.
	// metav1.Object have all kinds of method to get/set k8s object metadata,
	// such like: GetNamespace/SetNamespace, GetName/SetName, GetLabels/SetLabels.
	accessor, err := meta.Accessor(obj)
	if err != nil {
		return false
	}

	// the label only contains label key, only to check whether the labels of
	// the k8s object contains the label key.
	if len(val) == 0 {
		for k := range accessor.GetLabels() {
			if k == key {
				return true
			}
		}
	}
	// the label contains label key and value, and check whether the labels of
	// the k8s object contains the label key and value.
	for k, v := range accessor.GetLabels() {
		if k == key && v == val {
			return true
		}
	}
	return false
}

// Get get the label value of the provided k8s object for the specified label key.
// Return empty string if the object is not runtime.Object.
func Get(obj runtime.Object, label string) string {
	key, _, _ := parseLabel(label)

	accessor, err := meta.Accessor(obj)
	if err != nil {
		return ""
	}
	return accessor.GetLabels()[key]
}

// GetAll get all labels of the provided k8s object.
// Return nil if the provided object is not runtime.Object.
func GetAll(obj runtime.Object) map[string]string {
	accessor, err := meta.Accessor(obj)
	if err != nil {
		return nil
	}
	return accessor.GetLabels()
}

// Set set labels for the provided k8s object.
// Label key and value separated by "=".
// If label already exist, it will update the label.
// If label not exist, it will add the label.
func Set(obj runtime.Object, label ...string) error {
	for _, l := range label {
		key, val, err := parseLabel(l)
		if err != nil {
			return err
		}

		accessor, err := meta.Accessor(obj)
		if err != nil {
			return err
		}
		labels := accessor.GetLabels()
		labels[key] = val
		accessor.SetLabels(labels)
	}
	return nil
}

// Remove will delete the provided label for the k8s object if contains it.
func Remove(obj runtime.Object, label ...string) error {
	for _, l := range label {
		key, val, err := parseLabel(l)
		if err != nil {
			return err
		}

		accessor, err := meta.Accessor(obj)
		if err != nil {
			return err
		}

		// label only contains label key
		newLabels := make(map[string]string)
		if len(val) == 0 {
			for k, v := range accessor.GetLabels() {
				if k == key {
					continue
				}
				newLabels[k] = v
			}
		}
		// label contains label key and label value
		for k, v := range accessor.GetLabels() {
			if k == key && v == val {
				continue
			}
			newLabels[k] = v
		}
		accessor.SetLabels(newLabels)
	}
	return nil
}

// RemoveAll will remove all labels of the k8s object.
func RemoveAll(obj runtime.Object) error {
	accessor, err := meta.Accessor(obj)
	if err != nil {
		return err
	}
	emptyLabels := make(map[string]string)
	accessor.SetLabels(emptyLabels)
	return nil
}

// parseLabel parse the label and return label key and label value.
// Label key and value separated by "=".
func parseLabel(label string) (key, val string, err error) {
	parts := strings.Split(label, "=")
	// label only contains label key
	if len(parts) == 1 {
		// return the label key
		return parts[0], "", nil
	}
	// label contains the label key and value, and validate label value.
	if errs := validation.IsValidLabelValue(parts[1]); len(errs) != 0 {
		return "", "", fmt.Errorf("invalid label value: %q: %v", label, strings.Join(errs, ";"))
	}
	// return the label key and value
	return parts[0], parts[1], nil
}
