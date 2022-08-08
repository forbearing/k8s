package types

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// k8s resource name
const (
	ResourceClusterRole           = "clusterroles"
	ResourceClusterRoleBinding    = "clusterrolebindings"
	ResourceConfigMap             = "configmaps"
	ResourceCronJob               = "cronjobs"
	ResourceDaemonSet             = "daemonsets"
	ResourceDeployment            = "deployments"
	ResourceIngress               = "ingresses"
	ResourceIngressClass          = "ingressclasses"
	ResourceJob                   = "jobs"
	ResourceNamespace             = "namespaces"
	ResourceNetworkPolicy         = "networkpolicies"
	ResourceNode                  = "nodes"
	ResourcePersistentVolume      = "persistentvolumes"
	ResourcePersistentVolumeClaim = "persistentvolumeclaims"
	ResourcePod                   = "pods"
	ResourceReplicaSet            = "replicasets"
	ResourceReplicationController = "replicationcontrollers"
	ResourceRole                  = "roles"
	ResourceRoleBinding           = "rolebindings"
	ResourceSecret                = "secrets"
	ResourceService               = "services"
	ResourceServiceAccount        = "serviceaccounts"
	ResourceStatefulSet           = "statefulsets"
	ResourceStorageClass          = "storageclasses"
)

// k8s resource kind
const (
	KindClusterRole           = "ClusterRole"
	KindClusterRoleBinding    = "ClusterRoleBinding"
	KindConfigMap             = "ConfigMap"
	KindCronJob               = "CronJob"
	KindDaemonSet             = "DaemonSet"
	KindDeployment            = "Deployment"
	KindIngress               = "Ingress"
	KindIngressClass          = "IngressClass"
	KindJob                   = "Job"
	KindNamespace             = "Namespace"
	KindNetworkPolicy         = "NetworkPolicy"
	KindNode                  = "Node"
	KindPersistentVolume      = "PersistentVolume"
	KindPersistentVolumeClaim = "PersistentVolumeClaim"
	KindPod                   = "Pod"
	KindReplicaSet            = "ReplicaSet"
	KindReplicationController = "ReplicationController"
	KindRole                  = "Role"
	KindRoleBinding           = "RoleBinding"
	KindSecret                = "Secret"
	KindService               = "Service"
	KindServiceAccount        = "ServiceAccount"
	KindStatefulSet           = "StatefulSet"
	KindStorageClass          = "StorageClass"
)

type HandlerOptions struct {
	ListOptions   metav1.ListOptions
	GetOptions    metav1.GetOptions
	CreateOptions metav1.CreateOptions
	DeleteOptions metav1.DeleteOptions
	ApplyOptions  metav1.ApplyOptions
	UpdateOptions metav1.UpdateOptions
	PatchOptions  metav1.PatchOptions
}
