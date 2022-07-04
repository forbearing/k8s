package typed

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// k8s resource kind
const (
	ResourceKindPod                   = "pod"
	ResourceKindDeployment            = "deployment"
	ResourceKindDaemonSet             = "daemonset"
	ResourceKindStatefulSet           = "statefulset"
	ResourceKindJob                   = "job"
	ResourceKindCronJob               = "cronjob"
	ResourceKindReplicaSet            = "replicaset"
	ResourceKindReplicationController = "replicationcontroller"

	ResourceKindPods                   = "pods"
	ResourceKindDeployments            = "deployments"
	ResourceKindDaemonSets             = "daemonsets"
	ResourceKindStatefulSets           = "statefulsets"
	ResourceKindJobs                   = "jobs"
	ResourceKindCronJobs               = "cronjobs"
	ResourceKindReplicaSets            = "replicasets"
	ResourceKindReplicationControllers = "replicationcontrollers"

	//CLUSTERROLEBINDING    = "clusterrolebinding"
	//CLUSTERROLE           = "clusterrole"
	//CONFIGMAP             = "configmap"
	//INGRES                = "ingress"
	//INGRESSCLASS          = "ingressclass"
	//NAMESPACE             = "namespace"
	//NETWORKPOLICY         = "networkpolicy"
	//NODE                  = "node"
	//PERSISTENTVOLUMECLAIM = "persistentvolumeclaim"
	//PERSISTENTVOLUME      = "persistentvolume"
	//ROLEBINDING           = "rolebinding"
	//ROLE                  = "role"
	//SECRET                = "secret"
	//SERVICEACCOUNT        = "serviceaccount"
	//SERVICE               = "service"
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

const (
	FieldManager = "client-go"
)
