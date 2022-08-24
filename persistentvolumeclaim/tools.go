package persistentvolumeclaim

import (
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
)

// GetStatus get the status phase of the persistentvolumeclaim.
// All supported status are: Pending, Bound, Lost
// pending used for PersistentVolumeClaims that are not yet bound.
// Bound used for PersistentVolumeClaims that are bound.
// Lost used for PersistentVolumeClaims that lost their underlying.
func (h *Handler) GetStatus(object interface{}) (phase string, err error) {
	switch val := object.(type) {
	case string:
		pvc, err := h.Get(val)
		if err != nil {
			return "", err
		}
		return string(pvc.Status.Phase), nil
	case *corev1.PersistentVolumeClaim:
		return string(val.Status.Phase), nil
	case corev1.PersistentVolumeClaim:
		return string(val.Status.Phase), nil
	default:
		return "", ErrInvalidToolsType
	}
}

// GetVolume simply calls GetPV.
func (h *Handler) GetVolume(object interface{}) (string, error) {
	return h.GetPV(object)
}

// GetPV get the pv name of the persistentvolumeclaim
func (h *Handler) GetPV(object interface{}) (string, error) {
	switch val := object.(type) {
	case string:
		pvc, err := h.Get(val)
		if err != nil {
			return "", err
		}
		return pvc.Spec.VolumeName, nil
	case *corev1.PersistentVolumeClaim:
		return val.Spec.VolumeName, nil
	case corev1.PersistentVolumeClaim:
		return val.Spec.VolumeName, nil
	default:
		return "", ErrInvalidToolsType
	}
}

// GetCapacity get the storage capacity of the persistentvolumeclaim.
func (h *Handler) GetCapacity(object interface{}) (int64, error) {
	switch val := object.(type) {
	case string:
		pvc, err := h.Get(val)
		if err != nil {
			return 0, err
		}
		return h.getCapacity(pvc), nil
	case *corev1.PersistentVolumeClaim:
		return h.getCapacity(val), nil
	case corev1.PersistentVolumeClaim:
		return h.getCapacity(&val), nil
	default:
		return 0, ErrInvalidToolsType
	}
}
func (h *Handler) getCapacity(pvc *corev1.PersistentVolumeClaim) int64 {
	storage := pvc.Status.Capacity[corev1.ResourceName(corev1.ResourceStorage)]
	//capacity = storage.Value()
	//capacity = storage.MilliValue()
	//capacity = storage.ScaledValue(resource.Kilo)
	//capacity = storage.ScaledValue(resource.Mega)
	//capacity = storage.ScaledValue(resource.Giga)
	//capacity = storage.ScaledValue(resource.Tera)
	//capacity = storage.ScaledValue(resource.Peta)
	//capacity = storage.ScaledValue(resource.Exa)
	return storage.Value()
}

// GetAccessModes get the access modes of the persistentvolumeclaim.
func (h *Handler) GetAccessModes(object interface{}) ([]string, error) {
	switch val := object.(type) {
	case string:
		pvc, err := h.Get(val)
		if err != nil {
			return nil, err
		}
		return h.getAccessModes(pvc), nil
	case *corev1.PersistentVolumeClaim:
		return h.getAccessModes(val), nil
	case corev1.PersistentVolumeClaim:
		return h.getAccessModes(&val), nil
	default:
		return nil, ErrInvalidToolsType
	}
}
func (h *Handler) getAccessModes(pvc *corev1.PersistentVolumeClaim) []string {
	var accessModes []string
	for _, accessMode := range pvc.Status.AccessModes {
		accessModes = append(accessModes, string(accessMode))
	}
	return accessModes
}

// GetStorageClass get the storageclass name of the persistentvolumeclaim.
func (h *Handler) GetStorageClass(object interface{}) (string, error) {
	switch val := object.(type) {
	case string:
		pvc, err := h.Get(val)
		if err != nil {
			return "", err
		}
		if pvc.Spec.StorageClassName == nil {
			return "", fmt.Errorf("pvc/%s doesn't have storageclass", pvc.Name)
		}
		return *(pvc.Spec.StorageClassName), nil
	case *corev1.PersistentVolumeClaim:
		if val.Spec.StorageClassName == nil {
			return "", fmt.Errorf("pvc/%s doesn't have storageclass", val.Name)
		}
		return *(val.Spec.StorageClassName), nil
	case corev1.PersistentVolumeClaim:
		if val.Spec.StorageClassName == nil {
			return "", fmt.Errorf("pvc/%s doesn't have storageclass", val.Name)
		}
		return *(val.Spec.StorageClassName), nil
	default:
		return "", ErrInvalidToolsType
	}
}

// GetAge returns age of the persistentvolumeclaim.
func (h *Handler) GetAge(object interface{}) (time.Duration, error) {
	switch val := object.(type) {
	case string:
		pvc, err := h.Get(val)
		if err != nil {
			return time.Duration(0), err
		}
		return time.Now().Sub(pvc.CreationTimestamp.Time), nil
	case *corev1.PersistentVolumeClaim:
		return time.Now().Sub(val.CreationTimestamp.Time), nil
	case corev1.PersistentVolumeClaim:
		return time.Now().Sub(val.CreationTimestamp.Time), nil
	default:
		return time.Duration(0), ErrInvalidToolsType
	}
}

// GetVolumeMode get volume mode of the persistentvolumeclaim.
// volumeMode defines what type of volume is required by the claim.
// Value of Filesystem is implied when not included in claim spec.
func (h *Handler) GetVolumeMode(object interface{}) (string, error) {
	switch val := object.(type) {
	case string:
		pvc, err := h.Get(val)
		if err != nil {
			return "", err
		}
		return string(*pvc.Spec.VolumeMode), nil
	case *corev1.PersistentVolumeClaim:
		return string(*val.Spec.VolumeMode), nil
	case corev1.PersistentVolumeClaim:
		return string(*val.Spec.VolumeMode), nil
	default:
		return "", ErrInvalidToolsType
	}
}
