package networkpolicy

import (
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
)

// ApplyFromRaw apply netpol from map[string]interface{}.
func (h *Handler) ApplyFromRaw(raw map[string]interface{}) (*networkingv1.NetworkPolicy, error) {
	netpol := &networkingv1.NetworkPolicy{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(raw, netpol)
	if err != nil {
		return nil, err
	}

	var namespace string
	if len(netpol.Namespace) != 0 {
		namespace = netpol.Namespace
	} else {
		namespace = h.namespace
	}

	netpol, err = h.clientset.NetworkingV1().NetworkPolicies(namespace).Create(h.ctx, netpol, h.Options.CreateOptions)
	if k8serrors.IsAlreadyExists(err) {
		netpol, err = h.clientset.NetworkingV1().NetworkPolicies(namespace).Update(h.ctx, netpol, h.Options.UpdateOptions)
	}
	return netpol, err
}

// ApplyFromBytes apply networkpolicy from bytes.
func (h *Handler) ApplyFromBytes(data []byte) (netpol *networkingv1.NetworkPolicy, err error) {
	netpol, err = h.CreateFromBytes(data)
	if errors.IsAlreadyExists(err) {
		netpol, err = h.UpdateFromBytes(data)
	}
	return
}

// ApplyFromFile apply netpol from yaml file.
func (h *Handler) ApplyFromFile(filename string) (netpol *networkingv1.NetworkPolicy, err error) {
	netpol, err = h.CreateFromFile(filename)
	if errors.IsAlreadyExists(err) {
		netpol, err = h.UpdateFromFile(filename)
	}
	return
}

// Apply apply networkpolicy from yaml file, alias to "ApplyFromFile".
func (h *Handler) Apply(filename string) (*networkingv1.NetworkPolicy, error) {
	return h.ApplyFromFile(filename)
}
