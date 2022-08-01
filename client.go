package k8s

import (
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// RESTConfig creates a *rest.Config for the given kubeconfig.
func RESTConfig(kubeconfig string) (*rest.Config, error) {
	var config *rest.Config
	var err error
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
	return config, nil
}

// RESTConfigOrDie creates a *rest.Config for the given kubeconfig.
// panic if there is any error occurs.
func RESTConfigOrDie(kubeconfig string) *rest.Config {
	config, err := RESTConfig(kubeconfig)
	if err != nil {
		panic(err)
	}
	return config
}

// RESTClient creates a *rest.RESTClient for the given kubeconfig.
func RESTClient(kubeconfig string) (*rest.RESTClient, error) {
	config, err := RESTConfig(kubeconfig)
	if err != nil {
		return nil, err
	}
	restClient, err := rest.RESTClientFor(config)
	if err != nil {
		return nil, err
	}
	return restClient, nil
}

// RESTClientOrDie creates a *rest.RESTClient for the given kubeconfig.
// panic if there is any error occurs.
func RESTClientOrDie(kubeconfig string) *rest.RESTClient {
	restClient, err := RESTClient(kubeconfig)
	if err != nil {
		panic(err)
	}
	return restClient
}

// Clientset creates a *kubernetes.Clientset for the given kubeconfig.
func Clientset(kubeconfig string) (*kubernetes.Clientset, error) {
	config, err := RESTConfig(kubeconfig)
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return clientset, nil
}

// ClientsetOrDie creates a *kubernetes.Clientset for the given kubeconfig.
// panic if there is any error occurs.
func ClientsetOrDie(kubeconfig string) *kubernetes.Clientset {
	clientset, err := Clientset(kubeconfig)
	if err != nil {
		panic(err)
	}
	return clientset
}

// DynamicClient creates a dynamic.Interface for the given kubeconfig.
func DynamicClient(kubeconfig string) (dynamic.Interface, error) {
	config, err := RESTConfig(kubeconfig)
	if err != nil {
		return nil, err
	}
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return dynamicClient, nil
}

// DynamicClient creates a dynamic.Interface for the given kubeconfig.
// panic if there is any error occurs.
func DynamicClientOrDie(kubeconfig string) dynamic.Interface {
	dynamicClient, err := DynamicClient(kubeconfig)
	if err != nil {
		panic(err)
	}
	return dynamicClient
}

// DiscoveryClient creates a *discovery.DiscoveryClient for the given kubeconfig.
func DiscoveryClient(kubeconfig string) (*discovery.DiscoveryClient, error) {
	config, err := RESTConfig(kubeconfig)
	if err != nil {
		return nil, err
	}
	discoveryClient, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		return nil, err
	}
	return discoveryClient, nil
}

// DiscoveryClientOrDie creates a *discovery.DiscoveryClient for the given kubeconfig.
// panic if there is any error occurs.
func DiscoveryClientOrDie(kubeconfig string) *discovery.DiscoveryClient {
	discoveryClient, err := DiscoveryClient(kubeconfig)
	if err != nil {
		panic(err)
	}
	return discoveryClient
}
