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

//type ListOptions struct {
//    TypeMeta             `json:",inline"`
//    LabelSelector        string               `json:"labelSelector,omitempty" protobuf:"bytes,1,opt,name=labelSelector"`
//    FieldSelector        string               `json:"fieldSelector,omitempty" protobuf:"bytes,2,opt,name=fieldSelector"`
//    Watch                bool                 `json:"watch,omitempty" protobuf:"varint,3,opt,name=watch"`
//    AllowWatchBookmarks  bool                 `json:"allowWatchBookmarks,omitempty" protobuf:"varint,9,opt,name=allowWatchBookmarks"`
//    ResourceVersion      string               `json:"resourceVersion,omitempty" protobuf:"bytes,4,opt,name=resourceVersion"`
//    ResourceVersionMatch ResourceVersionMatch `json:"resourceVersionMatch,omitempty" protobuf:"bytes,10,opt,name=resourceVersionMatch,casttype=ResourceVersionMatch"`
//    TimeoutSeconds       *int64               `json:"timeoutSeconds,omitempty" protobuf:"varint,5,opt,name=timeoutSeconds"`
//    Limit                int64                `json:"limit,omitempty" protobuf:"varint,7,opt,name=limit"`
//    Continue             string               `json:"continue,omitempty" protobuf:"bytes,8,opt,name=continue"`
//}
//type GetOptions struct {
//    TypeMeta        `json:",inline"`
//    ResourceVersion string `json:"resourceVersion,omitempty" protobuf:"bytes,1,opt,name=resourceVersion"`
//}
//type DeleteOptions struct {
//    TypeMeta           `json:",inline"`
//    GracePeriodSeconds *int64               `json:"gracePeriodSeconds,omitempty" protobuf:"varint,1,opt,name=gracePeriodSeconds"`
//    Preconditions      *Preconditions       `json:"preconditions,omitempty" protobuf:"bytes,2,opt,name=preconditions"`
//    OrphanDependents   *bool                `json:"orphanDependents,omitempty" protobuf:"varint,3,opt,name=orphanDependents"`
//    PropagationPolicy  *DeletionPropagation `json:"propagationPolicy,omitempty" protobuf:"varint,4,opt,name=propagationPolicy"`
//    DryRun             []string             `json:"dryRun,omitempty" protobuf:"bytes,5,rep,name=dryRun"`
//}
//type CreateOptions struct {
//    TypeMeta        `json:",inline"`
//    DryRun          []string `json:"dryRun,omitempty" protobuf:"bytes,1,rep,name=dryRun"`
//    FieldManager    string   `json:"fieldManager,omitempty" protobuf:"bytes,3,name=fieldManager"`
//    FieldValidation string   `json:"fieldValidation,omitempty" protobuf:"bytes,4,name=fieldValidation"`
//}
//type PatchOptions struct {
//    TypeMeta        `json:",inline"`
//    DryRun          []string `json:"dryRun,omitempty" protobuf:"bytes,1,rep,name=dryRun"`
//    Force           *bool    `json:"force,omitempty" protobuf:"varint,2,opt,name=force"`
//    FieldManager    string   `json:"fieldManager,omitempty" protobuf:"bytes,3,name=fieldManager"`
//    FieldValidation string   `json:"fieldValidation,omitempty" protobuf:"bytes,4,name=fieldValidation"`
//}
//type ApplyOptions struct {
//    TypeMeta     `json:",inline"`
//    DryRun       []string `json:"dryRun,omitempty" protobuf:"bytes,1,rep,name=dryRun"`
//    Force        bool     `json:"force" protobuf:"varint,2,opt,name=force"`
//    FieldManager string   `json:"fieldManager" protobuf:"bytes,3,name=fieldManager"`
//}
//type UpdateOptions struct {
//    TypeMeta        `json:",inline"`
//    DryRun          []string `json:"dryRun,omitempty" protobuf:"bytes,1,rep,name=dryRun"`
//    FieldManager    string   `json:"fieldManager,omitempty" protobuf:"bytes,2,name=fieldManager"`
//    FieldValidation string   `json:"fieldValidation,omitempty" protobuf:"bytes,3,name=fieldValidation"`
//}

type HandlerOptions struct {
	ListOptions   metav1.ListOptions
	GetOptions    metav1.GetOptions
	CreateOptions metav1.CreateOptions
	DeleteOptions metav1.DeleteOptions
	ApplyOptions  metav1.ApplyOptions
	UpdateOptions metav1.UpdateOptions
	PatchOptions  metav1.PatchOptions
}
