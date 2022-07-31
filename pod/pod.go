package pod

// TODO:
// 1. 加了 Options 参数之后可能要修改的东西: List, Watch, WaitReady
// 2. 精简代码, 比如返回值
import (
	"context"
	"net/http"
	"sync"

	//_ "k8s.io/client-go/kubernetes/types/core/v1"

	"github.com/forbearing/k8s/types"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Handler struct {
	kubeconfig string
	namespace  string

	ctx             context.Context
	config          *rest.Config
	httpClient      *http.Client
	restClient      *rest.RESTClient
	clientset       *kubernetes.Clientset
	dynamicClient   dynamic.Interface
	discoveryClient *discovery.DiscoveryClient
	informerFactory informers.SharedInformerFactory
	client          typedcorev1.PodInterface

	Options *types.HandlerOptions

	l sync.RWMutex
}

// NewOrDie simply call New() to get a pod handler.
// panic if there is any error occurs.
func NewOrDie(ctx context.Context, kubeconfig, namespace string) *Handler {
	handler, err := New(ctx, kubeconfig, namespace)
	if err != nil {
		panic(err)
	}
	return handler
}

// New returns a pod handler from kubeconfig or in-cluster config.
func New(ctx context.Context, kubeconfig, namespace string) (handler *Handler, err error) {
	var (
		config          *rest.Config
		httpClient      *http.Client
		restClient      *rest.RESTClient
		clientset       *kubernetes.Clientset
		dynamicClient   dynamic.Interface
		discoveryClient *discovery.DiscoveryClient
		informerFactory informers.SharedInformerFactory
	)
	handler = &Handler{}

	// create rest config
	if len(kubeconfig) != 0 {
		// use the current context in kubeconfig
		if config, err = clientcmd.BuildConfigFromFlags("", kubeconfig); err != nil {
			return nil, err
		}
	} else {
		// creates the in-cluster config
		if config, err = rest.InClusterConfig(); err != nil {
			return nil, err
		}
	}

	// setup APIPath, GroupVersion and NegotiatedSerializer before initializing a RESTClient
	config.APIPath = "api"
	config.GroupVersion = &corev1.SchemeGroupVersion
	//config.GroupVersion = &schema.GroupVersion{Group: "", Version: "v1"}
	config.NegotiatedSerializer = scheme.Codecs
	//config.UserAgent = rest.DefaultKubernetesUserAgent()

	//// k8s cluster endpoint, eg: https://10.250.16.10:8443
	//config.Host = "127.0.0.1"
	//config.ContentConfig = rest.ContentConfig{
	//    GroupVersion:         &corev1.SchemeGroupVersion,
	//    NegotiatedSerializer: scheme.Codecs,
	//}

	// create a http client for the given config.
	if httpClient, err = rest.HTTPClientFor(config); err != nil {
		return nil, err
	}
	// create a RESTClient for the given config and http client.
	if restClient, err = rest.RESTClientForConfigAndClient(config, httpClient); err != nil {
		return nil, err
	}
	// create a Clientset for the given config and http client.
	if clientset, err = kubernetes.NewForConfigAndClient(config, httpClient); err != nil {
		return nil, err
	}
	// create a dynamic client for the given config and http client.
	if dynamicClient, err = dynamic.NewForConfigAndClient(config, httpClient); err != nil {
		return nil, err
	}
	// create a DiscoveryClient for the given config and http client.
	if discoveryClient, err = discovery.NewDiscoveryClientForConfigAndClient(config, httpClient); err != nil {
		return nil, err
	}
	// create a sharedInformerFactory for all namespaces.
	// NewSharedInformerFactory constructs a new instance of sharedInformerFactory for all namespaces.
	informerFactory = informers.NewSharedInformerFactory(clientset, 0)

	if len(namespace) == 0 {
		namespace = metav1.NamespaceDefault
	}
	handler.kubeconfig = kubeconfig
	handler.namespace = namespace
	handler.ctx = ctx
	handler.config = config
	handler.httpClient = httpClient
	handler.restClient = restClient
	handler.clientset = clientset
	handler.dynamicClient = dynamicClient
	handler.discoveryClient = discoveryClient
	handler.informerFactory = informerFactory
	handler.client = clientset.CoreV1().Pods(namespace)
	handler.Options = &types.HandlerOptions{}

	return handler, nil
}
func (p *Handler) Namespace() string {
	return p.namespace
}
func (in *Handler) DeepCopy() *Handler {
	if in == nil {
		return nil
	}
	out := new(Handler)

	out.kubeconfig = in.kubeconfig
	out.namespace = in.namespace

	out.ctx = in.ctx
	out.config = in.config
	out.httpClient = in.httpClient
	out.restClient = in.restClient
	out.clientset = in.clientset
	out.dynamicClient = in.dynamicClient
	out.discoveryClient = in.discoveryClient
	out.informerFactory = in.informerFactory
	out.client = in.client

	out.Options = &types.HandlerOptions{}
	out.Options.ListOptions = *in.Options.ListOptions.DeepCopy()
	out.Options.GetOptions = *in.Options.GetOptions.DeepCopy()
	out.Options.CreateOptions = *in.Options.CreateOptions.DeepCopy()
	out.Options.UpdateOptions = *in.Options.UpdateOptions.DeepCopy()
	out.Options.PatchOptions = *in.Options.PatchOptions.DeepCopy()
	out.Options.ApplyOptions = *in.Options.ApplyOptions.DeepCopy()

	return out
}

func (p *Handler) resetNamespace(namespace string) {
	p.l.Lock()
	defer p.l.Unlock()
	p.namespace = namespace
}
func (p *Handler) WithNamespace(namespace string) *Handler {
	handler := p.DeepCopy()
	handler.resetNamespace(namespace)
	return handler
}
func (p *Handler) WithDryRun() *Handler {
	handler := p.DeepCopy()
	handler.Options.CreateOptions.DryRun = []string{metav1.DryRunAll}
	handler.Options.UpdateOptions.DryRun = []string{metav1.DryRunAll}
	handler.Options.DeleteOptions.DryRun = []string{metav1.DryRunAll}
	handler.Options.PatchOptions.DryRun = []string{metav1.DryRunAll}
	handler.Options.ApplyOptions.DryRun = []string{metav1.DryRunAll}
	return handler
}
func (p *Handler) SetLimit(limit int64) {
	p.l.Lock()
	defer p.l.Unlock()
	p.Options.ListOptions.Limit = limit
}
func (p *Handler) SetTimeout(timeout int64) {
	p.l.Lock()
	defer p.l.Unlock()
	p.Options.ListOptions.TimeoutSeconds = &timeout
}
func (p *Handler) SetForceDelete(force bool) {
	p.l.Lock()
	defer p.l.Unlock()
	if force {
		gracePeriodSeconds := int64(0)
		p.Options.DeleteOptions.GracePeriodSeconds = &gracePeriodSeconds
	} else {
		p.Options.DeleteOptions = metav1.DeleteOptions{}
	}
}

// RESTConfig returns underlying rest config.
func (h *Handler) RESTConfig() *rest.Config {
	return h.config
}

// RESTClient returns underlying rest client.
func (h *Handler) RESTClient() *rest.RESTClient {
	return h.restClient
}

// Clientset returns underlying clientset.
func (h *Handler) Clientset() *kubernetes.Clientset {
	return h.clientset
}

// DynamicClient returns underlying dynamic client.
func (h *Handler) DynamicClient() dynamic.Interface {
	return h.dynamicClient
}

// DiscoveryClient returns underlying discovery client.
func (h *Handler) DiscoveryClient() *discovery.DiscoveryClient {
	return h.discoveryClient
}

// GVR returns the name of Group, Version, Resource of pod resource.
func GVR() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    corev1.SchemeGroupVersion.Group,
		Version:  corev1.SchemeGroupVersion.Version,
		Resource: "pods",
	}
}

// Group returns Group name of pod resource.
func Group() string {
	return corev1.SchemeGroupVersion.Group
}

// Version returns Version name of pod resource.
func Version() string {
	return corev1.SchemeGroupVersion.Version
}
