package dynamic

import (
	"context"
	"net/http"
	"strings"
	"sync"

	"github.com/forbearing/k8s/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Handler struct {
	kubeconfig         string
	namespace          string
	gvr                schema.GroupVersionResource
	namespacedResource bool

	ctx           context.Context
	config        *rest.Config
	httpClient    *http.Client
	restClient    *rest.RESTClient
	dynamicClient dynamic.Interface

	Options *types.HandlerOptions

	l sync.RWMutex
}

// NewOrDie creates a dynamic client.
// Panic if there is any error.
func NewOrDie(ctx context.Context, kubeconfig, namespace string, gvr schema.GroupVersionResource) *Handler {
	handler, err := New(ctx, kubeconfig, namespace, gvr)
	if err != nil {
		panic(err)
	}
	return handler
}

// New creates a dynamic client.
// If provided namespace is empty, it means the k8s resources created/updated/deleted
// by dynamic client is cluster scope. or it's namespaced scope.
// The dynamic client is reuseable, WithNamespace(), WithGVR()
func New(ctx context.Context, kubeconfig, namespace string, gvr schema.GroupVersionResource) (*Handler, error) {
	var (
		config        *rest.Config
		httpClient    *http.Client
		restClient    *rest.RESTClient
		dynamicClient dynamic.Interface
	)
	handler := &Handler{}
	var err error

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

	config.APIPath = "api"
	config.GroupVersion = &schema.GroupVersion{Group: gvr.Group, Version: gvr.Version}
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

	//if len(namespace) == 0 {
	//    namespace = metav1.NamespaceDefault
	//}
	if len(gvr.Version) == 0 {
		return nil, ErrVersionEmpty
	}
	if len(gvr.Resource) == 0 {
		return nil, ErrResourceEmpty
	}

	handler.kubeconfig = kubeconfig
	handler.namespacedResource = true
	handler.namespace = namespace
	handler.gvr = gvr
	handler.ctx = ctx
	handler.config = config
	handler.httpClient = httpClient
	handler.restClient = restClient
	handler.dynamicClient = dynamicClient
	handler.Options = &types.HandlerOptions{}

	return handler, nil
}

func (in *Handler) DeepCopy() *Handler {
	if in == nil {
		return nil
	}
	out := new(Handler)

	out.kubeconfig = in.kubeconfig
	out.namespacedResource = in.namespacedResource
	out.namespace = in.namespace
	out.gvr = in.gvr
	out.ctx = in.ctx
	out.config = in.config
	out.httpClient = in.httpClient
	out.restClient = in.restClient
	out.dynamicClient = in.dynamicClient

	out.Options = &types.HandlerOptions{}
	out.Options.CreateOptions = *in.Options.CreateOptions.DeepCopy()
	out.Options.UpdateOptions = *in.Options.UpdateOptions.DeepCopy()
	out.Options.ApplyOptions = *in.Options.ApplyOptions.DeepCopy()
	out.Options.DeleteOptions = *in.Options.DeleteOptions.DeepCopy()
	out.Options.GetOptions = *in.Options.GetOptions.DeepCopy()
	out.Options.ListOptions = *in.Options.ListOptions.DeepCopy()
	out.Options.PatchOptions = *in.Options.PatchOptions.DeepCopy()

	if in.gvr.Resource == "jobs" || in.gvr.Resource == "cronjobs" {
		out.setPropagationPolicy("background")
	}

	return out
}
func (h *Handler) setPropagationPolicy(policy string) {
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
func (h *Handler) IsNamespacedResource() bool {
	if len(h.namespace) == 0 {
		return false
	}
	return true
}

//// SetNamespacedResource sets whether this unstructured resource is namespaced
//// k8s resource, default is namespaced k8s resource.
//// If this unstructured k8s resource is global scope but SetNamespacedResource(false)
//// not called, then create/delete/update/delete/get will be failed.

//// For example, if you create "deployment" resource with dynamic client handler,
//// just call Create() method, but if you create "namespace" resource with dynamic
//// client handler without call SetNamespacedResource(false) before to mark it is
//// global k8s resource, it will fail.

//func (h *Handler) SetNamespacedResource(namespaced bool) {
//    h.l.Lock()
//    defer h.l.Unlock()
//    h.namespacedResource = namespaced
//}

// WithNamespace returns the same handler but with provided namespace.
func (h *Handler) WithNamespace(namespace string) *Handler {
	handler := h.DeepCopy()
	handler.namespace = namespace
	return handler
}

// WithGVR returns the same handler but with provided group, version and resource.
func (h *Handler) WithGVR(gvr schema.GroupVersionResource) *Handler {
	handler := h.DeepCopy()
	handler.gvr = gvr
	return handler
}

// DynamicClient returns the underlying dynamic client.
func (h *Handler) DynamicClient() dynamic.Interface {
	return h.dynamicClient
}
