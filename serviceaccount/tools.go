package serviceaccount

import (
	"time"

	corev1 "k8s.io/api/core/v1"
)

// GetNumSecrets get the number of secret referenced by this serviceaccount.
func (h *Handler) GetNumSecrets(object interface{}) (int, error) {
	switch val := object.(type) {
	case string:
		sa, err := h.Get(val)
		if err != nil {
			return 0, err
		}
		return len(sa.Secrets), nil
	case *corev1.ServiceAccount:
		return len(val.Secrets), nil
	case corev1.ServiceAccount:
		return len(val.Secrets), nil
	default:
		return 0, ErrInvalidToolsType
	}
}

// GetAge returns the age of the serviceaccount.
func (h *Handler) GetAge(object interface{}) (time.Duration, error) {
	switch val := object.(type) {
	case string:
		sa, err := h.Get(val)
		if err != nil {
			return time.Duration(0), err
		}
		return time.Now().Sub(sa.CreationTimestamp.Time), nil
	case *corev1.ServiceAccount:
		return time.Now().Sub(val.CreationTimestamp.Time), nil
	case corev1.ServiceAccount:
		return time.Now().Sub(val.CreationTimestamp.Time), nil
	default:
		return time.Duration(0), ErrInvalidToolsType
	}
}
