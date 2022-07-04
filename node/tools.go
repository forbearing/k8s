package node

import (
	"fmt"
	"strings"

	"github.com/forbearing/k8s/pod"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/util/sets"
)

const (
	// LabelNodeRolePrefix is a label prefix for node roles
	// It's copied over to here until it's merged in core: https://github.com/kubernetes/kubernetes/pull/39112
	LabelNodeRolePrefix = "node-role.kubernetes.io/"

	// LabelNodeRole specifies the role of a node
	LabelNodeRole = "kubernetes.io/role"

	//LabelNodeRoleMaster       = "kubernetes.io/role=master"
	//LabelNodeRoleControlPlane = "kubernetes.io/role=control-plane"
	//LabelNodeRoleWorker       = "!kubernetes.io/role=master"

	NodeRoleMaster       = "master"
	NodeRoleControlPlane = "control-plane"
)

type NodeStatus struct {
	Status  string
	Message string
	Reason  string
}
type NodeInfo struct {
	Hostname           string
	IPAddress          string
	AllocatableCpu     string
	AllocatableMemory  string
	AllocatableStorage string
	TotalCpu           string
	TotalMemory        string
	TotalStorage       string

	Architecture            string
	BootID                  string
	ContainerRuntimeVersion string
	KernelVersion           string
	KubeProxyVersion        string
	KubeletVersion          string
	MachineID               string
	OperatingSystem         string
	OSImage                 string
	SystemUUID              string
}

// check if the node status is ready
func (h *Handler) IsReady(name string) bool {
	// get *corev1.Node
	node, err := h.Get(name)
	if err != nil {
		return false
	}
	for _, cond := range node.Status.Conditions {
		if cond.Status == corev1.ConditionTrue && cond.Type == corev1.NodeReady {
			return true
		}
	}
	return false
}

// check if the node is master
func (h *Handler) IsMaster(name string) bool {
	roles := h.GetRoles(name)
	for _, role := range roles {
		if strings.ToLower(role) == NodeRoleMaster {
			return true
		}
	}
	return false
}

// check if the node is control-plane
func (h *Handler) IsControlPlane(name string) bool {
	roles := h.GetRoles(name)
	for _, role := range roles {
		if strings.ToLower(role) == NodeRoleControlPlane {
			return true
		}
	}
	return false
}

// get the node status
func (h *Handler) GetStatus(name string) *NodeStatus {
	nodeStatus := &NodeStatus{
		Message: "Unknow",
		Reason:  "Unknow",
		Status:  string(corev1.ConditionUnknown),
	}

	// get *corev1.Node
	node, err := h.Get(name)
	if err != nil {
		return nodeStatus
	}

	for _, cond := range node.Status.Conditions {
		if cond.Type == corev1.NodeReady {
			nodeStatus = &NodeStatus{
				Message: cond.Message,
				Reason:  cond.Reason,
				Status:  string(cond.Status),
			}
		}
	}

	return nodeStatus
}

// GetRoles returns the roles of a given node.
// The roles are determined by looking for:
//   node-role.kubernetes.io/<role>=""
//   kubernetes.io/role="<role>"
func (h *Handler) GetRoles(name string) []string {
	roles := sets.NewString()

	// get *corev1.Node
	node, err := h.Get(name)
	if err != nil {
		return roles.List()
	}

	for label, value := range node.Labels {
		switch {
		case strings.HasPrefix(label, LabelNodeRolePrefix):
			if role := strings.TrimPrefix(label, LabelNodeRolePrefix); len(role) > 0 {
				roles.Insert(role)
			}
		case label == LabelNodeRole && len(value) > 0:
			roles.Insert(value)
		}
	}

	return roles.List()
}

// get all pods in the node
func (h *Handler) GetPods(name string) (*corev1.PodList, error) {
	// ParseSelector takes a string representing a selector and returns an
	// object suitable for matching, or an error.
	fieldSelector, err := fields.ParseSelector(fmt.Sprintf("spec.nodeName=%s", name))
	if err != nil {
		return nil, err
	}

	podHandler, err := pod.New(h.ctx, "", h.kubeconfig)
	if err != nil {
		return nil, err
	}
	podHandler.Options.ListOptions = metav1.ListOptions{FieldSelector: fieldSelector.String()}
	//podHandler.SetNamespace(metav1.NamespaceAll)
	//return podHandler.List("")
	return podHandler.WithNamespace(metav1.NamespaceAll).List("")
}

// get not terminated pod in the node.
func (h *Handler) GetNonTerminatedPods(name string) (*corev1.PodList, error) {
	// PodSucceeded 表示 containers 成功退出, pod 终止
	// PodSucceeded 表示 containers 失败退出, pod 也终止
	// PodPending, PodRunning, PodUnknown 都表示 pod 正在运行
	selector := fmt.Sprintf("spec.nodeName=%s,status.phase!=%s,status.phase!=%s",
		name, string(corev1.PodSucceeded), string(corev1.PodFailed))
	// ParseSelector takes a string representing a selector and returns an
	// object suitable for matching, or an error.
	fieldSelector, err := fields.ParseSelector(selector)
	if err != nil {
		return nil, err
	}
	podHandler, err := pod.New(h.ctx, "", h.kubeconfig)
	if err != nil {
		return nil, err
	}
	podHandler.Options.ListOptions = metav1.ListOptions{FieldSelector: fieldSelector.String()}
	return podHandler.WithNamespace(metav1.NamespaceAll).List("")
}

// get the node ip
func (h *Handler) GetIP(name string) (ip string, err error) {
	node, err := h.Get(name)
	if err != nil {
		return
	}
	for _, address := range node.Status.Addresses {
		if address.Type == corev1.NodeInternalIP {
			ip = address.Address
		}
	}
	return
}

