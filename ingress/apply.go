package ingress

import (
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
)

// ApplyFromRaw apply ingress from map[string]interface{}.
func (h *Handler) ApplyFromRaw(raw map[string]interface{}) (*networkingv1.Ingress, error) {
	ingress := &networkingv1.Ingress{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(raw, ingress)
	if err != nil {
		return nil, err
	}

	var namespace string
	if len(ingress.Namespace) != 0 {
		namespace = ingress.Namespace
	} else {
		namespace = h.namespace
	}

	ingress, err = h.clientset.NetworkingV1().Ingresses(namespace).Create(h.ctx, ingress, h.Options.CreateOptions)
	if k8serrors.IsAlreadyExists(err) {
		ingress, err = h.clientset.NetworkingV1().Ingresses(namespace).Update(h.ctx, ingress, h.Options.UpdateOptions)
	}
	return ingress, err
}

// ApplyFromBytes apply ingress from bytes.
func (h *Handler) ApplyFromBytes(data []byte) (ingress *networkingv1.Ingress, err error) {
	ingress, err = h.CreateFromBytes(data)
	if errors.IsAlreadyExists(err) {
		ingress, err = h.UpdateFromBytes(data)
	}
	return
}

// ApplyFromFile apply ingress from yaml file.
func (h *Handler) ApplyFromFile(filename string) (ingress *networkingv1.Ingress, err error) {
	ingress, err = h.CreateFromFile(filename)
	if errors.IsAlreadyExists(err) {
		ingress, err = h.UpdateFromFile(filename)
	}
	return
}

// Apply apply ingress from file, alias to "ApplyFromFile".
func (h *Handler) Apply(filename string) (*networkingv1.Ingress, error) {
	return h.ApplyFromFile(filename)
}
