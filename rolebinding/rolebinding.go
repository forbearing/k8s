package rolebinding

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/forbearing/k8s/types"
	"github.com/forbearing/k8s/util/client"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/informers/internalinterfaces"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type Handler struct {
	ctx        context.Context
	kubeconfig string
	namespace  string

	config          *rest.Config
	httpClient      *http.Client
	restClient      *rest.RESTClient
	clientset       *kubernetes.Clientset
	dynamicClient   dynamic.Interface
	discoveryClient *discovery.DiscoveryClient

	resyncPeriod     time.Duration
	informerScope    string
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	informerFactory  informers.SharedInformerFactory

	Options *types.HandlerOptions

	l sync.RWMutex
}

// NewOrDie simply call New() to get a rolebinding handler.
// panic if there is any error occurs.
func NewOrDie(ctx context.Context, kubeconfig, namespace string) *Handler {
	handler, err := New(ctx, kubeconfig, namespace)
	if err != nil {
		panic(err)
	}
	return handler
}

// New returns a Handler handler from kubeconfig or in-cluster config.
// The kubeconfig precedence is:
// * kubeconfig variable passed.
// * KUBECONFIG environment variable pointing at a file.
// * $HOME/.kube/config if exists.
// * In-cluster config if running in cluster.
func New(ctx context.Context, kubeconfig, namespace string) (*Handler, error) {
	var (
		err             error
		config          *rest.Config
		httpClient      *http.Client
		restClient      *rest.RESTClient
		clientset       *kubernetes.Clientset
		dynamicClient   dynamic.Interface
		discoveryClient *discovery.DiscoveryClient
		informerFactory informers.SharedInformerFactory
	)

	// create rest config, and config precedence.
	// * kubeconfig variable passed.
	// * KUBECONFIG environment variable pointing at a file.
	// * $HOME/.kube/config if exists.
	// * In-cluster config if running in cluster.
	if config, err = client.RESTConfig(kubeconfig); err != nil {
		return nil, err
	}
	// setup APIPath, GroupVersion and NegotiatedSerializer before initializing a RESTClient
	config.APIPath = "api"
	config.GroupVersion = &rbacv1.SchemeGroupVersion
	config.NegotiatedSerializer = scheme.Codecs

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
	if len(namespace) == 0 {
		namespace = metav1.NamespaceDefault
	}
	// create a sharedInformerFactory for all namespaces.
	informerFactory = informers.NewSharedInformerFactory(clientset, 0)

	return &Handler{
		ctx:             ctx,
		kubeconfig:      kubeconfig,
		namespace:       namespace,
		config:          config,
		httpClient:      httpClient,
		restClient:      restClient,
		clientset:       clientset,
		dynamicClient:   dynamicClient,
		discoveryClient: discoveryClient,
		informerFactory: informerFactory,
		Options:         &types.HandlerOptions{},
	}, nil
}

// WithNamespace deep copies a new handler, but set the handler.namespace to
// the provided namespace.
func (h *Handler) WithNamespace(namespace string) *Handler {
	handler := h.DeepCopy()
	handler.ResetNamespace(namespace)
	return handler
}

// WithDryRun deep copies a new handler and prints the create/update/apply/delete
// operations, without sending it to apiserver.
func (h *Handler) WithDryRun() *Handler {
	handler := h.DeepCopy()
	handler.Options.CreateOptions.DryRun = []string{metav1.DryRunAll}
	handler.Options.UpdateOptions.DryRun = []string{metav1.DryRunAll}
	handler.Options.ApplyOptions.DryRun = []string{metav1.DryRunAll}
	handler.Options.PatchOptions.DryRun = []string{metav1.DryRunAll}
	handler.Options.DeleteOptions.DryRun = []string{metav1.DryRunAll}
	return handler
}
func (in *Handler) DeepCopy() *Handler {
	if in == nil {
		return nil
	}
	return &Handler{
		ctx:              in.ctx,
		kubeconfig:       in.kubeconfig,
		namespace:        in.namespace,
		config:           in.config,
		httpClient:       in.httpClient,
		restClient:       in.restClient,
		clientset:        in.clientset,
		dynamicClient:    in.dynamicClient,
		discoveryClient:  in.discoveryClient,
		informerFactory:  in.informerFactory,
		resyncPeriod:     in.resyncPeriod,
		informerScope:    in.informerScope,
		tweakListOptions: in.tweakListOptions,
		Options: &types.HandlerOptions{
			CreateOptions: *in.Options.CreateOptions.DeepCopy(),
			UpdateOptions: *in.Options.UpdateOptions.DeepCopy(),
			ApplyOptions:  *in.Options.ApplyOptions.DeepCopy(),
			DeleteOptions: *in.Options.DeleteOptions.DeepCopy(),
			GetOptions:    *in.Options.GetOptions.DeepCopy(),
			ListOptions:   *in.Options.ListOptions.DeepCopy(),
			PatchOptions:  *in.Options.PatchOptions.DeepCopy(),
		},
	}
}
func (h *Handler) ResetNamespace(namespace string) {
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

// GVK contains the Group, Version, Kind name of rolebinding.
var GVK = schema.GroupVersionKind{
	Group:   rbacv1.SchemeGroupVersion.Group,
	Version: rbacv1.SchemeGroupVersion.Version,
	Kind:    types.KindRoleBinding,
}

// GVR contains the Group, Version and Resource name of rolebinding.
var GVR = schema.GroupVersionResource{
	Group:    rbacv1.SchemeGroupVersion.Group,
	Version:  rbacv1.SchemeGroupVersion.Version,
	Resource: types.ResourceRoleBinding,
}

// Kind is the rolebinding Kind name.
var Kind = GVK.Kind

// Group is the rolebinding Group name.
var Group = GVK.Group

// Version is the rolebinding Version name.
var Version = GVK.Version

// Resource is the rolebinding Resource name.
var Resource = GVR.Resource
