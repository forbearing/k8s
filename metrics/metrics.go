package metrics

import (
	"context"
	"encoding/json"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	v1beta1 "k8s.io/metrics/pkg/apis/metrics/v1beta1"
	metricsv "k8s.io/metrics/pkg/client/clientset/versioned"
	//metricsv "k8s.io/metrics/pkg/client/clientset/versioned"
	//metricsv1alpha1 "k8s.io/metrics/pkg/client/clientset/versioned/typed/metrics/v1alpha1"
	//metricsv1beta1 "k8s.io/metrics/pkg/client/clientset/versioned/typed/metrics/v1beta1"
	//v1beta1 "k8s.io/metrics/pkg/apis/metrics/v1beta1"
)

/*
reference:
	https://stackoverflow.com/questions/52029656/how-to-retrieve-kubernetes-metrics-via-client-go-and-golang

you can use following endpoints to retrieve the metrics as you want:
	Nodes: apis/metrics.k8s.io/v1beta1/nodes
	Pods: apis/metrics.k8s.io/v1beta1/pods
	Pods of default namespace: apis/metrics.k8s.io/v1beta1/namespaces/default/pods
	Specific Pod: /apis/metrics.k8s.io/v1beta1/namespaces/default/pods/<POD-NAME>

kubectl commnd:
	kubectl get --raw /apis/metrics.k8s.io/v1beta1/pods
*/

// node metrics object
type NodeMetrics struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Timestamp         string `json:"timpstamp"`
	Window            string `json:"window"`
	Usage             map[corev1.ResourceName]int64
}

// pod metrics object
type PodMetrics struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Timestamp         string `json:"timestamp"`
	Window            string `json:"window"`

	Containers []ContainerMetrics `json:"containers"`
}

// container metrics object
type ContainerMetrics struct {
	Name  string `json:"name,omitempty"`
	Usage map[corev1.ResourceName]int64
}

// metrics handler, used to get metrics from node or pod
type MetricsHandler struct {
	kubeconfig string
	namespace  string

	ctx       context.Context
	config    *rest.Config
	clientset *metricsv.Clientset
}

// Namespace return the the MetricsHandler working namespace.
// namespace will be required when query pods metrics
func (m *MetricsHandler) Namespace() string {
	return m.namespace
}

// DeepCopy copy a new MetricsHandler
func (in *MetricsHandler) DeepCopy() *MetricsHandler {
	out := new(MetricsHandler)
	out.kubeconfig = in.kubeconfig
	out.namespace = in.namespace
	out.ctx = in.ctx
	out.config = in.config
	out.clientset = in.clientset
	return out
}

// WithNamespace working with specific namespace
func (m *MetricsHandler) WithNamespace(namespace string) *MetricsHandler {
	metricsHandler := m.DeepCopy()
	metricsHandler.namespace = namespace
	return metricsHandler
}

// NewMetrics new a metrics handler from kubeconfig or in-cluster config
func NewMetrics(ctx context.Context, namespace, kubeconfig string) (metrics *MetricsHandler, err error) {
	var (
		config    *rest.Config
		clientset *metricsv.Clientset
	)

	if len(kubeconfig) != 0 {
		// create a rest config from kubeconfig
		if config, err = clientcmd.BuildConfigFromFlags("", kubeconfig); err != nil {
			return
		}
	} else {
		// create a rest config in-cluster config
		if config, err = rest.InClusterConfig(); err != nil {
			return
		}
	}

	// create a metrics clientset from rest config
	clientset, err = metricsv.NewForConfig(config)
	if err != nil {
		return
	}

	metrics = &MetricsHandler{}
	metrics.kubeconfig = kubeconfig
	metrics.namespace = namespace
	metrics.ctx = ctx
	metrics.config = config
	metrics.clientset = clientset
	return
}

