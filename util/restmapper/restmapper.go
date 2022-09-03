package restmapper

import (
	"github.com/forbearing/k8s/util/client"
	"k8s.io/apimachinery/pkg/api/meta"
	memory "k8s.io/client-go/discovery/cached"
	"k8s.io/client-go/restmapper"
)

// NewRESTMapper alias to NewDeferredRESTMapper
var NewRESTMapper = NewDeferredRESTMapper

// NewDeferredRESTMapper
func NewDeferredRESTMapper(kubeconfig string) (meta.RESTMapper, error) {
	discoveryClient, err := client.DiscoveryClient(kubeconfig)
	if err != nil {
		return nil, err
	}

	// NewMemCacheClient creates a new CachedDiscoveryInterface which caches
	// discovery information in memory and will stay up-to-date if Invalidate is
	// called with regularity.
	//
	// NewDeferredDiscoveryRESTMapper returns a
	// DeferredDiscoveryRESTMapper that will lazily query the provided
	// client for discovery information to do REST mappings.
	return restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(discoveryClient)), nil
}

// NewPriorityRESTMapper
func NewPriorityRESTMapper(kubeconfig string) (meta.RESTMapper, error) {
	discoveryClient, err := client.DiscoveryClient(kubeconfig)
	if err != nil {
		return nil, err
	}
	// GetAPIGroupResources uses the provided discovery client to gather
	// discovery information and populate a slice of APIGroupResources.
	apiGroupResources, err := restmapper.GetAPIGroupResources(discoveryClient)
	if err != nil {
		return nil, err
	}
	// NewDiscoveryRESTMapper returns a PriorityRESTMapper based on the discovered
	// groups and resources passed in.
	return restmapper.NewDiscoveryRESTMapper(apiGroupResources), nil
}
