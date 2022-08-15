package daemonset

import (
	"context"
	"net/http"
	"os"
	"sync"

	"github.com/forbearing/k8s/types"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
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

	Options *types.HandlerOptions

	l sync.RWMutex
}

// NewOrDie simply call New() to get a daemonset handler.
// panic if there is any error occurs.
func NewOrDie(ctx context.Context, kubeconfig, namespace string) *Handler {
	handler, err := New(ctx, kubeconfig, namespace)
	if err != nil {
		panic(err)
	}
	return handler
}

// New returns a handler handler from kubeconfig or in-cluster config.
// The kubeconfig precedence is:
// * kubeconfig variable passed.
// * KUBECONFIG environment variable pointing at a file
// * $HOME/.kube/config if exists.
// * In-cluster config if running in cluster
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

	// create rest config, and config precedence.
	// * kubeconfig variable passed.
	// * KUBECONFIG environment variable pointing at a file
	// * $HOME/.kube/config if exists.
	// * In-cluster config if running in cluster
	//
	// create the outside-cluster config
	if len(kubeconfig) != 0 {
		if config, err = clientcmd.BuildConfigFromFlags("", kubeconfig); err != nil {
			return nil, err
		}
	} else if len(os.Getenv(clientcmd.RecommendedConfigPathEnvVar)) != 0 {
		if config, err = clientcmd.BuildConfigFromFlags("", os.Getenv(clientcmd.RecommendedConfigPathEnvVar)); err != nil {
			return nil, err
		}
	} else if len(clientcmd.RecommendedHomeFile) != 0 {
		if config, err = clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile); err != nil {
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
	config.NegotiatedSerializer = scheme.Codecs

	// create a http client for the given config.
	httpClient, err = rest.HTTPClientFor(config)
	if err != nil {
		return nil, err
	}
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
	handler.Options = &types.HandlerOptions{}

	return handler, nil
}

// WithNamespace deep copies a new handler, but set the handler.namespace to
// the provided namespace.
func (h *Handler) WithNamespace(namespace string) *Handler {
	handler := h.DeepCopy()
	handler.resetNamespace(namespace)
	return handler
}

// WithDryRun deep copies a new handler and prints the create/update/apply/delete
// operations, without sending it to apiserver.
func (h *Handler) WithDryRun() *Handler {
	handler := h.DeepCopy()
	handler.Options.CreateOptions.DryRun = []string{metav1.DryRunAll}
	handler.Options.UpdateOptions.DryRun = []string{metav1.DryRunAll}
	handler.Options.DeleteOptions.DryRun = []string{metav1.DryRunAll}
	handler.Options.PatchOptions.DryRun = []string{metav1.DryRunAll}
	handler.Options.ApplyOptions.DryRun = []string{metav1.DryRunAll}
	return handler
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

	out.Options = &types.HandlerOptions{}
	out.Options.ListOptions = *in.Options.ListOptions.DeepCopy()
	out.Options.GetOptions = *in.Options.GetOptions.DeepCopy()
	out.Options.CreateOptions = *in.Options.CreateOptions.DeepCopy()
	out.Options.UpdateOptions = *in.Options.UpdateOptions.DeepCopy()
	out.Options.PatchOptions = *in.Options.PatchOptions.DeepCopy()
	out.Options.ApplyOptions = *in.Options.ApplyOptions.DeepCopy()

	return out
}
func (h *Handler) resetNamespace(namespace string) {
	h.l.Lock()
	defer h.l.Unlock()
	h.namespace = namespace
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
		h.Options.DeleteOptions.GracePeriodSeconds = new(int64)
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

// GVK returns the name of Group, Version, Kind of daemonset resource.
func GVK() schema.GroupVersionKind {
	return schema.GroupVersionKind{
		Group:   appsv1.SchemeGroupVersion.Group,
		Version: appsv1.SchemeGroupVersion.Version,
		Kind:    types.KindDaemonSet,
	}
}

// GVR returns the name of Group, Version, Resource of daemonset resource.
func GVR() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    appsv1.SchemeGroupVersion.Group,
		Version:  appsv1.SchemeGroupVersion.Version,
		Resource: types.ResourceDaemonSet,
	}
}

// Kind returns the Kind name of daemonset resource.
func Kind() string {
	return GVK().Kind
}

// Group returns the Group name of daemonset resource.
func Group() string {
	return GVR().Group
}

// Version returns the Version name of daemonset resource.
func Version() string {
	return GVR().Version
}

// Resource returns the Resource name of daemonset resource.
func Resource() string {
	return GVR().Resource
}
