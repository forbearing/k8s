package networkpolicy

import (
	"time"

	networkingv1 "k8s.io/api/networking/v1"
)

// GetAge get the networkpolicy age.
func (h *Handler) GetAge(object interface{}) (time.Duration, error) {
	switch val := object.(type) {
	case string:
		ns, err := h.Get(val)
		if err != nil {
			return time.Duration(0), err
		}
		return time.Now().Sub(ns.CreationTimestamp.Time), nil
	case *networkingv1.NetworkPolicy:
		return time.Now().Sub(val.CreationTimestamp.Time), nil
	case networkingv1.NetworkPolicy:
		return time.Now().Sub(val.CreationTimestamp.Time), nil
	default:
		return time.Duration(0), ErrInvalidToolsType
	}
}
