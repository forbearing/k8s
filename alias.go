package k8s

import (
	"github.com/forbearing/k8s/util/client"
)

// RESTConfig creates a *rest.Config for the given kubeconfig.
// create rest config, and config precedence.
// * kubeconfig variable passed.
// * KUBECONFIG environment variable pointing at a file
// * $HOME/.kube/config if exists.
// * In-cluster config if running in cluster
var RESTConfig = client.RESTConfig

// RESTConfigOrDie creates a *rest.Config for the given kubeconfig.
// panic if there is any error occurs.
var RESTConfigOrDie = client.RESTConfigOrDie

// RESTClient creates a *rest.RESTClient for the given kubeconfig.
var RESTClient = client.RESTClient

// RESTClientOrDie creates a *rest.RESTClient for the given kubeconfig.
// panic if there is any error occurs.
var RESTClientOrDie = client.RESTClientOrDie

// Clientset creates a *kubernetes.Clientset for the given kubeconfig.
var Clientset = client.Clientset

// ClientsetOrDie creates a *kubernetes.Clientset for the given kubeconfig.
// panic if there is any error occurs.
var ClientsetOrDie = client.ClientsetOrDie

// DynamicClient creates a dynamic.Interface for the given kubeconfig.
var DynamicClient = client.DynamicClient

// DynamicClient creates a dynamic.Interface for the given kubeconfig.
// panic if there is any error occurs.
var DynamicClientOrDie = client.DynamicClientOrDie

// DiscoveryClient creates a *discovery.DiscoveryClient for the given kubeconfig.
var DiscoveryClient = client.DiscoveryClient

// DiscoveryClientOrDie creates a *discovery.DiscoveryClient for the given kubeconfig.
// panic if there is any error occurs.
var DiscoveryClientOrDie = client.DiscoveryClientOrDie

// RawConfig holds the information needed to build connect to remote kubernetes
// clusters as a given user
//
// ref: https://stackoverflow.com/questions/70885022/how-to-get-current-k8s-context-name-using-client-go
var RawConfig = client.RawConfig
