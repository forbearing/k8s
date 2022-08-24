package ingressclass

import (
	"time"

	networkingv1 "k8s.io/api/networking/v1"
)

// GetAge
func (h *Handler) GetAge(object interface{}) (time.Duration, error) {
	switch val := object.(type) {
	case string:
		ingc, err := h.Get(val)
		if err != nil {
			return time.Duration(0), err
		}
		return time.Now().Sub(ingc.CreationTimestamp.Time), nil
	case *networkingv1.IngressClass:
		return time.Now().Sub(val.CreationTimestamp.Time), nil
	case networkingv1.IngressClass:
		return time.Now().Sub(val.CreationTimestamp.Time), nil
	default:
		return time.Duration(0), ErrInvalidToolsType
	}
}
