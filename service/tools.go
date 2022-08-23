package service

import (
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type ServicePort struct {
	Name       string
	Protocol   string
	Port       int32
	TargetPort intstr.IntOrString
	NodePort   int32
}

// GetType get the type of the service.
func (h *Handler) GetType(object interface{}) (string, error) {
	switch val := object.(type) {
	case string:
		svc, err := h.Get(val)
		if err != nil {
			return "", err
		}
		return string(svc.Spec.Type), nil
	case *corev1.Service:
		return string(val.Spec.Type), nil
	case corev1.Service:
		return string(val.Spec.Type), nil
	default:
		return "", ErrInvalidToolsType
	}
}

// GetClusterIP returns the cluster ip of the service.
func (h *Handler) GetClusterIP(object interface{}) (string, error) {
	switch val := object.(type) {
	case string:
		svc, err := h.Get(val)
		if err != nil {
			return "", err
		}
		return svc.Spec.ClusterIP, nil
	case *corev1.Service:
		return val.Spec.ClusterIP, nil
	case corev1.Service:
		return val.Spec.ClusterIP, nil
	default:
		return "", ErrInvalidToolsType
	}
}

// GetExternalIPs returns the external ip list of the service.
func (h *Handler) GetExternalIPs(object interface{}) ([]string, error) {
	switch val := object.(type) {
	case string:
		svc, err := h.Get(val)
		if err != nil {
			return nil, err
		}
		return svc.Spec.ExternalIPs, nil
	case *corev1.Service:
		return val.Spec.ExternalIPs, nil
	case corev1.Service:
		return val.Spec.ExternalIPs, nil
	default:
		return nil, ErrInvalidToolsType
	}
}

// GetPorts get the port list of the service.
func (h *Handler) GetPorts(object interface{}) ([]ServicePort, error) {
	switch val := object.(type) {
	case string:
		svc, err := h.Get(val)
		if err != nil {
			return nil, err
		}
		return h.getPorts(svc), nil
	case *corev1.Service:
		return h.getPorts(val), nil
	case corev1.Service:
		return h.getPorts(&val), nil
	default:
		return nil, ErrInvalidToolsType
	}
}
func (h *Handler) getPorts(svc *corev1.Service) []ServicePort {
	var pl []ServicePort

	for _, port := range svc.Spec.Ports {
		p := ServicePort{
			Name:       port.Name,
			Protocol:   string(port.Protocol),
			Port:       port.Port,
			TargetPort: port.TargetPort,
			NodePort:   port.NodePort,
		}
		pl = append(pl, p)
	}
	return pl
}

// GetAge returns the age of the service.
func (h *Handler) GetAge(object interface{}) (time.Duration, error) {
	switch val := object.(type) {
	case string:
		svc, err := h.Get(val)
		if err != nil {
			return time.Duration(int64(0)), err
		}
		return time.Now().Sub(svc.CreationTimestamp.Time), nil
	case *corev1.Service:
		return time.Now().Sub(val.CreationTimestamp.Time), nil
	case corev1.Service:
		return time.Now().Sub(val.CreationTimestamp.Time), nil
	default:
		return time.Duration(int64(0)), ErrInvalidToolsType
	}
}
