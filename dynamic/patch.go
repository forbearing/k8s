package dynamic

import (
	"encoding/json"
	"errors"
	"os"

	utilrestmapper "github.com/forbearing/k8s/util/restmapper"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// Patch use the default patch type(Strategic Merge Patch) to patch typed resource
// or unstructured object. Supported patch types are: "StrategicMergePatchType",
// "MergePatchType", "JSONPatchType".
//
// Note: if the patched object is Custom Resource, you should calling this method
// with types.MergePatchType, because "Strategic Merge Patch" does not work
// with CRD types.
//
// It's not necessary to explicitly specify the GVK or GVR by calling WithGVK(),
// Patch() will find the GVK and GVR by RESTMapper and patch the k8s resource
// that defined in original *unstructured.Unstructured.
//
// For further more Strategic Merge patch, see:
//     https://kubernetes.io/docs/tasks/manage-kubernetes-objects/update-api-object-kubectl-patch/#before-you-begin
// For a comparison of JSON patch and JSON merge patch, see:
//     https://erosb.github.io/post/json-patch-vs-merge-patch/
func (h *Handler) Patch(original *unstructured.Unstructured, patch interface{}, patchOptions ...types.PatchType) (*unstructured.Unstructured, error) {
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

	case *unstructured.Unstructured:
		return h.diffMergePatch(original, val, patchOptions...)

	case unstructured.Unstructured:
		return h.diffMergePatch(original, &val, patchOptions...)

	case map[string]interface{}:
		return h.diffMergePatch(original, &unstructured.Unstructured{Object: val}, patchOptions...)

	case runtime.Object:
		modified, ok := patch.(*unstructured.Unstructured)
		if !ok {
			return nil, errors.New("patch data type is not *unstructured.Unstructured")
		}
		return h.diffMergePatch(original, modified, patchOptions...)

	default:
		return nil, ErrInvalidPatchType
	}
}

// strategicMergePatch use the "Strategic Merge Patch" patch type to patch pod.
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
func (h *Handler) strategicMergePatch(original *unstructured.Unstructured, patchData []byte) (*unstructured.Unstructured, error) {
	if len(patchData) == 0 || string(patchData) == "{}" {
		return original, nil
	}
	return h.patchUnstructured(original, patchData, types.StrategicMergePatchType)
}

// jsonMergePatch use the "JSON Merge Patch" patch type to patch pod.
// A JSON merge patch is different from strategic merge patch, With a JSON merge patch,
// If you want to update a list, you have to specify the entire new list.
// And the new list completely replicas the existing list.
//
// For a comparison of JSON patch and JSON merge patch, see:
//     https://erosb.github.io/post/json-patch-vs-merge-patch/
// For further more Json Patch see:
//     https://kubernetes.io/docs/tasks/manage-kubernetes-objects/update-api-object-kubectl-patch/#before-you-begin
//     https://tools.ietf.org/html/rfc6902
func (h *Handler) jsonMergePatch(original *unstructured.Unstructured, patchData []byte) (*unstructured.Unstructured, error) {
	if len(patchData) == 0 || string(patchData) == "{}" {
		return original, nil
	}
	return h.patchUnstructured(original, patchData, types.MergePatchType)
}

// jsonPatch use "JSON Patch" patch type to patch pod.
//
// For a comparison of JSON patch and JSON merge patch, see:
//     https://erosb.github.io/post/json-patch-vs-merge-patch/
// For further more Json Merge Patch see:
//     https://kubernetes.io/docs/tasks/manage-kubernetes-objects/update-api-object-kubectl-patch/#before-you-begin
//     https://tools.ietf.org/html/rfc7386
func (h *Handler) jsonPatch(original *unstructured.Unstructured, patchData []byte) (*unstructured.Unstructured, error) {
	return h.patchUnstructured(original, patchData, types.JSONPatchType)
}

// diffMergePatch will tak the difference data between original and modified pod object,
// and use the default patch type(Strategic Merge Patch) patch the differen pod.
// You can set patchOptions to MergePatchType to use the "JSON Merge Patch" to
// patch pod.
//
// Note: if the patched object is Custom Resource, you should calling this method
// with types.MergePatchType, because "Strategic Merge Patch" does not work
// with CRD types.
func (h *Handler) diffMergePatch(original, modified *unstructured.Unstructured, patchOptions ...types.PatchType) (*unstructured.Unstructured, error) {
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
	if patchData, err = strategicpatch.CreateTwoWayMergePatch(originalJson, modifiedJson, unstructured.Unstructured{}); err != nil {
		return nil, err
	}
	if len(patchData) == 0 || string(patchData) == "{}" {
		return original, nil
	}

	if len(patchOptions) != 0 && patchOptions[0] == types.MergePatchType {
		return h.patchUnstructured(original, patchData, types.MergePatchType)
	}
	return h.patchUnstructured(original, patchData, types.StrategicMergePatchType)
}

// patchUnstructured
func (h *Handler) patchUnstructured(obj *unstructured.Unstructured, patchData []byte, patchType types.PatchType) (*unstructured.Unstructured, error) {
	var (
		err          error
		gvk          schema.GroupVersionKind
		gvr          schema.GroupVersionResource
		isNamespaced bool
	)

	if gvr, err = utilrestmapper.FindGVR(h.restMapper, obj); err != nil {
		return nil, err
	}
	if gvk, err = utilrestmapper.FindGVK(h.restMapper, obj); err != nil {
		return nil, err
	}
	if isNamespaced, err = utilrestmapper.IsNamespaced(h.restMapper, gvk); err != nil {
		return nil, err
	}

	var namespace string
	if isNamespaced {
		if len(obj.GetNamespace()) != 0 {
			namespace = obj.GetNamespace()
		} else {
			namespace = h.namespace
		}
		return h.dynamicClient.Resource(gvr).Namespace(namespace).Patch(h.ctx, obj.GetName(), patchType, patchData, h.Options.PatchOptions)
	}
	return h.dynamicClient.Resource(gvr).Patch(h.ctx, obj.GetName(), patchType, patchData, h.Options.PatchOptions)
}
