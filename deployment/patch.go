package deployment

import (
	"encoding/json"
	"errors"
	"os"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
	"k8s.io/apimachinery/pkg/util/yaml"
)

/*
reference:
	https://loft.sh/blog/kubectl-patch-what-you-can-use-it-for-and-how-to-do-it
	https://github.com/kmodules/client-go/blob/201f259584dbffc8e4bb0c78fa96efdf812ff605/apps/v1/deployment.go#L63
	https://github.com/tamalsaha/patch-demo/blob/master/main.go
	https://github.com/kubernetes/kubectl/blob/master/pkg/cmd/patch/patch.go
	https://kubernetes.io/docs/tasks/manage-kubernetes-objects/update-api-object-kubectl-patch/
	https://erosb.github.io/post/json-patch-vs-merge-patch/
	https://stackoverflow.com/questions/61653702/scale-deployment-replicas-with-kubernetes-go-client
	https://stackoverflow.com/questions/53891862/patching-deployments-via-kubernetes-client-go
	https://github.com/kubernetes/client-go/issues/236
	https://dwmkerr.com/patching-kubernetes-resources-in-golang/
	https://caiorcferreira.github.io/post/the-kubernetes-dynamic-client/

Merge-patch: With a JSON merge patch, if you want to update a list, you have to
specify the entire new list. And the new list completely replaces the existing list.

Strategic-merge-patch: With a strategic merge patch, a list is either replaced
or merged depending on its patch strategy. The patch strategy is specified by
the value of the patchStrategy key in a field tag in the Kubernetes source code.
For example, the Containers field of PodSpec struct has a patchStrategy of merge:
type PodSpec struct {
    ...
    Containers []Container `json:"containers" patchStrategy:"merge" patchMergeKey:"name" ...`
*/

// Patch use the default patch type(Strategic Merge Patch) to patch deployment.
// Supported patch types are: "StrategicMergePatchType", "MergePatchType", "JSONPatchType".
//
// For further more Strategic Merge patch, see:
//     https://kubernetes.io/docs/tasks/manage-kubernetes-objects/update-api-object-kubectl-patch/#before-you-begin
// For a comparison of JSON patch and JSON merge patch, see:
//     https://erosb.github.io/post/json-patch-vs-merge-patch/
func (h *Handler) Patch(original *appsv1.Deployment, patch interface{}, patchOptions ...types.PatchType) (*appsv1.Deployment, error) {
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

	case *appsv1.Deployment:
		return h.diffMergePatch(original, val, patchOptions...)

	case appsv1.Deployment:
		return h.diffMergePatch(original, &val, patchOptions...)

	case map[string]interface{}:
		modified := &appsv1.Deployment{}
		if err := runtime.DefaultUnstructuredConverter.FromUnstructured(val, modified); err != nil {
			return nil, err
		}
		return h.diffMergePatch(original, modified, patchOptions...)

	case *unstructured.Unstructured:
		modified := &appsv1.Deployment{}
		if err := runtime.DefaultUnstructuredConverter.FromUnstructured(val.UnstructuredContent(), modified); err != nil {
			return nil, err
		}
		return h.diffMergePatch(original, modified, patchOptions...)

	case unstructured.Unstructured:
		modified := &appsv1.Deployment{}
		if err := runtime.DefaultUnstructuredConverter.FromUnstructured(val.UnstructuredContent(), modified); err != nil {
			return nil, err
		}
		return h.diffMergePatch(original, modified, patchOptions...)

	case runtime.Object:
		modified, ok := patch.(*appsv1.Deployment)
		if !ok {
			return nil, errors.New("patch data type is not *appsv1.Deployment")
		}
		return h.diffMergePatch(original, modified, patchOptions...)

	default:
		return nil, ErrInvalidPathType
	}
}

// strategicMergePatch use the "Strategic Merge Patch" patch type to patch deployment.
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
func (h *Handler) strategicMergePatch(original *appsv1.Deployment, patchData []byte) (*appsv1.Deployment, error) {
	if len(patchData) == 0 || string(patchData) == "{}" {
		return original, nil
	}

	var namespace string
	if len(original.Namespace) != 0 {
		namespace = original.Namespace
	} else {
		namespace = h.namespace
	}
	return h.clientset.AppsV1().Deployments(namespace).
		Patch(h.ctx, original.Name, types.StrategicMergePatchType, patchData, h.Options.PatchOptions)
}

// jsonMergePatch use the "JSON Merge Patch" patch type to patch deployment.
// A JSON merge patch is different from strategic merge patch, With a JSON merge patch,
// If you want to update a list, you have to specify the entire new list.
// And the new list completely replicas the existing list.
//
// For a comparison of JSON patch and JSON merge patch, see:
//     https://erosb.github.io/post/json-patch-vs-merge-patch/
// For further more Json Patch see:
//     https://kubernetes.io/docs/tasks/manage-kubernetes-objects/update-api-object-kubectl-patch/#before-you-begin
//     https://tools.ietf.org/html/rfc6902
func (h *Handler) jsonMergePatch(original *appsv1.Deployment, patchData []byte) (*appsv1.Deployment, error) {
	if len(patchData) == 0 || string(patchData) == "{}" {
		return original, nil
	}

	var namespace string
	if len(original.Namespace) != 0 {
		namespace = original.Namespace
	} else {
		namespace = h.namespace
	}
	return h.clientset.AppsV1().Deployments(namespace).
		Patch(h.ctx, original.Name, types.MergePatchType, patchData, h.Options.PatchOptions)
}

// jsonPatch use "JSON Patch" patch type to patch deployment.
//
// For a comparison of JSON patch and JSON merge patch, see:
//     https://erosb.github.io/post/json-patch-vs-merge-patch/
// For further more Json Merge Patch see:
//     https://kubernetes.io/docs/tasks/manage-kubernetes-objects/update-api-object-kubectl-patch/#before-you-begin
//     https://tools.ietf.org/html/rfc7386
func (h *Handler) jsonPatch(original *appsv1.Deployment, patchData []byte) (*appsv1.Deployment, error) {
	var namespace string
	if len(original.Namespace) != 0 {
		namespace = original.Namespace
	} else {
		namespace = h.namespace
	}
	return h.clientset.AppsV1().Deployments(namespace).Patch(h.ctx,
		original.Name, types.JSONPatchType, patchData, h.Options.PatchOptions)
}

// diffMergePatch will tak the difference data between original and modified deployment object,
// and use the default patch type(Strategic Merge Patch) patch the differen deployment.
// You can set patchOptions to MergePatchType to use the "JSON Merge Patch" to
// patch deployment.
func (h *Handler) diffMergePatch(original, modified *appsv1.Deployment, patchOptions ...types.PatchType) (*appsv1.Deployment, error) {
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
	if patchData, err = strategicpatch.CreateTwoWayMergePatch(originalJson, modifiedJson, appsv1.Deployment{}); err != nil {
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
		return h.clientset.AppsV1().Deployments(namespace).
			Patch(h.ctx, original.Name, types.MergePatchType, patchData, h.Options.PatchOptions)
	}
	return h.clientset.AppsV1().Deployments(namespace).
		Patch(h.ctx, original.Name, types.StrategicMergePatchType, patchData, h.Options.PatchOptions)
}
