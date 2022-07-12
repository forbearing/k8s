package clusterrole

import (
	"fmt"
	"time"

	rabcv1 "k8s.io/api/rbac/v1"
)

var ERR_TYPE = fmt.Errorf("type must be *rbacv1.ClusterRole, rbacv1.ClusterRole or string")

// GetAge returns clusterrole age.
func (h *Handler) GetAge(object interface{}) (time.Duration, error) {
	switch val := object.(type) {
	case string:
		cr, err := h.Get(val)
		if err != nil {
			return time.Duration(int64(0)), err
		}
		return time.Now().Sub(cr.CreationTimestamp.Time), nil
	case *rabcv1.ClusterRole:
		return time.Now().Sub(val.CreationTimestamp.Time), nil
	case rabcv1.ClusterRole:
		return time.Now().Sub(val.CreationTimestamp.Time), nil
	default:
		return time.Duration(int64(0)), ERR_TYPE
	}
}