// get the node hostname
func (h *Handler) GetHostname(name string) (hostname string, err error) {
	node, err := h.Get(name)
	if err != nil {
		return
	}
	for _, address := range node.Status.Addresses {
		if address.Type == corev1.NodeHostName {
			hostname = address.Address
		}
	}
	return
}

// get the node podCIDR
func (h *Handler) GetCIDR(name string) (string, error) {
	node, err := h.Get(name)
	if err != nil {
		return "", err
	}
	return node.Spec.PodCIDR, nil
}

// get the node podCIDRs
func (h *Handler) GetCIDRs(name string) ([]string, error) {
	node, err := h.Get(name)
	if err != nil {
		return nil, err
	}
	return node.Spec.PodCIDRs, nil
}

// get all master node info
func (h *Handler) GetMasterInfo() ([]NodeInfo, error) {
	var nodeInfo NodeInfo
	var nodeInfoList []NodeInfo

	masterNodes, err := h.List(LabelNodeRolePrefix + "master")
	if err != nil {
		return nil, err
	}
	for _, node := range masterNodes.Items {
		nodeInfo.Hostname = node.ObjectMeta.Name
		nodeInfo.IPAddress, _ = h.GetIP(nodeInfo.Hostname)
		nodeInfo.AllocatableCpu = node.Status.Allocatable.Cpu().String()
		nodeInfo.AllocatableMemory = node.Status.Allocatable.Memory().String()
		nodeInfo.AllocatableStorage = node.Status.Allocatable.StorageEphemeral().String()
		nodeInfo.Architecture = node.Status.NodeInfo.Architecture
		nodeInfo.TotalCpu = node.Status.Capacity.Cpu().String()
		nodeInfo.TotalMemory = node.Status.Capacity.Memory().String()
		nodeInfo.TotalStorage = node.Status.Capacity.StorageEphemeral().String()
		nodeInfo.BootID = node.Status.NodeInfo.BootID
		nodeInfo.ContainerRuntimeVersion = node.Status.NodeInfo.ContainerRuntimeVersion
		nodeInfo.KernelVersion = node.Status.NodeInfo.KernelVersion
		nodeInfo.KubeProxyVersion = node.Status.NodeInfo.KubeProxyVersion
		nodeInfo.KubeletVersion = node.Status.NodeInfo.KubeletVersion
		nodeInfo.MachineID = node.Status.NodeInfo.MachineID
		nodeInfo.OperatingSystem = node.Status.NodeInfo.OperatingSystem
		nodeInfo.OSImage = node.Status.NodeInfo.OSImage
		nodeInfo.SystemUUID = node.Status.NodeInfo.SystemUUID
		// map 的 key 就是 node.ObjectMeta.Name, 即 k8s 节点的 ip 地址
		nodeInfoList = append(nodeInfoList, nodeInfo)
	}

	return nodeInfoList, nil
}

// get all worker node info
func (h *Handler) GetWorkerInfo() ([]NodeInfo, error) {
	var nodeInfo NodeInfo
	var nodeInfoList []NodeInfo

	workerNodes, err := h.List("!" + LabelNodeRolePrefix + "master")
	if err != nil {
		return nil, err
	}
	for _, node := range workerNodes.Items {
		nodeInfo.Hostname = node.ObjectMeta.Name
		nodeInfo.IPAddress, _ = h.GetIP(nodeInfo.Hostname)
		nodeInfo.AllocatableCpu = node.Status.Allocatable.Cpu().String()
		nodeInfo.AllocatableMemory = node.Status.Allocatable.Memory().String()
		nodeInfo.AllocatableStorage = node.Status.Allocatable.StorageEphemeral().String()
		nodeInfo.Architecture = node.Status.NodeInfo.Architecture
		nodeInfo.TotalCpu = node.Status.Capacity.Cpu().String()
		nodeInfo.TotalMemory = node.Status.Capacity.Memory().String()
		nodeInfo.TotalStorage = node.Status.Capacity.StorageEphemeral().String()
		nodeInfo.BootID = node.Status.NodeInfo.BootID
		nodeInfo.ContainerRuntimeVersion = node.Status.NodeInfo.ContainerRuntimeVersion
		nodeInfo.KernelVersion = node.Status.NodeInfo.KernelVersion
		nodeInfo.KubeProxyVersion = node.Status.NodeInfo.KubeProxyVersion
		nodeInfo.KubeletVersion = node.Status.NodeInfo.KubeletVersion
		nodeInfo.MachineID = node.Status.NodeInfo.MachineID
		nodeInfo.OperatingSystem = node.Status.NodeInfo.OperatingSystem
		nodeInfo.OSImage = node.Status.NodeInfo.OSImage
		nodeInfo.SystemUUID = node.Status.NodeInfo.SystemUUID
		// map 的 key 就是 node.ObjectMeta.Name, 即 k8s 节点的 ip 地址
		nodeInfoList = append(nodeInfoList, nodeInfo)
	}
	return nodeInfoList, nil
}

// get all k8s node info
func (h *Handler) GetAllInfo() ([]NodeInfo, error) {
	var nodeInfoList []NodeInfo
	masterInfo, err := h.GetMasterInfo()
	if err != nil {
		return nil, err
	}
	workerInfo, err := h.GetWorkerInfo()
	if err != nil {
		return nil, err
	}

	for _, info := range masterInfo {
		nodeInfoList = append(nodeInfoList, info)
	}
	for _, info := range workerInfo {
		nodeInfoList = append(nodeInfoList, info)
	}

	return nodeInfoList, nil
}
