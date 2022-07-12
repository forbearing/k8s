package ingress

import (
	"fmt"
	"time"

	networkingv1 "k8s.io/api/networking/v1"
)

var ERR_TYPE = fmt.Errorf("type must be *networkingv1.Ingress, networkingv1.Ingress or string")

// getHosts
func (h *Handler) getHosts(ing *networkingv1.Ingress) []string {
	var hl []string
	for _, rule := range ing.Spec.Rules {
		hl = append(hl, rule.Host)
	}
	return hl
}

// GetClass get ingressclass of the ingress.
func (h *Handler) GetClass(object interface{}) (string, error) {
	switch val := object.(type) {
	case string:
		ing, err := h.Get(val)
		if err != nil {
			return "", nil
		}
		if ing.Spec.IngressClassName == nil {
			return "", fmt.Errorf("ingress/%s don't have IngressClassName", val)
		}
		return *(ing.Spec.IngressClassName), nil
	case *networkingv1.Ingress:
		if val.Spec.IngressClassName == nil {
			return "", fmt.Errorf("ingress/%s don't have IngressClassName", val.Name)
		}
		return *(val.Spec.IngressClassName), nil
	case networkingv1.Ingress:
		if val.Spec.IngressClassName == nil {
			return "", fmt.Errorf("ingress/%s don't have IngressClassName", val.Name)
		}
		return *(val.Spec.IngressClassName), nil
	default:
		return "", ERR_TYPE
	}
}

// GetHosts get ingress server host
func (h *Handler) GetHosts(object interface{}) ([]string, error) {
	switch val := object.(type) {
	case string:
		ing, err := h.Get(val)
		if err != nil {
			return nil, err
		}
		return h.getHosts(ing), nil
	case *networkingv1.Ingress:
		return h.getHosts(val), nil
	case networkingv1.Ingress:
		return h.getHosts(&val), nil
	default:
		return nil, ERR_TYPE
	}
}

func (h *Handler) getAddress(ing *networkingv1.Ingress) []string {
	var al []string
	if ing.Status.LoadBalancer.Ingress == nil {
		return nil
	}
	for _, i := range ing.Status.LoadBalancer.Ingress {
		al = append(al, i.IP)
	}
	return al
}

// GetAddress
func (h *Handler) GetAddress(object interface{}) ([]string, error) {
	switch val := object.(type) {
	case string:
		ing, err := h.Get(val)
		if err != nil {
			return nil, err
		}
		return h.getAddress(ing), nil
	case *networkingv1.Ingress:
		return h.getAddress(val), nil
	case networkingv1.Ingress:
		return h.getAddress(&val), nil
	default:
		return nil, ERR_TYPE
	}
}

// GetAge get the ingress age.
func (h *Handler) GetAge(object interface{}) (time.Duration, error) {
	switch val := object.(type) {
	case string:
		ing, err := h.Get(val)
		if err != nil {
			return time.Duration(int64(0)), err
		}
		return time.Now().Sub(ing.CreationTimestamp.Time), nil
	case *networkingv1.Ingress:
		return time.Now().Sub(val.CreationTimestamp.Time), nil
	case networkingv1.Ingress:
		return time.Now().Sub(val.CreationTimestamp.Time), nil
	default:
		return time.Duration(int64(0)), ERR_TYPE
	}
}
