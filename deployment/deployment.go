package deployment

import (
	"context"
	"net/http"
	"sync"
	"time"

	//_ "k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	//_ "k8s.io/client-go/applyconfigurations/apps/v1"
	//_ "k8s.io/client-go/applyconfigurations/meta/v1"

	"github.com/forbearing/k8s/types"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	listersappsv1 "k8s.io/client-go/listers/apps/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

type Handler struct {
	kubeconfig string
	namespace  string

	ctx                context.Context
	config             *rest.Config
	httpClient         *http.Client
	restClient         *rest.RESTClient
	clientset          *kubernetes.Clientset
	dynamicClient      dynamic.Interface
	discoveryClient    *discovery.DiscoveryClient
	discoveryInterface discovery.DiscoveryInterface
	informerFactory    informers.SharedInformerFactory
	informer           cache.SharedIndexInformer
	lister             listersappsv1.DeploymentLister

	Options *types.HandlerOptions

	l sync.RWMutex
}

//// Discovery retrieves the DiscoveryClient
//func (c *Clientset) Discovery() discovery.DiscoveryInterface {
//    if c == nil {
//        return nil
//    }
//    return c.DiscoveryClient
//}
// clientset 调用 Discovery 方法可以得到一个 discovery.DiscoveryInterface
// discovery.DiscoveryClient 其实就是 discovery.DiscoveryInterface 的一个实现

// NewOrDie simply call New() to get a deployment handler.
// panic if there is any error occurs.
func NewOrDie(ctx context.Context, kubeconfig, namespace string) *Handler {
	handler, err := New(ctx, kubeconfig, namespace)
	if err != nil {
		panic(err)
	}
	return handler
}

// New returns a deployment handler from kubeconfig or in-cluster config.
func New(ctx context.Context, kubeconfig, namespace string) (handler *Handler, err error) {
	var (
		config             *rest.Config
		httpClient         *http.Client
		restClient         *rest.RESTClient
		clientset          *kubernetes.Clientset
		dynamicClient      dynamic.Interface
		discoveryClient    *discovery.DiscoveryClient
		discoveryInterface discovery.DiscoveryInterface
		informerFactory    informers.SharedInformerFactory
	)
	handler = &Handler{}

	// create rest config
	if len(kubeconfig) != 0 {
		// use the current context in kubeconfig
		if config, err = clientcmd.BuildConfigFromFlags("", kubeconfig); err != nil {
			return nil, err
		}
	} else {
		// create the in-cluster config
		if config, err = rest.InClusterConfig(); err != nil {
			return nil, err
		}
	}

	// setup APIPath, GroupVersion and NegotiatedSerializer before initializing a RESTClient
	config.APIPath = "api"
	config.GroupVersion = &appsv1.SchemeGroupVersion
	//config.GroupVersion = &schema.GroupVersion{Group: "apps", Version: "v1"}
	config.NegotiatedSerializer = scheme.Codecs
	//config.UserAgent = rest.DefaultKubernetesUserAgent()

	//// k8s cluster endpoint, eg: https://10.250.16.10:8443
	//config.Host = "127.0.0.1"
	//config.ContentConfig = rest.ContentConfig{
	//    GroupVersion:         &corev1.SchemeGroupVersion,
	//    NegotiatedSerializer: scheme.Codecs,
	//}

	// create a http client for the given config
	if httpClient, err = rest.HTTPClientFor(config); err != nil {
		return nil, err
	}

	//// create a RESTClient for the given config
	//restClient, err = rest.RESTClientFor(config)
	//if err != nil {
	//    return nil, err
	//}
	// create a RESTClient for the given config and http client
	if restClient, err = rest.RESTClientForConfigAndClient(config, httpClient); err != nil {
		return nil, err
	}

	//// create a Clientset for the given config
	//clientset, err = kubernetes.NewForConfig(config)
	//if err != nil {
	//    return nil, err
	//}
	// create a clientset for the given config and http client.
	if clientset, err = kubernetes.NewForConfigAndClient(config, httpClient); err != nil {
		return nil, err
	}
	//// create a dynamic client for the given config
	//dynamicClient, err = dynamic.NewForConfig(config)
	//if err != nil {
	//    return nil, err
	//}
	// create a dynamic client for the given config and http client.
	if dynamicClient, err = dynamic.NewForConfigAndClient(config, httpClient); err != nil {
		return nil, err
	}

	//// create a DiscoveryClient for the given config
	//discoveryClient, err = discovery.NewDiscoveryClientForConfig(config)
	//if err != nil {
	//    return nil, err
	//}
	// create a DiscoveryClient for the given config and http client.
	if discoveryClient, err = discovery.NewDiscoveryClientForConfigAndClient(config, httpClient); err != nil {
		return nil, err
	}

	// create a sharedInformerFactory for all namespaces.
	informerFactory = informers.NewSharedInformerFactory(clientset, time.Second*30)
	//discoveryClient = clientset.DiscoveryClient
	//discoveryInterface = clientset.Discovery()

	// default namespace is meatv1.NamespaceDefault ("default")
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
	handler.informer = informerFactory.Apps().V1().Deployments().Informer()
	handler.lister = informerFactory.Apps().V1().Deployments().Lister()
	//handler.discoveryInterface = discoveryInterface
	_ = discoveryInterface

	handler.Options = &types.HandlerOptions{}

	return handler, nil
}
func (h *Handler) Namespace() string {
	return h.namespace
}
func (in *Handler) DeepCopy() *Handler {
	if in == nil {
		return nil
	}
	out := new(Handler)

	out.kubeconfig = in.kubeconfig
	out.namespace = in.namespace

	// 和几个字段都是共用的, 不需要深拷贝
	out.ctx = in.ctx
	out.config = in.config
	out.httpClient = in.httpClient
	out.restClient = in.restClient
	out.clientset = in.clientset
	out.dynamicClient = in.dynamicClient
	out.discoveryClient = in.discoveryClient
	out.informerFactory = in.informerFactory
	out.informer = in.informer
	out.lister = in.lister

	out.Options = &types.HandlerOptions{}
	out.Options.ListOptions = *in.Options.ListOptions.DeepCopy()
	out.Options.GetOptions = *in.Options.GetOptions.DeepCopy()
	out.Options.CreateOptions = *in.Options.CreateOptions.DeepCopy()
	out.Options.UpdateOptions = *in.Options.UpdateOptions.DeepCopy()
	out.Options.PatchOptions = *in.Options.PatchOptions.DeepCopy()
	out.Options.ApplyOptions = *in.Options.ApplyOptions.DeepCopy()

	// 锁 sync.RWMutex 不需要拷贝, 也不能拷贝. 拷贝 sync.RWMutex 会直接 panic

	return out
}
func (h *Handler) resetNamespace(namespace string) {
	h.l.Lock()
	h.l.Unlock()
	h.namespace = namespace
}
func (h *Handler) WithNamespace(namespace string) *Handler {
	handler := h.DeepCopy()
	handler.resetNamespace(namespace)
	return handler
}
func (h *Handler) WithDryRun() *Handler {
	handler := h.DeepCopy()
	handler.Options.CreateOptions.DryRun = []string{metav1.DryRunAll}
	handler.Options.UpdateOptions.DryRun = []string{metav1.DryRunAll}
	handler.Options.DeleteOptions.DryRun = []string{metav1.DryRunAll}
	handler.Options.PatchOptions.DryRun = []string{metav1.DryRunAll}
	handler.Options.ApplyOptions.DryRun = []string{metav1.DryRunAll}
	return handler
}
func (h *Handler) SetTimeout(timeout int64) {
	h.l.Lock()
	defer h.l.Unlock()
	h.Options.ListOptions.TimeoutSeconds = &timeout
}
func (h *Handler) SetLimit(limit int64) {
	h.l.Lock()
	defer h.l.Unlock()
	h.Options.ListOptions.Limit = limit
}
func (h *Handler) SetForceDelete(force bool) {
	h.l.Lock()
	defer h.l.Unlock()
	if force {
		gracePeriodSeconds := int64(0)
		h.Options.DeleteOptions.GracePeriodSeconds = &gracePeriodSeconds
	} else {
		h.Options.DeleteOptions = metav1.DeleteOptions{}
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

// GVR returns the name of Group, Version, Resource of deployment resource.
func (h *Handler) GVR() *schema.GroupVersionResource {
	return &schema.GroupVersionResource{
		Group:    h.config.GroupVersion.Group,
		Version:  h.config.GroupVersion.Version,
		Resource: "deployments",
	}
}