// Pod query pod metrics by name
func (m *MetricsHandler) Pod(name string) (*PodMetrics, error) {
	podMetrics, err := m.clientset.MetricsV1beta1().PodMetricses(m.namespace).Get(m.ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return convertPodMetrics(podMetrics), nil
}

// Pods query multiple pod metrics by labels
func (m *MetricsHandler) Pods(label string) ([]PodMetrics, error) {
	//podMetricsList, err := m.clientset.MetricsV1beta1.PodMetricses(m.namespace).List()
	podMetricsList, err := m.clientset.MetricsV1beta1().PodMetricses(m.namespace).List(m.ctx, metav1.ListOptions{LabelSelector: label})
	if err != nil {
		return nil, err
	}

	pms := []PodMetrics{}
	for _, podMetrics := range podMetricsList.Items {
		pm := convertPodMetrics(&podMetrics)
		pms = append(pms, *pm)
	}
	return pms, nil
}

// PodRaw query pod metrics by name using REST API
func (m *MetricsHandler) PodRaw(name string) (*PodMetrics, error) {
	apiPath := "apis/metrics.k8s.io/v1beta1"
	if len(m.namespace) == 0 {
		return nil, fmt.Errorf("not set the namespace")
	}
	if len(name) == 0 {
		return nil, fmt.Errorf("pod name is empty")
	}
	apiPath = apiPath + "/namespaces/" + m.namespace + "/pods/" + name
	data, err := m.clientset.RESTClient().Get().AbsPath(apiPath).DoRaw(m.ctx)
	if err != nil {
		return nil, err
	}

	podMetrics := &v1beta1.PodMetrics{}
	if err = json.Unmarshal(data, podMetrics); err != nil {
		return nil, err
	}

	return convertPodMetrics(podMetrics), nil
}

// PodsRaw query the metrics of all pods in the namespace where the pod running, using REST API
func (m *MetricsHandler) PodsRaw() ([]PodMetrics, error) {
	apiPath := "apis/metrics.k8s.io/v1beta1"
	if len(m.namespace) == 0 {
		return nil, fmt.Errorf("not set the namespace")
	}
	apiPath = apiPath + "/namespaces/" + m.namespace + "/pods"
	data, err := m.clientset.RESTClient().Get().AbsPath(apiPath).DoRaw(m.ctx)
	if err != nil {
		return nil, err
	}

	podMetricsList := new(v1beta1.PodMetricsList)
	if err = json.Unmarshal(data, podMetricsList); err != nil {
		return nil, err
	}

	pms := []PodMetrics{}
	for _, podMetrics := range podMetricsList.Items {
		pm := convertPodMetrics(&podMetrics)
		pms = append(pms, *pm)
	}
	return pms, nil
}

// PodAllRaw query the metrics of all pods in the k8s cluster where the pod running, using REST API
func (m *MetricsHandler) PodAllRaw() ([]PodMetrics, error) {
	apiPath := "apis/metrics.k8s.io/v1beta1/pods"
	data, err := m.clientset.RESTClient().Get().AbsPath(apiPath).DoRaw(m.ctx)
	if err != nil {
		return nil, err
	}

	podMetricsList := new(v1beta1.PodMetricsList)
	if err = json.Unmarshal(data, podMetricsList); err != nil {
		return nil, err
	}

	pms := []PodMetrics{}
	for _, podMetrics := range podMetricsList.Items {
		pm := convertPodMetrics(&podMetrics)
		pms = append(pms, *pm)
	}
	return pms, nil
}

// Node query k8s node metrics by name
func (m *MetricsHandler) Node(name string) (*NodeMetrics, error) {
	nodeMetrics, err := m.clientset.MetricsV1beta1().NodeMetricses().Get(m.ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return convertNodeMetrics(nodeMetrics), nil
}

// Nodes query multiple k8s node metrics by label
func (m *MetricsHandler) Nodes(label string) ([]NodeMetrics, error) {
	nodeMetricsList, err := m.clientset.MetricsV1beta1().NodeMetricses().List(m.ctx, metav1.ListOptions{LabelSelector: label})
	if err != nil {
		return nil, err
	}

	nms := []NodeMetrics{}
	for _, nodeMetrics := range nodeMetricsList.Items {
		nm := convertNodeMetrics(&nodeMetrics)
		nms = append(nms, *nm)
	}
	return nms, nil
}

// NodeRaw query k8s node metrics by name using REST API
func (m *MetricsHandler) NodeRaw(name string) (*NodeMetrics, error) {
	apiPath := "apis/metrics.k8s.io/v1beta1"
	if len(name) == 0 {
		return nil, fmt.Errorf("k8s node hostname is empty")
	}
	apiPath = apiPath + "/nodes/" + name
	data, err := m.clientset.RESTClient().Get().AbsPath(apiPath).DoRaw(m.ctx)
	if err != nil {
		return nil, err
	}

	nodeMetrics := &v1beta1.NodeMetrics{}
	if err = json.Unmarshal(data, nodeMetrics); err != nil {
		return nil, err
	}

	return convertNodeMetrics(nodeMetrics), nil
}

// NodeAllRaw query all k8s node metrics using REST API
func (m *MetricsHandler) NodeAllRaw() ([]NodeMetrics, error) {
	apiPath := "apis/metrics.k8s.io/v1beta1/nodes"
	data, err := m.clientset.RESTClient().Get().AbsPath(apiPath).DoRaw(m.ctx)
	if err != nil {
		return nil, err
	}

	nodeMetricsList := new(v1beta1.NodeMetricsList)
	if err = json.Unmarshal(data, nodeMetricsList); err != nil {
		return nil, err
	}

	nms := []NodeMetrics{}
	for _, nodeMetrics := range nodeMetricsList.Items {
		nm := convertNodeMetrics(&nodeMetrics)
		nms = append(nms, *nm)
	}
	return nms, nil
}

// convertNodeMetrics convert *v1beta1.NodeMetrics to *NodeMetrics
func convertNodeMetrics(nodeMetrics *v1beta1.NodeMetrics) *NodeMetrics {
	nm := &NodeMetrics{}

	nm.TypeMeta = nodeMetrics.TypeMeta
	nm.ObjectMeta = nodeMetrics.ObjectMeta
	nm.Timestamp = nodeMetrics.Timestamp.Time.String()
	nm.Window = nodeMetrics.Window.Duration.String()
	nm.Usage = make(map[corev1.ResourceName]int64)
	for resourceName, resourceQuantity := range nodeMetrics.Usage {
		switch resourceName {
		case corev1.ResourceCPU:
			nm.Usage[resourceName] = resourceQuantity.MilliValue()
		case corev1.ResourceMemory:
			nm.Usage[resourceName] = resourceQuantity.Value()
		case corev1.ResourceStorage:
			nm.Usage[resourceName] = resourceQuantity.Value()
		case corev1.ResourceEphemeralStorage:
			nm.Usage[resourceName] = resourceQuantity.Value()
		default:
			nm.Usage[resourceName] = resourceQuantity.Value()
		}
	}

	return nm
}

// convertPodMetrics convert *v1beta1.PodMetrics to *PodMetrics
func convertPodMetrics(podMetrics *v1beta1.PodMetrics) *PodMetrics {
	pm := &PodMetrics{}

	pm.TypeMeta = podMetrics.TypeMeta
	pm.ObjectMeta = podMetrics.ObjectMeta
	pm.Timestamp = podMetrics.Timestamp.Time.String()
	pm.Window = podMetrics.Window.Duration.String()
	pm.Containers = convertContainersMetrics(podMetrics.Containers)

	return pm
}

// convertContainersMetrics convert []v1beta1.ContainerMetrics to []ContainerMetrics
func convertContainersMetrics(containersMetrics []v1beta1.ContainerMetrics) []ContainerMetrics {
	cms := []ContainerMetrics{}

	for _, metrics := range containersMetrics {
		cm := ContainerMetrics{}
		cm.Name = metrics.Name
		cm.Usage = make(map[corev1.ResourceName]int64)
		for resourceName, resourceQuantity := range metrics.Usage {
			switch resourceName {
			case corev1.ResourceCPU:
				cm.Usage[resourceName] = resourceQuantity.MilliValue()
			case corev1.ResourceMemory:
				cm.Usage[resourceName] = resourceQuantity.Value()
			case corev1.ResourceStorage:
				cm.Usage[resourceName] = resourceQuantity.Value()
			case corev1.ResourceEphemeralStorage:
				cm.Usage[resourceName] = resourceQuantity.Value()
			default:
				cm.Usage[resourceName] = resourceQuantity.Value()
			}
		}
		cms = append(cms, cm)
	}

	return cms
}
