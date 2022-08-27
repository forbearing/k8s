package annotations

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
https://github.com/kubernetes/kubectl/blob/master/pkg/cmd/annotate/annotate.go
https://github.com/fenggolang/client-go-example/blob/master/vendor/k8s.io/apimachinery/pkg/annotations/annotations.go
*/

// Has will return true if the k8s object has specified annotation, otherwise return false.
// Label key and value separated by "=".
//
// If the provided annotation only contains annotation key, then only to check whether
// the annotations of the k8s object contains the annotation key.
// If the provided annotation contains annotation key and value, then check whether
// the annotations of the k8s object contains the annotation key and value.
func Has(obj runtime.Object, annotation string) bool {
	// meta.Accessor convert runtime.Object to metav1.Object.
	// metav1.Object have all kinds of method to get/set k8s object metadata,
	// such like: GetNamespace/SetNamespace, GetName/SetName, GetLabels/SetLabels, etc.
	accessor, err := meta.Accessor(obj)
	if err != nil {
		return false
	}

	key, val, err := parseAnnotations(annotation)
	if err != nil {
		return false
	}

	// the annotation only contains annotation key, only to check whether the annotations of
	// the k8s object contains the annotation key.
	if len(val) == 0 {
		for k := range accessor.GetAnnotations() {
			if k == key {
				return true
			}
		}
	}
	// the annotation contains annotation key and value, and check whether the annotations of
	// the k8s object contains the annotation key and value.
	for k, v := range accessor.GetAnnotations() {
		if k == key && v == val {
			return true
		}
	}
	return false
}

// Get get the annotation value of the provided k8s object for the specified annotation key.
// Return empty string if the object is not runtime.Object.
func Get(obj runtime.Object, annotation string) string {
	accessor, err := meta.Accessor(obj)
	if err != nil {
		return ""
	}
	key, _, _ := parseAnnotations(annotation)
	return accessor.GetAnnotations()[key]
}

// GetAll get all annotations of the provided k8s object.
// Return nil if the provided object is not runtime.Object.
func GetAll(obj runtime.Object) map[string]string {
	accessor, err := meta.Accessor(obj)
	if err != nil {
		return nil
	}
	return accessor.GetAnnotations()
}

// Set set annotations for the provided k8s object.
// Label key and value separated by "=".
// If annotation already exist, it will update the annotation.
// If annotation not exist, it will add the annotation.
func Set(obj runtime.Object, annotation ...string) error {
	accessor, err := meta.Accessor(obj)
	if err != nil {
		return err
	}

	for _, l := range annotation {
		key, val, err := parseAnnotations(l)
		if err != nil {
			return err
		}
		annotations := accessor.GetAnnotations()
		annotations[key] = val
		accessor.SetAnnotations(annotations)
	}
	return nil
}

// Remove will delete the provided annotation for the k8s object if contains it.
func Remove(obj runtime.Object, annotation ...string) error {
	accessor, err := meta.Accessor(obj)
	if err != nil {
		return err
	}

	for _, l := range annotation {
		key, val, err := parseAnnotations(l)
		if err != nil {
			return err
		}
		// annotation only contains annotation key
		newLabels := make(map[string]string)
		if len(val) == 0 {
			for k, v := range accessor.GetAnnotations() {
				if k == key {
					continue
				}
				newLabels[k] = v
			}
		}
		// annotation contains annotation key and annotation value
		for k, v := range accessor.GetAnnotations() {
			if k == key && v == val {
				continue
			}
			newLabels[k] = v
		}
		accessor.SetAnnotations(newLabels)
	}
	return nil
}

// RemoveAll will remove all annotations of the k8s object.
func RemoveAll(obj runtime.Object) error {
	accessor, err := meta.Accessor(obj)
	if err != nil {
		return err
	}
	emptyLabels := make(map[string]string)
	accessor.SetAnnotations(emptyLabels)
	return nil
}

// parseAnnotations parse the annotation and return annotation key and value.
// Label key and value separated by "=".
func parseAnnotations(annotation string) (key, val string, err error) {
	parts := strings.Split(annotation, "=")
	// the annotation only contains key
	if len(parts) == 1 {
		// return the annotation key
		return parts[0], "", nil
	}
	// the annotation contains the key and value, and validate annotation value.
	if errs := validation.IsValidLabelValue(parts[1]); len(errs) != 0 {
		return "", "", fmt.Errorf("invalid annotation value: %q: %v", annotation, strings.Join(errs, ";"))
	}
	// return the annotation key and value
	return parts[0], parts[1], nil
}
