package secret

import (
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
)

var ERR_TYPE = fmt.Errorf("type must be *corev1.Secret, corev1.Secret or string")

// GetType returns the secret type.
func (h *Handler) GetType(object interface{}) (string, error) {
	switch val := object.(type) {
	case string:
		secret, err := h.Get(val)
		if err != nil {
			return "", err
		}
		return string(secret.Type), nil
	case *corev1.Secret:
		return string(val.Type), nil
	case corev1.Secret:
		return string(val.Type), nil
	default:
		return "", ERR_TYPE
	}
}

// GetNumData returns the number of the secret.
func (h *Handler) GetNumData(object interface{}) (int, error) {
	switch val := object.(type) {
	case string:
		secret, err := h.Get(val)
		if err != nil {
			return 0, err
		}
		return len(secret.Data), nil
	case *corev1.Secret:
		return len(val.Data), nil
	case corev1.Secret:
		return len(val.Data), nil
	default:
		return 0, ERR_TYPE
	}
}

// GetAge returns the age of the secret.
func (h *Handler) GetAge(object interface{}) (time.Duration, error) {
	switch val := object.(type) {
	case string:
		secret, err := h.Get(val)
		if err != nil {
			return time.Duration(int64(0)), err
		}
		return time.Now().Sub(secret.CreationTimestamp.Time), nil
	case *corev1.Secret:
		return time.Now().Sub(val.CreationTimestamp.Time), nil
	case corev1.Secret:
		return time.Now().Sub(val.CreationTimestamp.Time), nil
	default:
		return time.Duration(int64(0)), ERR_TYPE
	}
}
