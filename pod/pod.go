package pod

// TODO:
// 1. 加了 Options 参数之后可能要修改的东西: List, Watch, WaitReady
// 2. 精简代码, 比如返回值
import (
	"context"
	"sync"
	"time"

	//_ "k8s.io/client-go/kubernetes/typed/core/v1"

	"github.com/forbearing/k8s/typed"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

type Handler struct {
	kubeconfig string
	namespace  string

	ctx             context.Context
	config          *rest.Config
	restClient      *rest.RESTClient
	clientset       *kubernetes.Clientset
	dynamicClient   dynamic.Interface
	discoveryClient *discovery.DiscoveryClient
	informerFactory informers.SharedInformerFactory
	informer        cache.SharedIndexInformer
	client          typedcorev1.PodInterface

	Options *typed.HandlerOptions

	sync.Mutex
}

// New returns a pod handler from kubeconfig or in-cluster config.
func New(ctx context.Context, namespace, kubeconfig string) (handler *Handler, err error) {
	var (
		config          *rest.Config
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
		// creates the in-cluster config
		config, err = rest.InClusterConfig()
		if err != nil {
			return
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

	// create a RESTClient for the given config
	restClient, err = rest.RESTClientFor(config)
	if err != nil {
		return
	}
	// create a Clientset for the given config
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		return
	}
	// create a dynamic client for the given config
	dynamicClient, err = dynamic.NewForConfig(config)
	if err != nil {
		return
	}
	// create a DiscoveryClient for the given config
	discoveryClient, err = discovery.NewDiscoveryClientForConfig(config)
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
	handler.restClient = restClient
	handler.clientset = clientset
	handler.dynamicClient = dynamicClient
	handler.discoveryClient = discoveryClient
	handler.informerFactory = informerFactory
	handler.informer = informerFactory.Core().V1().Pods().Informer()
	handler.client = clientset.CoreV1().Pods(namespace)
	handler.Options = &typed.HandlerOptions{}

	return
}
func (p *Handler) Namespace() string {
	return p.namespace
}
func (in *Handler) DeepCopy() *Handler {
	out := new(Handler)

	out.kubeconfig = in.kubeconfig
	out.namespace = in.namespace

	out.ctx = in.ctx
	out.config = in.config
	out.restClient = in.restClient
	out.clientset = in.clientset
	out.dynamicClient = in.dynamicClient
	out.discoveryClient = in.discoveryClient
	out.informerFactory = in.informerFactory
	out.informer = in.informer
	out.client = in.client

	out.Options = &typed.HandlerOptions{}
	out.Options.ListOptions = *in.Options.ListOptions.DeepCopy()
	out.Options.GetOptions = *in.Options.GetOptions.DeepCopy()
	out.Options.CreateOptions = *in.Options.CreateOptions.DeepCopy()
	out.Options.UpdateOptions = *in.Options.UpdateOptions.DeepCopy()
	out.Options.PatchOptions = *in.Options.PatchOptions.DeepCopy()
	out.Options.ApplyOptions = *in.Options.ApplyOptions.DeepCopy()

	return out
}

func (p *Handler) setNamespace(namespace string) {
	p.Lock()
	defer p.Unlock()
	p.namespace = namespace
}
func (p *Handler) WithNamespace(namespace string) *Handler {
	handler := p.DeepCopy()
	handler.setNamespace(namespace)
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
	p.Lock()
	defer p.Unlock()
	p.Options.ListOptions.Limit = limit
}
func (p *Handler) SetTimeout(timeout int64) {
	p.Lock()
	defer p.Unlock()
	p.Options.ListOptions.TimeoutSeconds = &timeout
}
func (p *Handler) SetForceDelete(force bool) {
	p.Lock()
	defer p.Unlock()
	if force {
		gracePeriodSeconds := int64(0)
		p.Options.DeleteOptions.GracePeriodSeconds = &gracePeriodSeconds
	} else {
		p.Options.DeleteOptions = metav1.DeleteOptions{}
	}
}
