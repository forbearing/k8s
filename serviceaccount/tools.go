package serviceaccount

import (
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
)

var ERR_TYPE = fmt.Errorf("type must be *corev1.ServiceAccount, corev1.ServiceAccount or string")

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
		return 0, ERR_TYPE
	}
}

// GetAge returns the age of the serviceaccount.
func (h *Handler) GetAge(object interface{}) (time.Duration, error) {
	switch val := object.(type) {
	case string:
		sa, err := h.Get(val)
		if err != nil {
			return time.Duration(int64(0)), err
		}
		return time.Now().Sub(sa.CreationTimestamp.Time), nil
	case *corev1.ServiceAccount:
		return time.Now().Sub(val.CreationTimestamp.Time), nil
	case corev1.ServiceAccount:
		return time.Now().Sub(val.CreationTimestamp.Time), nil
	default:
		return time.Duration(int64(0)), ERR_TYPE
	}
}
