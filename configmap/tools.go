package configmap

import (
	"time"

	corev1 "k8s.io/api/core/v1"
)

// GetData get configmap data.
func (h *Handler) GetData(object interface{}) (map[string]string, error) {
	switch val := object.(type) {
	case string:
		cm, err := h.Get(val)
		if err != nil {
			return nil, err
		}
		return cm.Data, nil
	case *corev1.ConfigMap:
		return val.Data, nil
	case corev1.ConfigMap:
		return val.Data, nil
	default:
		return nil, ErrInvalidToolsType
	}
}

// NumData get the number of configmap data.
func (h *Handler) NumData(object interface{}) (int, error) {
	switch val := object.(type) {
	case string:
		cm, err := h.Get(val)
		if err != nil {
			return 0, err
		}
		return len(cm.Data), nil
	case *corev1.ConfigMap:
		return len(val.Data), nil
	case corev1.ConfigMap:
		return len(val.Data), nil
	default:
		return 0, ErrInvalidToolsType
	}
}

// GetAge returns configmap age.
func (h *Handler) GetAge(object interface{}) (time.Duration, error) {
	switch val := object.(type) {
	case string:
		cm, err := h.Get(val)
		if err != nil {
			return time.Duration(0), err
		}
		return time.Now().Sub(cm.CreationTimestamp.Time), nil
	case *corev1.ConfigMap:
		return time.Now().Sub(val.CreationTimestamp.Time), nil
	case corev1.ConfigMap:
		return time.Now().Sub(val.CreationTimestamp.Time), nil
	default:
		return time.Duration(0), ErrInvalidToolsType
	}
}
