package ingressclass

import (
	"fmt"
	"time"

	networkingv1 "k8s.io/api/networking/v1"
)

var ERR_TYPE = fmt.Errorf("type must be *networkingv1.IngressClass, networkingv1.IngressClass or string")

func (h *Handler) GetAge(object interface{}) (time.Duration, error) {
	switch val := object.(type) {
	case string:
		ingc, err := h.Get(val)
		if err != nil {
			return time.Duration(int64(0)), err
		}
		return time.Now().Sub(ingc.CreationTimestamp.Time), nil
	case *networkingv1.IngressClass:
		return time.Now().Sub(val.CreationTimestamp.Time), nil
	case networkingv1.IngressClass:
		return time.Now().Sub(val.CreationTimestamp.Time), nil
	default:
		return time.Duration(int64(0)), ERR_TYPE
	}
}
