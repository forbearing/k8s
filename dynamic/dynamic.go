package dynamic

import (
	"context"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/forbearing/k8s/types"
	"github.com/forbearing/k8s/util/client"
	utilrestmapper "github.com/forbearing/k8s/util/restmapper"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type Handler struct {
	ctx        context.Context
	gvk        schema.GroupVersionKind
	kubeconfig string
	namespace  string

	config        *rest.Config
	httpClient    *http.Client
	restClient    *rest.RESTClient
	dynamicClient dynamic.Interface
	restMapper    meta.RESTMapper

	resyncPeriod     time.Duration
	informerScope    string
	tweakListOptions dynamicinformer.TweakListOptionsFunc
	informerFactory  dynamicinformer.DynamicSharedInformerFactory

	Options *types.HandlerOptions

	l sync.RWMutex
}

// NewOrDie creates a dynamic client.
// Panic if there is any error.
func NewOrDie(ctx context.Context, kubeconfig string, namespace string) *Handler {
	handler, err := New(ctx, kubeconfig, namespace)
	if err != nil {
		panic(err)
	}
	return handler
}

// New creates a dynamic client from kubeconfig or in-cluster config.
// If provided namespace is empty, it means the k8s resources created/updated/deleted
// by dynamic client is cluster scope. or it's namespaced scope.
//
// The kubeconfig precedence is:
// * kubeconfig variable passed.
// * KUBECONFIG environment variable pointing at a file.
// * $HOME/.kube/config if exists.
// * In-cluster config if running in cluster.
func New(ctx context.Context, kubeconfig string, namespace string) (*Handler, error) {
	var (
		err           error
		config        *rest.Config
		httpClient    *http.Client
		restClient    *rest.RESTClient
		dynamicClient dynamic.Interface
		restMapper    meta.RESTMapper
	)

	// create rest config, and config precedence.
	// * kubeconfig variable passed.
	// * KUBECONFIG environment variable pointing at a file.
	// * $HOME/.kube/config if exists.
	// * In-cluster config if running in cluster.
	if config, err = client.RESTConfig(kubeconfig); err != nil {
		return nil, err
	}

	config.APIPath = "api"
	config.GroupVersion = &schema.GroupVersion{}
	config.NegotiatedSerializer = scheme.Codecs

	// create a http client for the given config.
	if httpClient, err = rest.HTTPClientFor(config); err != nil {
		return nil, err
	}
	// create a RESTClient for the given config and http client.
	if restClient, err = rest.RESTClientForConfigAndClient(config, httpClient); err != nil {
		return nil, err
	}
	// create a dynamic client for the given config and http client.
	if dynamicClient, err = dynamic.NewForConfigAndClient(config, httpClient); err != nil {
		return nil, err
	}
	if restMapper, err = utilrestmapper.NewRESTMapper(kubeconfig); err != nil {
		return nil, err
	}
	if len(namespace) == 0 {
		namespace = metav1.NamespaceDefault
	}

	return &Handler{
		ctx:           ctx,
		kubeconfig:    kubeconfig,
		namespace:     namespace,
		config:        config,
		httpClient:    httpClient,
		restClient:    restClient,
		dynamicClient: dynamicClient,
		restMapper:    restMapper,
		Options:       &types.HandlerOptions{},
	}, nil
}

// WithNamespace returns the same handler but with provided namespace.
// If the k8s resource is namespace scope, it will create/delete/update/apply
// k8s resource in the new namespace.
// If the k8s resource is cluster scope, it will ignore the namespace.
//
// But the namespace defined in yaml file have higher precedence than namespace specified here.
//
// If no namespace is defined in yaml file and no namespace is specified using
// WithNamespace() method, then the namespace default to metav1.NamespaceDefault("default").
//
// namespace precedence:
// * namespace defined in yaml file or json file.
// * namespace specified by WithNamespace() method.
// * namespace specified in dynamic.New() or dynamic.NewOrDie() funciton.
// * namespace will be ignored if k8s resource is cluster scope.
func (h *Handler) WithNamespace(namespace string) *Handler {
	handler := h.DeepCopy()
	if len(namespace) == 0 {
		namespace = metav1.NamespaceDefault
	}
	handler.namespace = namespace
	return handler
}

// WithGVK returns the same handler but with provided group, version and resource.
func (h *Handler) WithGVK(gvk schema.GroupVersionKind) *Handler {
	handler := h.DeepCopy()
	handler.gvk = gvk
	return handler
}

// WithDryRun deep copies a new handler and prints the create/update/apply/delete
// operations, without sending it to apiserver.
func (h *Handler) WithDryRun() *Handler {
	handler := h.DeepCopy()
	handler.Options.CreateOptions.DryRun = []string{metav1.DryRunAll}
	handler.Options.UpdateOptions.DryRun = []string{metav1.DryRunAll}
	handler.Options.ApplyOptions.DryRun = []string{metav1.DryRunAll}
	handler.Options.DeleteOptions.DryRun = []string{metav1.DryRunAll}
	handler.Options.PatchOptions.DryRun = []string{metav1.DryRunAll}
	return handler
}

// DeepCopy
func (in *Handler) DeepCopy() *Handler {
	if in == nil {
		return nil
	}
	return &Handler{
		ctx:              in.ctx,
		gvk:              in.gvk,
		kubeconfig:       in.kubeconfig,
		namespace:        in.namespace,
		config:           in.config,
		httpClient:       in.httpClient,
		restClient:       in.restClient,
		dynamicClient:    in.dynamicClient,
		informerFactory:  in.informerFactory,
		resyncPeriod:     in.resyncPeriod,
		informerScope:    in.informerScope,
		tweakListOptions: in.tweakListOptions,
		restMapper:       in.restMapper,
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

// SetTimeout
func (h *Handler) SetTimeout(timeout int64) {
	h.l.Lock()
	defer h.l.Unlock()
	h.Options.ListOptions.TimeoutSeconds = &timeout
}

// SetLimit
func (h *Handler) SetLimit(limit int64) {
	h.l.Lock()
	defer h.l.Unlock()
	h.Options.ListOptions.Limit = limit
}

// SetForceDelete
func (h *Handler) SetForceDelete(force bool) {
	h.l.Lock()
	defer h.l.Unlock()
	if force {
		h.Options.DeleteOptions.GracePeriodSeconds = new(int64)
	}
}

// SetPropagationPolicy will set the PropagationPolicy.
// If we delete job or/and cronjob, we should always set the PropagationPolicy to
// DeletePropagationBackground to delete all pods managed by that job or/and cronjob.
// Default to "DeletePropagationBackground" to job and/or cronjob.
func (h *Handler) SetPropagationPolicy(policy string) {
	h.l.Lock()
	defer h.l.Unlock()
	switch strings.ToLower(policy) {
	case strings.ToLower(string(metav1.DeletePropagationBackground)):
		propagationPolicy := metav1.DeletePropagationBackground
		h.Options.DeleteOptions.PropagationPolicy = &propagationPolicy
	case strings.ToLower(string(metav1.DeletePropagationForeground)):
		propagationPolicy := metav1.DeletePropagationForeground
		h.Options.DeleteOptions.PropagationPolicy = &propagationPolicy
	case strings.ToLower(string(metav1.DeletePropagationOrphan)):
		propagationPolicy := metav1.DeletePropagationOrphan
		h.Options.DeleteOptions.PropagationPolicy = &propagationPolicy
	default:
		propagationPolicy := metav1.DeletePropagationBackground
		h.Options.DeleteOptions.PropagationPolicy = &propagationPolicy
	}
}

// DynamicClient returns the underlying dynamic client used by this dynamic handler.
func (h *Handler) DynamicClient() dynamic.Interface {
	return h.dynamicClient
}
