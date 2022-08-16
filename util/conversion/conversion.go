package conversion

import (
	"strings"

	. "github.com/forbearing/k8s/types"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	allKind []string = []string{
		KindClusterRole,
		KindClusterRoleBinding,
		KindClusterRole,
		KindClusterRoleBinding,
		KindConfigMap,
		KindCronJob,
		KindDaemonSet,
		KindDeployment,
		KindIngress,
		KindIngressClass,
		KindJob,
		KindNamespace,
		KindNetworkPolicy,
		KindNode,
		KindPersistentVolume,
		KindPersistentVolumeClaim,
		KindPod,
		KindReplicaSet,
		KindReplicationController,
		KindRole,
		KindRoleBinding,
		KindSecret,
		KindService,
		KindServiceAccount,
		KindStatefulSet,
		KindStorageClass,
	}

	allResource []string = []string{
		ResourceClusterRole,
		ResourceClusterRoleBinding,
		ResourceConfigMap,
		ResourceCronJob,
		ResourceDaemonSet,
		ResourceDeployment,
		ResourceIngress,
		ResourceIngressClass,
		ResourceJob,
		ResourceNamespace,
		ResourceNetworkPolicy,
		ResourceNode,
		ResourcePersistentVolume,
		ResourcePersistentVolumeClaim,
		ResourcePod,
		ResourceReplicaSet,
		ResourceReplicationController,
		ResourceRole,
		ResourceRoleBinding,
		ResourceSecret,
		ResourceService,
		ResourceServiceAccount,
		ResourceStatefulSet,
		ResourceStorageClass,
	}
)

// KindToResource convert k8s kind name to k8s resource name.
// like:
//     KindToResource(Pod) --> pods
//     KindToResource(Ingress) -> ingresses
func KindToResource(kind string) string {
	for _, k := range allKind {
		if kind == k {
			switch kind {
			case KindIngress:
				return ResourceIngress
			case KindIngressClass:
				return ResourceIngressClass
			case KindNetworkPolicy:
				return ResourceNetworkPolicy
			case KindStorageClass:
				return ResourceStorageClass
			default:
				return strings.ToLower(kind) + "s"
			}
		}
	}

	return ""
}

// ResourceToKind convert k8s resource name to k8s kind name.
// like:
//     ResourceToKind(pods) -> Pod
//     ResourceToKind(networkpolicies) -> NetworkPolicy
func ResourceToKind(resource string) string {
	for _, r := range allResource {
		if resource == r {
			switch resource {
			case ResourceIngress:
				return KindIngress
			case ResourceIngressClass:
				return KindIngressClass
			case ResourceNetworkPolicy:
				return KindNetworkPolicy
			case ResourceStorageClass:
				return KindStorageClass
			default:
				return strings.TrimSuffix(strings.ToTitle(resource), "s")
			}
		}
	}

	return ""
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
