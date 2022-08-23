package role

import (
	"time"

	rabcv1 "k8s.io/api/rbac/v1"
)

// GetAge returns role age.
func (h *Handler) GetAge(object interface{}) (time.Duration, error) {
	switch val := object.(type) {
	case string:
		role, err := h.Get(val)
		if err != nil {
			return time.Duration(int64(0)), err
		}
		return time.Now().Sub(role.CreationTimestamp.Time), nil
	case *rabcv1.Role:
		return time.Now().Sub(val.CreationTimestamp.Time), nil
	case rabcv1.Role:
		return time.Now().Sub(val.CreationTimestamp.Time), nil
	default:
		return time.Duration(int64(0)), ErrInvalidToolsType
	}
}
