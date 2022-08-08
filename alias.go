package k8s

import (
	"github.com/forbearing/k8s/clusterrole"
	"github.com/forbearing/k8s/clusterrolebinding"
	"github.com/forbearing/k8s/configmap"
	"github.com/forbearing/k8s/cronjob"
	"github.com/forbearing/k8s/daemonset"
	"github.com/forbearing/k8s/deployment"
	"github.com/forbearing/k8s/ingress"
	"github.com/forbearing/k8s/ingressclass"
	"github.com/forbearing/k8s/job"
	"github.com/forbearing/k8s/namespace"
	"github.com/forbearing/k8s/networkpolicy"
	"github.com/forbearing/k8s/node"
	"github.com/forbearing/k8s/persistentvolume"
	"github.com/forbearing/k8s/persistentvolumeclaim"
	"github.com/forbearing/k8s/pod"
	"github.com/forbearing/k8s/replicaset"
	"github.com/forbearing/k8s/replicationcontroller"
	"github.com/forbearing/k8s/role"
	"github.com/forbearing/k8s/rolebinding"
	"github.com/forbearing/k8s/secret"
	"github.com/forbearing/k8s/service"
	"github.com/forbearing/k8s/serviceaccount"
	"github.com/forbearing/k8s/statefulset"
	"github.com/forbearing/k8s/storageclass"
	"github.com/forbearing/k8s/tools/signals"
)

type ClusterRoleHandler = clusterrole.Handler
type ClusterRoleBindingHandler = clusterrolebinding.Handler
type ConfigMapHandler = configmap.Handler
type CronJobHandler = cronjob.Handler
type DaemonSetHandler = daemonset.Handler
type DeploymentHandler = deployment.Handler
type IngressHandler = ingress.Handler
type IngressClassHandler = ingressclass.Handler
type JobHandler = job.Handler
type NamespaceHandler = namespace.Handler
type NetworkPolicyHandler = networkpolicy.Handler
type NodeHandler = node.Handler
type PersistentVolumeHandler = persistentvolume.Handler
type PersistentVolumeClaimHandler = persistentvolumeclaim.Handler
type PodHandler = pod.Handler
type ReplicaSetHandler = replicaset.Handler
type ReplicationControllerHandler = replicationcontroller.Handler
type RoleHandler = role.Handler
type RoleBindingHandler = rolebinding.Handler
type SecretHandler = secret.Handler
type ServiceHandler = service.Handler
type ServiceAccountHandler = serviceaccount.Handler
type StatefulSetHandler = statefulset.Handler
type StorageClassHandler = storageclass.Handler

var (
	NewClusterRole           = clusterrole.New
	NewClusterRoleBinding    = clusterrolebinding.New
	NewConfigMap             = configmap.New
	NewCronJob               = cronjob.New
	NewDaemonSet             = daemonset.New
	NewDeployment            = deployment.New
	NewIngress               = ingress.New
	NewIngressClass          = ingressclass.New
	NewJob                   = job.New
	NewNamespace             = namespace.New
	NewNetworkPolicy         = networkpolicy.New
	NewNode                  = node.New
	NewPersistentVolume      = persistentvolume.New
	NewPersistentVolumeClaim = persistentvolumeclaim.New
	NewPod                   = pod.New
	NewReplicaSet            = replicaset.New
	NewReplicationController = replicationcontroller.New
	NewRole                  = role.New
	NewRoleBinding           = rolebinding.New
	NewSecret                = secret.New
	NewService               = service.New
	NewServiceAccount        = serviceaccount.New
	NewStatefulSet           = statefulset.New
	NewStorageClass          = storageclass.New
)

var (
	SetupSignalContext = signals.SetupSignalContext
	SetupSignalChannel = signals.SetupSignalChannel
)
