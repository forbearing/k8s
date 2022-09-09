package rolebinding

import (
	"encoding/json"
	"errors"
	"os"

	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// Patch use the default patch type(Strategic Merge Patch) to patch rolebinding.
// Supported patch types are: "StrategicMergePatchType", "MergePatchType", "JSONPatchType".
//
// For further more Strategic Merge patch, see:
//     https://kubernetes.io/docs/tasks/manage-kubernetes-objects/update-api-object-kubectl-patch/#before-you-begin
// For a comparison of JSON patch and JSON merge patch, see:
//     https://erosb.github.io/post/json-patch-vs-merge-patch/
func (h *Handler) Patch(original *rbacv1.RoleBinding, patch interface{}, patchOptions ...types.PatchType) (*rbacv1.RoleBinding, error) {
	switch val := patch.(type) {
	case string:
		var err error
		var patchData []byte
		var jsonData []byte

		if patchData, err = os.ReadFile(val); err != nil {
			return nil, err
		}
		if jsonData, err = yaml.ToJSON(patchData); err != nil {
			return nil, err
		}
		if len(patchOptions) != 0 && patchOptions[0] == types.JSONPatchType {
			return h.jsonPatch(original, jsonData)
		}
		if len(patchOptions) != 0 && patchOptions[0] == types.MergePatchType {
			return h.jsonMergePatch(original, jsonData)
		}
		return h.strategicMergePatch(original, jsonData)

	case []byte:
		var err error
		var jsonData []byte

		if jsonData, err = yaml.ToJSON(val); err != nil {
			return nil, err
		}
		if len(patchOptions) != 0 && patchOptions[0] == types.JSONPatchType {
			return h.jsonPatch(original, jsonData)
		}
		if len(patchOptions) != 0 && patchOptions[0] == types.MergePatchType {
			return h.jsonMergePatch(original, jsonData)
		}
		return h.strategicMergePatch(original, jsonData)

	case *rbacv1.RoleBinding:
		return h.diffMergePatch(original, val, patchOptions...)

	case rbacv1.RoleBinding:
		return h.diffMergePatch(original, &val, patchOptions...)

	case map[string]interface{}:
		modified := &rbacv1.RoleBinding{}
		if err := runtime.DefaultUnstructuredConverter.FromUnstructured(val, modified); err != nil {
			return nil, err
		}
		return h.diffMergePatch(original, modified, patchOptions...)

	case *unstructured.Unstructured:
		modified := &rbacv1.RoleBinding{}
		if err := runtime.DefaultUnstructuredConverter.FromUnstructured(val.UnstructuredContent(), modified); err != nil {
			return nil, err
		}
		return h.diffMergePatch(original, modified, patchOptions...)

	case unstructured.Unstructured:
		modified := &rbacv1.RoleBinding{}
		if err := runtime.DefaultUnstructuredConverter.FromUnstructured(val.UnstructuredContent(), modified); err != nil {
			return nil, err
		}
		return h.diffMergePatch(original, modified, patchOptions...)

	case runtime.Object:
		modified, ok := patch.(*rbacv1.RoleBinding)
		if !ok {
			return nil, errors.New("patch data type is not *rbacv1.RoleBinding")
		}
		return h.diffMergePatch(original, modified, patchOptions...)

	default:
		return nil, ErrInvalidPathType
	}
}

// strategicMergePatch use the "Strategic Merge Patch" patch type to patch rolebinding.
//
// Notice that the patch did not replace the containers list. Instead it added
// a new Container to the list. In other words, the list in the patch was merged
// with the existing list.
//
// This is not always what happens when you use a strategic merge patch on a list.
// In some cases, the list is replaced, not merged.
//
// Note: Strategic merge patch is not supported for custom resources.
// For further more Strategic Merge patch, see:
//     https://kubernetes.io/docs/tasks/manage-kubernetes-objects/update-api-object-kubectl-patch/#before-you-begin
func (h *Handler) strategicMergePatch(original *rbacv1.RoleBinding, patchData []byte) (*rbacv1.RoleBinding, error) {
	if len(patchData) == 0 || string(patchData) == "{}" {
		return original, nil
	}

	var namespace string
	if len(original.Namespace) != 0 {
		namespace = original.Namespace
	} else {
		namespace = h.namespace
	}
	return h.clientset.RbacV1().RoleBindings(namespace).
		Patch(h.ctx, original.Name, types.StrategicMergePatchType, patchData, h.Options.PatchOptions)
}

// jsonMergePatch use the "JSON Merge Patch" patch type to patch rolebinding.
// A JSON merge patch is different from strategic merge patch, With a JSON merge patch,
// If you want to update a list, you have to specify the entire new list.
// And the new list completely replicas the existing list.
//
// For a comparison of JSON patch and JSON merge patch, see:
//     https://erosb.github.io/post/json-patch-vs-merge-patch/
// For further more Json Patch see:
//     https://kubernetes.io/docs/tasks/manage-kubernetes-objects/update-api-object-kubectl-patch/#before-you-begin
//     https://tools.ietf.org/html/rfc6902
func (h *Handler) jsonMergePatch(original *rbacv1.RoleBinding, patchData []byte) (*rbacv1.RoleBinding, error) {
	if len(patchData) == 0 || string(patchData) == "{}" {
		return original, nil
	}

	var namespace string
	if len(original.Namespace) != 0 {
		namespace = original.Namespace
	} else {
		namespace = h.namespace
	}
	return h.clientset.RbacV1().RoleBindings(namespace).
		Patch(h.ctx, original.Name, types.StrategicMergePatchType, patchData, h.Options.PatchOptions)
}

// jsonPatch use "JSON Patch" patch type to patch rolebinding.
//
// For a comparison of JSON patch and JSON merge patch, see:
//     https://erosb.github.io/post/json-patch-vs-merge-patch/
// For further more Json Merge Patch see:
//     https://kubernetes.io/docs/tasks/manage-kubernetes-objects/update-api-object-kubectl-patch/#before-you-begin
//     https://tools.ietf.org/html/rfc7386
func (h *Handler) jsonPatch(original *rbacv1.RoleBinding, patchData []byte) (*rbacv1.RoleBinding, error) {
	var namespace string
	if len(original.Namespace) != 0 {
		namespace = original.Namespace
	} else {
		namespace = h.namespace
	}
	return h.clientset.RbacV1().RoleBindings(namespace).Patch(h.ctx,
		original.Name, types.JSONPatchType, patchData, h.Options.PatchOptions)
}

// diffMergePatch will tak the difference data between original and modified rolebinding object,
// and use the default patch type(Strategic Merge Patch) patch the differen rolebinding.
// You can set patchOptions to MergePatchType to use the "JSON Merge Patch" to
// patch rolebinding.
func (h *Handler) diffMergePatch(original, modified *rbacv1.RoleBinding, patchOptions ...types.PatchType) (*rbacv1.RoleBinding, error) {
	var (
		err          error
		originalJson []byte
		modifiedJson []byte
		patchData    []byte
	)

	if originalJson, err = json.Marshal(original); err != nil {
		return nil, err
	}
	if modifiedJson, err = json.Marshal(modified); err != nil {
		return nil, err
	}
	if patchData, err = strategicpatch.CreateTwoWayMergePatch(originalJson, modifiedJson, rbacv1.RoleBinding{}); err != nil {
		return nil, err
	}
	if len(patchData) == 0 || string(patchData) == "{}" {
		return original, nil
	}

	var namespace string
	if len(original.Namespace) != 0 {
		namespace = original.Namespace
	} else {
		namespace = h.namespace
	}
	if len(patchOptions) != 0 && patchOptions[0] == types.MergePatchType {
		return h.clientset.RbacV1().RoleBindings(namespace).
			Patch(h.ctx, original.Name, types.MergePatchType, patchData, h.Options.PatchOptions)
	}
	return h.clientset.RbacV1().RoleBindings(namespace).
		Patch(h.ctx, original.Name, types.StrategicMergePatchType, patchData, h.Options.PatchOptions)
}
