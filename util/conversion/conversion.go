package conversion

import (
	"github.com/forbearing/k8s/types"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// KindToResource convert k8s kind name to k8s resource name.
// like:
//     KindToResource(Pod) --> pods
//     KindToResource(Ingress) -> ingresses
func KindToResource(kind string) string {
	return types.MapKindResource[kind]
}

// ResourceToKind convert k8s resource name to k8s kind name.
// like:
//     ResourceToKind(pods) -> Pod
//     ResourceToKind(networkpolicies) -> NetworkPolicy
func ResourceToKind(resource string) string {
	return types.MapResourceKind[resource]
}

// GVRToGVK convert schema.GroupVersionResource to schema.GroupVersionKind.
func GVRToGVK(gvr schema.GroupVersionResource) schema.GroupVersionKind {
	return schema.GroupVersionKind{
		Group:   gvr.Group,
		Version: gvr.Version,
		Kind:    ResourceToKind(gvr.Resource),
	}
}

// GVKToGVR convert schema.GroupVersionKind to schema.GroupVersionResource.
func GVKToGVR(gvk schema.GroupVersionKind) schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    gvk.Group,
		Version:  gvk.Version,
		Resource: KindToResource(gvk.Kind),
	}
}
