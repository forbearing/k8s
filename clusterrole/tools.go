package clusterrole

import (
	"time"

	rabcv1 "k8s.io/api/rbac/v1"
)

// GetAge returns clusterrole age.
func (h *Handler) GetAge(object interface{}) (time.Duration, error) {
	switch val := object.(type) {
	case string:
		cr, err := h.Get(val)
		if err != nil {
			return time.Duration(0), err
		}
		return time.Now().Sub(cr.CreationTimestamp.Time), nil
	case *rabcv1.ClusterRole:
		return time.Now().Sub(val.CreationTimestamp.Time), nil
	case rabcv1.ClusterRole:
		return time.Now().Sub(val.CreationTimestamp.Time), nil
	default:
		return time.Duration(0), ErrInvalidToolsType
	}
}
