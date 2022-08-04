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
	group              string
	version            string
	resource           string
	namespacedResource bool

	ctx           context.Context
	config        *rest.Config
	httpClient    *http.Client
	restClient    *rest.RESTClient
	dynamicClient dynamic.Interface

	Options *types.HandlerOptions

	l sync.RWMutex
}

func NewOrDie(ctx context.Context, kubeconfig, namespace, group, version, resource string) *Handler {
	handler, err := New(ctx, kubeconfig, namespace, group, version, resource)
	if err != nil {
		panic(err)
	}
	return handler
}

func New(ctx context.Context, kubeconfig, namespace, group, version, resource string) (*Handler, error) {
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
	config.GroupVersion = &schema.GroupVersion{Group: group, Version: version}
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
	if len(version) == 0 {
		return nil, ErrVersionEmpty
	}
	if len(resource) == 0 {
		return nil, ErrResourceEmpty
	}

	handler.kubeconfig = kubeconfig
	handler.namespacedResource = true
	handler.namespace = namespace
	handler.group = group
	handler.version = version
	handler.resource = resource
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
	out.group = in.group
	out.version = in.version
	out.resource = in.resource
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

	if in.resource == "jobs" || in.resource == "cronjobs" {
		in.setPropagationPolicy("background")
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

// Namespace returns the same handler but with provided namespace.
func (h *Handler) Namespace(namespace string) *Handler {
	handler := h.DeepCopy()
	handler.namespace = namespace
	return handler
}
func (h *Handler) N(namespace string) *Handler {
	return h.Namespace(namespace)
}

// Group returns the same handler but with provided group.
func (h *Handler) Group(group string) *Handler {
	handler := h.DeepCopy()
	handler.group = group
	return handler
}
func (h *Handler) G(group string) *Handler {
	return h.Group(group)
}

// Version returns the same handler but with provided version.
func (h *Handler) Version(version string) *Handler {
	handler := h.DeepCopy()
	handler.version = version
	return handler
}
func (h *Handler) V(version string) *Handler {
	return h.Version(version)
}

// Resource returns the same handler but with provided resource.
func (h *Handler) Resource(resource string) *Handler {
	handler := h.DeepCopy()
	handler.resource = resource
	return handler
}
func (h *Handler) R(resource string) *Handler {
	return h.Resource(resource)
}

// GVR returns the same handler but with provided group, version and resource.
func (h *Handler) GVR(gvr schema.GroupVersionResource) *Handler {
	handler := h.DeepCopy()
	handler.group = gvr.Group
	handler.version = gvr.Version
	handler.resource = gvr.Resource
	return handler
}

// DynamicClient returns the underlying dynamic client.
func (h *Handler) DynamicClient() dynamic.Interface {
	return h.dynamicClient
}

// gvr
func (h *Handler) gvr() schema.GroupVersionResource {
	h.l.RLock()
	defer h.l.RUnlock()
	return schema.GroupVersionResource{
		Group:    h.group,
		Version:  h.version,
		Resource: h.resource,
	}
}
