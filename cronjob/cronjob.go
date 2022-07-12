package cronjob

import (
	"context"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/forbearing/k8s/typed"
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
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
	informer        cache.SharedIndexInformer

	Options *typed.HandlerOptions

	l sync.Mutex
}

// New returns a cronjob handler from kubeconfig or in-cluster config.
func New(ctx context.Context, namespace, kubeconfig string) (handler *Handler, err error) {
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
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return
		}
	} else {
		// create the in-cluster config
		config, err = rest.InClusterConfig()
		if err != nil {
			return
		}
	}

	// setup APIPath, GroupVersion and NegotiatedSerializer before initializing a RESTClient
	config.APIPath = "api"
	config.GroupVersion = &batchv1.SchemeGroupVersion
	config.NegotiatedSerializer = scheme.Codecs

	// create a http client for the given config.
	httpClient, err = rest.HTTPClientFor(config)
	if err != nil {
		return nil, err
	}
	// create a RESTClient for the given config and http client.
	restClient, err = rest.RESTClientForConfigAndClient(config, httpClient)
	if err != nil {
		return
	}
	// create a Clientset for the given config and http client.
	clientset, err = kubernetes.NewForConfigAndClient(config, httpClient)
	if err != nil {
		return
	}
	// create a dynamic client for the given config and http client.
	dynamicClient, err = dynamic.NewForConfigAndClient(config, httpClient)
	if err != nil {
		return
	}
	// create a DiscoveryClient for the given config and http client.
	discoveryClient, err = discovery.NewDiscoveryClientForConfigAndClient(config, httpClient)
	if err != nil {
		return
	}
	// create a sharedInformerFactory for all namespaces.
	informerFactory = informers.NewSharedInformerFactory(clientset, time.Minute)

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
	handler.informer = informerFactory.Batch().V1().CronJobs().Informer()
	handler.Options = &typed.HandlerOptions{}

	return
}
func (h *Handler) Namespace() string {
	return h.namespace
}
func (in *Handler) DeepCopy() *Handler {
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
	out.informer = in.informer

	out.Options = &typed.HandlerOptions{}
	out.Options.ListOptions = *in.Options.ListOptions.DeepCopy()
	out.Options.GetOptions = *in.Options.GetOptions.DeepCopy()
	out.Options.CreateOptions = *in.Options.CreateOptions.DeepCopy()
	out.Options.UpdateOptions = *in.Options.UpdateOptions.DeepCopy()
	out.Options.PatchOptions = *in.Options.PatchOptions.DeepCopy()
	out.Options.ApplyOptions = *in.Options.ApplyOptions.DeepCopy()
	out.SetPropagationPolicy("background")

	return out
}
func (h *Handler) resetNamespace(namespace string) {
	h.l.Lock()
	defer h.l.Unlock()
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
	handler.SetPropagationPolicy("background")
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
	defer h.l.Lock()
	if force {
		gracePeriodSeconds := int64(0)
		h.Options.DeleteOptions.GracePeriodSeconds = &gracePeriodSeconds
		propagationPolicy := metav1.DeletePropagationBackground
		h.Options.DeleteOptions.PropagationPolicy = &propagationPolicy
	} else {
		h.Options.DeleteOptions = metav1.DeleteOptions{}
		propagationPolicy := metav1.DeletePropagationBackground
		h.Options.DeleteOptions.PropagationPolicy = &propagationPolicy
	}
}

// Whether and how garbage collection will be performed.
// support value are "Background", "Orphan", "Foreground",
// default value is "Background"
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

func (h *Handler) RESTClient() *rest.RESTClient {
	return h.restClient
}
func (h *Handler) Clientset() *kubernetes.Clientset {
	return h.clientset
}
func (h *Handler) DynamicClient() dynamic.Interface {
	return h.dynamicClient
}
func (h *Handler) DiscoveryClient() *discovery.DiscoveryClient {
	return h.discoveryClient
}
