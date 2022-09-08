package deployment

import (
	"encoding/json"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
)

/*
reference:
	https://ymmt2005.hatenablog.com/entry/2020/04/14/An_example_of_using_dynamic_client_of_k8s.io/client-go#Mapping-between-GVK-and-GVR
		types.ApplyPatchType
	https://kubernetes.io/docs/tasks/manage-kubernetes-objects/update-api-object-kubectl-patch/
	https://stackoverflow.com/questions/61653702/scale-deployment-replicas-with-kubernetes-go-client
	https://stackoverflow.com/questions/53891862/patching-deployments-via-kubernetes-client-go
	https://github.com/kubernetes/client-go/issues/236
	https://github.com/tamalsaha/patch-demo
	https://dwmkerr.com/patching-kubernetes-resources-in-golang/
	https://golang.hotexamples.com/examples/k8s.io.client-go.pkg.labels/-/Everything/golang-everything-function-examples.html
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

type PatchType string

const (
	JSONPatchType           PatchType = PatchType(types.JSONPatchType)
	MergePatchType          PatchType = PatchType(types.MergePatchType)
	StrategicMergePatchType PatchType = PatchType(types.StrategicMergePatchType)
	ApplyPatchType          PatchType = PatchType(types.ApplyPatchType)
)

// Path default use the Strategic Merge Path, its simply call StrategicMergePatch().
func (h *Handler) Path() {}

// StrategicMergePatch use the strategic merge patch to patch deployment.
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
func (h *Handler) StrategicMergePatch(original, modified *appsv1.Deployment) (*appsv1.Deployment, error) {
	var (
		err          error
		originalJson []byte
		modifiedJson []byte
		patchData    []byte
		namespace    string
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

	if len(original.Namespace) != 0 {
		namespace = original.Namespace
	} else {
		namespace = h.namespace
	}

	return h.clientset.AppsV1().Deployments(namespace).Patch(h.ctx, original.Name,
		types.StrategicMergePatchType, patchData, h.Options.PatchOptions)
}

// JsonMergePath use JSON merge patch to patch deployment.
// A JSON merge patch is different from strategic merge patch, With a JSON merge patch,
// If you want to update a list, you have to specify the entire new list.
// And the new list completely replicas the existing list.
//
// For a comparison of JSON patch and JSON merge patch, see:
//     https://erosb.github.io/post/json-patch-vs-merge-patch/
// For further more Json Merge Patch see:
//     https://kubernetes.io/docs/tasks/manage-kubernetes-objects/update-api-object-kubectl-patch/#before-you-begin
//     https://tools.ietf.org/html/rfc7386
func (h *Handler) JsonMergePath(deploy *appsv1.Deployment, patchData []byte) (*appsv1.Deployment, error) {
	var namespace string
	if len(deploy.Namespace) != 0 {
		namespace = deploy.Namespace
	} else {
		namespace = h.namespace
	}

	return h.clientset.AppsV1().Deployments(namespace).Patch(h.ctx,
		deploy.Name, types.JSONPatchType, patchData, h.Options.PatchOptions)
}

// JsonPath
// For a comparison of JSON patch and JSON merge patch, see:
//     https://erosb.github.io/post/json-patch-vs-merge-patch/
// For further more Json Patch see:
//     https://kubernetes.io/docs/tasks/manage-kubernetes-objects/update-api-object-kubectl-patch/#before-you-begin
//     https://tools.ietf.org/html/rfc6902
func (h *Handler) JsonPath() {}

// MergePatch
func (h *Handler) MergePatch(original, modified *appsv1.Deployment) (*appsv1.Deployment, error) {
	var (
		err          error
		originalJson []byte
		modifiedJson []byte
		patchData    []byte
		namespace    string
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

	if len(original.Namespace) != 0 {
		namespace = original.Namespace
	} else {
		namespace = h.namespace
	}

	return h.clientset.AppsV1().Deployments(namespace).Patch(h.ctx, original.Name,
		types.MergePatchType, patchData, h.Options.PatchOptions)
}
