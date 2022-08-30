package dynamic

import (
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic/dynamicinformer"
)

// SetInformerFactoryResyncPeriod will set informer resync period.
func (h *Handler) SetInformerFactoryResyncPeriod(resyncPeriod time.Duration) {
	h.l.Lock()
	defer h.l.Unlock()
	h.resyncPeriod = resyncPeriod
	if len(h.informerScope) == 0 {
		h.informerScope = metav1.NamespaceAll
	}
	h.informerFactory = dynamicinformer.NewFilteredDynamicSharedInformerFactory(
		h.dynamicClient, h.resyncPeriod, h.informerScope, h.tweakListOptions)
}

// SetInformerFactoryNamespace limit the scope of informer list-and-watch k8s resource.
// informer list-and-watch all namespace k8s resource by default.
func (h *Handler) SetInformerFactoryNamespace(namespace string) {
	h.l.Lock()
	defer h.l.Unlock()
	h.informerScope = namespace
	if len(h.informerScope) == 0 {
		h.informerScope = metav1.NamespaceAll
	}
	h.informerFactory = dynamicinformer.NewFilteredDynamicSharedInformerFactory(
		h.dynamicClient, h.resyncPeriod, h.informerScope, h.tweakListOptions)
}

// SetInformerFactoryTweakListOptions sets a custom filter on all listers of
// the configured SharedInformerFactory.
func (h *Handler) SetInformerFactoryTweakListOptions(tweakListOptions dynamicinformer.TweakListOptionsFunc) {
	h.l.Lock()
	defer h.l.Unlock()
	h.tweakListOptions = tweakListOptions
	if len(h.informerScope) == 0 {
		h.informerScope = metav1.NamespaceAll
	}
	h.informerFactory = dynamicinformer.NewFilteredDynamicSharedInformerFactory(
		h.dynamicClient, h.resyncPeriod, h.informerScope, h.tweakListOptions)
}

// InformerFactory returns underlying DyanmicSharedInformerFactory which provides
// access to a shared informer and lister for dynamic client
func (h *Handler) InformerFactory() dynamicinformer.DynamicSharedInformerFactory {
	return h.informerFactory
}
