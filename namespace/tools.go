package namespace

import (
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
)

var ERR_TYPE = fmt.Errorf("type must be *corev1.Namespace, corev1.Namespace or string")

func (h *Handler) GetAge(object interface{}) (time.Duration, error) {
	switch val := object.(type) {
	case string:
		ns, err := h.Get(val)
		if err != nil {
			return time.Duration(int64(0)), nil
		}
		return time.Now().Sub(ns.CreationTimestamp.Time), nil
	case *corev1.Namespace:
		return time.Now().Sub(val.CreationTimestamp.Time), nil
	case corev1.Namespace:
		return time.Now().Sub(val.CreationTimestamp.Time), nil
	default:
		return time.Duration(int64(0)), ERR_TYPE
	}
}
