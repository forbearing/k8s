package dynamic

import (
	"fmt"
	"strings"

	. "github.com/forbearing/k8s/types"
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
func KindToResource(kind string) (string, error) {
	for _, k := range allKind {
		if kind == k {
			switch kind {
			case KindIngress:
				return ResourceIngress, nil
			case KindIngressClass:
				return ResourceIngressClass, nil
			case KindNetworkPolicy:
				return ResourceNetworkPolicy, nil
			case KindStorageClass:
				return ResourceStorageClass, nil
			default:
				return strings.ToLower(kind) + "s", nil
			}
		}
	}

	return "", fmt.Errorf("invalid kind: %s", kind)
}

// ResourceToKind convert k8s resource name to k8s kind name.
// like:
//     ResourceToKind(pods) -> Pod
//     ResourceToKind(networkpolicies) -> NetworkPolicy
func ResourceToKind(resource string) (string, error) {
	for _, r := range allResource {
		if resource == r {
			switch resource {
			case ResourceIngress:
				return KindIngress, nil
			case ResourceIngressClass:
				return KindIngressClass, nil
			case ResourceNetworkPolicy:
				return KindNetworkPolicy, nil
			case ResourceStorageClass:
				return KindStorageClass, nil
			default:
				return strings.TrimSuffix(strings.ToTitle(resource), "s"), nil
			}
		}
	}

	return "", fmt.Errorf("invalid resource: %s", resource)
}
