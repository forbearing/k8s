package node

import (
	"fmt"
	"strings"
	"time"

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
	IPAddress          []string
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

// IsReady check whether the node is ready.
func (h *Handler) IsReady(name string) bool {
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

// These are the valid phases of node.
// Running, Pending, Terminated
func (h *Handler) GetPhase(object interface{}) (string, error) {
	switch val := object.(type) {
	case string:
		node, err := h.Get(val)
		if err != nil {
			return "", err
		}
		return string(node.Status.Phase), nil
	case *corev1.Node:
		return string(val.Status.Phase), nil
	case corev1.Node:
		return string(val.Status.Phase), nil
	default:
		return "", ErrInvalidToolsType
	}
}

//// GetStatus
//func (h *Handler) GetStatus(object interface{}) (string, error) {
//    switch val := object.(type) {
//    case string:
//        node, err := h.Get(val)
//        if err != nil {
//            return "", err
//        }
//        return string(node.Status.Phase), nil
//    case *corev1.Node:
//        return string(val.Status.Phase), nil
//    case corev1.Node:
//        return string(val.Status.Phase), nil
//    default:
//        return "", ErrInvalidToolsType
//    }
//}

// GetHostname returns the node ip
func (h *Handler) GetHostname(object interface{}) ([]string, error) {

	switch val := object.(type) {
	case string:
		node, err := h.Get(val)
		if err != nil {
			return nil, err
		}
		return h.getHostname(node), nil
	case *corev1.Node:
		return h.getHostname(val), nil
	case corev1.Node:
		return h.getHostname(&val), nil
	default:
		return nil, ErrInvalidToolsType
	}
}
func (h *Handler) getHostname(node *corev1.Node) []string {
	var al []string
	for _, addr := range node.Status.Addresses {
		if addr.Type == corev1.NodeHostName {
			al = append(al, addr.Address)
		}
	}
	return al
}

// GetInternalIP returns the node ip
func (h *Handler) GetInternalIP(object interface{}) ([]string, error) {

	switch val := object.(type) {
	case string:
		node, err := h.Get(val)
		if err != nil {
			return nil, err
		}
		return h.getInternalIP(node), nil
	case *corev1.Node:
		return h.getInternalIP(val), nil
	case corev1.Node:
		return h.getInternalIP(&val), nil
	default:
		return nil, ErrInvalidToolsType
	}
}
func (h *Handler) getInternalIP(node *corev1.Node) []string {
	var al []string
	for _, addr := range node.Status.Addresses {
		if addr.Type == corev1.NodeInternalIP {
			al = append(al, addr.Address)
		}
	}
	return al
}

// GetExternalIP returns the node ip
func (h *Handler) GetExternalIP(object interface{}) ([]string, error) {

	switch val := object.(type) {
	case string:
		node, err := h.Get(val)
		if err != nil {
			return nil, err
		}
		return h.getExternalIP(node), nil
	case *corev1.Node:
		return h.getExternalIP(val), nil
	case corev1.Node:
		return h.getExternalIP(&val), nil
	default:
		return nil, ErrInvalidToolsType
	}
}
func (h *Handler) getExternalIP(node *corev1.Node) []string {
	var al []string
	for _, addr := range node.Status.Addresses {
		if addr.Type == corev1.NodeExternalIP {
			al = append(al, addr.Address)
		}
	}
	return al
}

// GetInternalDNS returns the node ip
func (h *Handler) GetInternalDNS(object interface{}) ([]string, error) {

	switch val := object.(type) {
	case string:
		node, err := h.Get(val)
		if err != nil {
			return nil, err
		}
		return h.getInternalDNS(node), nil
	case *corev1.Node:
		return h.getInternalDNS(val), nil
	case corev1.Node:
		return h.getInternalDNS(&val), nil
	default:
		return nil, ErrInvalidToolsType
	}
}
func (h *Handler) getInternalDNS(node *corev1.Node) []string {
	var al []string
	for _, addr := range node.Status.Addresses {
		if addr.Type == corev1.NodeInternalDNS {
			al = append(al, addr.Address)
		}
	}
	return al
}

// GetExternaDNS returns the node ip
func (h *Handler) GetExternaDNS(object interface{}) ([]string, error) {

	switch val := object.(type) {
	case string:
		node, err := h.Get(val)
		if err != nil {
			return nil, err
		}
		return h.getExternalDNS(node), nil
	case *corev1.Node:
		return h.getExternalDNS(val), nil
	case corev1.Node:
		return h.getExternalDNS(&val), nil
	default:
		return nil, ErrInvalidToolsType
	}
}
func (h *Handler) getExternalDNS(node *corev1.Node) []string {
	var al []string
	for _, addr := range node.Status.Addresses {
		if addr.Type == corev1.NodeExternalDNS {
			al = append(al, addr.Address)
		}
	}
	return al
}

// IsMaster check whether the node is master.
func (h *Handler) IsMaster(object interface{}) bool {
	roles, _ := h.GetRoles(object)
	for _, role := range roles {
		if strings.ToLower(role) == NodeRoleMaster {
			return true
		}
	}
	return false
}

// IsControlPlane check whether the node is control-plane.
func (h *Handler) IsControlPlane(object interface{}) bool {
	roles, _ := h.GetRoles(object)
	for _, role := range roles {
		if strings.ToLower(role) == NodeRoleControlPlane {
			return true
		}
	}
	return false
}

// GetRoles returns the roles of a given node.
// The roles are determined by looking for:
//   node-role.kubernetes.io/<role>=""
//   kubernetes.io/role="<role>"
func (h *Handler) GetRoles(object interface{}) ([]string, error) {
	switch val := object.(type) {
	case string:
		node, err := h.Get(val)
		if err != nil {
			return nil, err
		}
		return h.getRoles(node), nil
	case *corev1.Node:
		return h.getRoles(val), nil
	case corev1.Node:
		return h.getRoles(&val), nil
	default:
		return nil, ErrInvalidToolsType
	}
}
func (h *Handler) getRoles(node *corev1.Node) []string {
	roles := sets.NewString()
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

// GetPods get all pods running in the node
func (h *Handler) GetPods(object interface{}) ([]*corev1.Pod, error) {
	switch val := object.(type) {
	case string:
		node, err := h.Get(val)
		if err != nil {
			return nil, err
		}
		return h.getPods(node)
	case *corev1.Node:
		return h.getPods(val)
	case corev1.Node:
		return h.getPods(&val)
	default:
		return nil, ErrInvalidToolsType
	}
}
func (h *Handler) getPods(node *corev1.Node) ([]*corev1.Pod, error) {
	// ParseSelector takes a string representing a selector and returns an
	// object suitable for matching, or an error.
	fieldSelector, err := fields.ParseSelector(fmt.Sprintf("spec.nodeName=%s", node.Name))
	if err != nil {
		return nil, err
	}

	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.FieldSelector = fieldSelector.String()
	podList, err := h.clientset.CoreV1().Pods(metav1.NamespaceAll).List(h.ctx, *listOptions)
	if err != nil {
		return nil, err
	}

	var pl []*corev1.Pod
	for i := range podList.Items {
		pl = append(pl, &podList.Items[i])
	}
	return pl, nil
}

// GetCIDR get the node podCIDR
func (h *Handler) GetCIDR(object interface{}) (string, error) {
	switch val := object.(type) {
	case string:
		node, err := h.Get(val)
		if err != nil {
			return "", err
		}
		return node.Spec.PodCIDR, nil
	case *corev1.Node:
		return val.Spec.PodCIDR, nil
	case corev1.Node:
		return val.Spec.PodCIDR, nil
	default:
		return "", ErrInvalidToolsType
	}
}

// GetCIDRs get the node podCIDRs
func (h *Handler) GetCIDRs(object interface{}) ([]string, error) {
	switch val := object.(type) {
	case string:
		node, err := h.Get(val)
		if err != nil {
			return nil, err
		}
		return node.Spec.PodCIDRs, nil
	case *corev1.Node:
		return val.Spec.PodCIDRs, nil
	case corev1.Node:
		return val.Spec.PodCIDRs, nil
	default:
		return nil, ErrInvalidToolsType
	}
}

// GetNodeInfo get given node info.
func (h *Handler) GetNodeInfo(object interface{}) (*NodeInfo, error) {
	switch val := object.(type) {
	case string:
		node, err := h.Get(val)
		if err != nil {
			return nil, err
		}
		return h.getNodeInfo(node), nil
	case *corev1.Node:
		return h.getNodeInfo(val), nil
	case corev1.Node:
		return h.getNodeInfo(&val), nil
	default:
		return nil, ErrInvalidToolsType
	}
}

func (h *Handler) getNodeInfo(node *corev1.Node) *NodeInfo {
	var nodeInfo NodeInfo

	nodeInfo.Hostname = node.ObjectMeta.Name
	nodeInfo.IPAddress, _ = h.GetInternalIP(nodeInfo.Hostname)
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

	return &nodeInfo
}

// GetMasterInfo
func (h *Handler) GetMasterInfo() ([]NodeInfo, error) {
	// TODO:
	// 1.考虑标签格式部位  "node-role.kubernetes.io/master=" 的情况
	// 2.测试下 List 获取的结果为空是, err 是否为 nil.
	masterList, err := h.ListByLabel(LabelNodeRolePrefix + "master")
	if err != nil {
		return nil, err
	}

	var nodeinfoList []NodeInfo
	for _, master := range masterList {
		nodeInfo, err := h.GetNodeInfo(master)
		if err == nil {
			nodeinfoList = append(nodeinfoList, *nodeInfo)
		}
	}

	return nodeinfoList, nil
}

// GetWorkerInfo
func (h *Handler) GetWorkerInfo() ([]NodeInfo, error) {
	// TODO:
	// 1.考虑标签格式部位  "node-role.kubernetes.io/master=" 的情况
	// 2.测试下 List 获取的结果为空是, err 是否为 nil.
	masterList, err := h.ListByLabel("!" + LabelNodeRolePrefix + "master")
	if err != nil {
		return nil, err
	}

	var nodeinfoList []NodeInfo
	for _, master := range masterList {
		nodeInfo, err := h.GetNodeInfo(master)
		if err == nil {
			nodeinfoList = append(nodeinfoList, *nodeInfo)
		}
	}

	return nodeinfoList, nil
}

// GetAge
func (h *Handler) GetAge(object interface{}) (time.Duration, error) {
	switch val := object.(type) {
	case string:
		node, err := h.Get(val)
		if err != nil {
			return time.Duration(0), err
		}
		return time.Now().Sub(node.CreationTimestamp.Time), nil
	case *corev1.Node:
		return time.Now().Sub(val.CreationTimestamp.Time), nil
	case corev1.Node:
		return time.Now().Sub(val.CreationTimestamp.Time), nil
	default:
		return time.Duration(0), ErrInvalidToolsType
	}
}
