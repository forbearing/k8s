package ingressclass

import (
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
)

// ApplyFromRaw apply ingressclass from map[string]interface{}.
func (h *Handler) ApplyFromRaw(raw map[string]interface{}) (*networkingv1.IngressClass, error) {
	ingressclass := &networkingv1.IngressClass{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(raw, ingressclass)
	if err != nil {
		return nil, err
	}

	ingressclass, err = h.clientset.NetworkingV1().IngressClasses().Create(h.ctx, ingressclass, h.Options.CreateOptions)
	if k8serrors.IsAlreadyExists(err) {
		ingressclass, err = h.clientset.NetworkingV1().IngressClasses().Update(h.ctx, ingressclass, h.Options.UpdateOptions)
	}
	return ingressclass, err
}

// ApplyFromBytes apply ingressclass from bytes.
func (h *Handler) ApplyFromBytes(data []byte) (ingc *networkingv1.IngressClass, err error) {
	ingc, err = h.CreateFromBytes(data)
	if errors.IsAlreadyExists(err) {
		ingc, err = h.UpdateFromBytes(data)
	}
	return
}

// ApplyFromFile apply ingressclass from yaml file.
func (h *Handler) ApplyFromFile(filename string) (ingc *networkingv1.IngressClass, err error) {
	ingc, err = h.CreateFromFile(filename)
	if errors.IsAlreadyExists(err) {
		ingc, err = h.UpdateFromFile(filename)
	}
	return
}

// Apply apply ingressclass from yaml file, alias to "ApplyFromFile".
func (h *Handler) Apply(filename string) (*networkingv1.IngressClass, error) {
	return h.ApplyFromFile(filename)
}
