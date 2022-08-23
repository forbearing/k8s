package storageclass

import (
	"time"

	storagev1 "k8s.io/api/storage/v1"
)

// GetProvisioner get the provisioner of the storageclass.
func (h *Handler) GetProvisioner(object interface{}) (string, error) {
	switch val := object.(type) {
	case string:
		sc, err := h.Get(val)
		if err != nil {
			return "", err
		}
		return sc.Provisioner, nil
	case *storagev1.StorageClass:
		return val.Provisioner, nil
	case storagev1.StorageClass:
		return val.Provisioner, nil
	default:
		return "", ErrInvalidToolsType
	}
}

// GetReclaimPolicy get the  reclaim policy of the storageclass.
func (h *Handler) GetReclaimPolicy(object interface{}) (string, error) {
	switch val := object.(type) {
	case string:
		sc, err := h.Get(val)
		if err != nil {
			return "", err
		}
		return string(*sc.ReclaimPolicy), nil
	case *storagev1.StorageClass:
		return string(*val.ReclaimPolicy), nil
	case storagev1.StorageClass:
		return string(*val.ReclaimPolicy), nil
	default:
		return "", ErrInvalidToolsType
	}
}

// IsAllowVolumeExpansion checks whether the storageclass allow volume expand.
func (h *Handler) IsAllowVolumeExpansion(object interface{}) (bool, error) {
	switch val := object.(type) {
	case string:
		sc, err := h.Get(val)
		if err != nil {
			return false, err
		}
		// if AllowVolumeExpansion not set, set it to false.
		if sc.AllowVolumeExpansion == nil {
			return false, nil
			//return false, fmt.Errorf("AllowVolumeExpansion field not set")
		}
		return *sc.AllowVolumeExpansion, nil
	case *storagev1.StorageClass:
		if val.AllowVolumeExpansion == nil {
			return false, nil
		}
		return *val.AllowVolumeExpansion, nil
	case storagev1.StorageClass:
		if val.AllowVolumeExpansion == nil {
			return false, nil
		}
		return *val.AllowVolumeExpansion, nil
	default:
		return false, ErrInvalidToolsType
	}
}

// GetAge returns this storageclass age.
func (h *Handler) GetAge(object interface{}) (time.Duration, error) {
	switch val := object.(type) {
	case string:
		sts, err := h.Get(val)
		if err != nil {
			return time.Duration(int64(0)), err
		}
		ctime := sts.CreationTimestamp.Time
		return time.Now().Sub(ctime), nil
	case *storagev1.StorageClass:
		ctime := val.CreationTimestamp.Time
		return time.Now().Sub(ctime), nil
	case storagev1.StorageClass:
		ctime := val.CreationTimestamp.Time
		return time.Now().Sub(ctime), nil
	default:
		return time.Duration(int64(0)), ErrInvalidToolsType
	}
}
