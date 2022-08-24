package persistentvolume

import (
	"reflect"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
)

// GetCapacity get the the storage capacity of the persistentvolume.
func (h *Handler) GetCapacity(object interface{}) (int64, error) {
	switch val := object.(type) {
	case string:
		pv, err := h.Get(val)
		if err != nil {
			return 0, err
		}
		return h.getCapacity(pv), nil
	case *corev1.PersistentVolume:
		return h.getCapacity(val), nil
	case corev1.PersistentVolume:
		return h.getCapacity(&val), nil
	default:
		return 0, ErrInvalidToolsType
	}
}
func (h *Handler) getCapacity(pv *corev1.PersistentVolume) int64 {
	storage := pv.Spec.Capacity[corev1.ResourceName(corev1.ResourceStorage)]
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

// GetAccessModes get the accessModes of the persistentvolume.
func (h *Handler) GetAccessModes(object interface{}) (accessModes []string, err error) {
	switch val := object.(type) {
	case string:
		pv, err := h.Get(val)
		if err != nil {
			return nil, err
		}
		return h.getAccessModes(pv), nil
	case *corev1.PersistentVolume:
		return h.getAccessModes(val), nil
	case corev1.PersistentVolume:
		return h.getAccessModes(&val), nil
	default:
		return nil, ErrInvalidToolsType
	}
}
func (h *Handler) getAccessModes(pv *corev1.PersistentVolume) []string {
	var accessModes []string
	for _, accessMode := range pv.Spec.AccessModes {
		accessModes = append(accessModes, string(accessMode))
	}
	return accessModes
}

// GetReclaimPolicy get the reclaim policy of the persistentvolume.
func (h *Handler) GetReclaimPolicy(object interface{}) (string, error) {
	switch val := object.(type) {
	case string:
		pv, err := h.Get(val)
		if err != nil {
			return "", err
		}
		return string(pv.Spec.PersistentVolumeReclaimPolicy), nil
	case *corev1.PersistentVolume:
		return string(val.Spec.PersistentVolumeReclaimPolicy), nil
	case corev1.PersistentVolume:
		return string(val.Spec.PersistentVolumeReclaimPolicy), nil
	default:
		return "", ErrInvalidToolsType
	}
}

// GetStatus get the status phase of the persistentvolume.
// All supported pv status are: Pending, Available, Bound, Released, Failed.
// Pending used for PersistentVolumes that are not available.
// Available used for PersistentVolumes that are not yet bound.
// Bound used for PersistentVolumes that are bound.
// Released used for PersistentVolumes where the bound PersistentVolumeClaim was deleted.
// Failed used for PersistentVolumes that failed to be correctly recycled or
// deleted after being released from a claim.
func (h *Handler) GetStatus(object interface{}) (string, error) {
	switch val := object.(type) {
	case string:
		pv, err := h.Get(val)
		if err != nil {
			return "", err
		}
		return string(pv.Status.Phase), nil
	case *corev1.PersistentVolume:
		return string(val.Status.Phase), nil
	case corev1.PersistentVolume:
		return string(val.Status.Phase), nil
	default:
		return "", ErrInvalidToolsType
	}
}

// GetClaim simply calls GetPVC.
func (h *Handler) GetClaim(object interface{}) (string, error) {
	return h.GetPVC(object)
}

// GetPVC get the pvc name of the persistentvolume.
func (h *Handler) GetPVC(object interface{}) (string, error) {
	switch val := object.(type) {
	case string:
		pv, err := h.Get(val)
		if err != nil {
			return "", err
		}
		return h.getPVC(pv), nil
	case *corev1.PersistentVolume:
		return h.getPVC(val), nil
	case corev1.PersistentVolume:
		return h.getPVC(&val), nil
	default:
		return "", ErrInvalidToolsType
	}
}
func (h *Handler) getPVC(pv *corev1.PersistentVolume) string {
	var pvc string
	if pv.Spec.ClaimRef != nil {
		if pv.Spec.ClaimRef.Kind == "PersistentVolumeClaim" {
			pvc = pv.Spec.ClaimRef.Name
		}
	}
	return pvc
}

// GetStorageClass get the storageclass name of the persistentvolume.
func (h *Handler) GetStorageClass(object interface{}) (string, error) {
	switch val := object.(type) {
	case string:
		pv, err := h.Get(val)
		if err != nil {
			return "", err
		}
		return pv.Spec.StorageClassName, nil
	case *corev1.PersistentVolume:
		return val.Spec.StorageClassName, nil
	case corev1.PersistentVolume:
		return val.Spec.StorageClassName, nil
	default:
		return "", ErrInvalidToolsType
	}
}

func (h *Handler) GetVolumeSource(object interface{}) (string, error) {
	switch val := object.(type) {
	case string:
		pv, err := h.Get(val)
		if err != nil {
			return "", err
		}
		return h.getVolumeSource(pv), nil
	case *corev1.PersistentVolume:
		return h.getVolumeSource(val), nil
	case corev1.PersistentVolume:
		return h.getVolumeSource(&val), nil
	default:
		return "", ErrInvalidToolsType
	}

}
func (h *Handler) getVolumeSource(pv *corev1.PersistentVolume) string {
	// 通过反射来做
	t := reflect.TypeOf(pv.Spec.PersistentVolumeSource)
	v := reflect.ValueOf(pv.Spec.PersistentVolumeSource)

	for i := 0; i < v.NumField(); i++ {
		//log.Println(t.Field(i).Tag.Get("protobuf"))
		//log.Println(t.Field(i).Tag.Get("json"))
		val := v.Field(i).Interface()
		if !reflect.ValueOf(val).IsNil() {
			tag := t.Field(i).Tag.Get("json") // nfs,omitempty
			source := strings.Split(tag, ",") // [nfs]
			return source[0]
		}

	}
	return ""
}

// GetAge returns age of the persistentvolume.
func (h *Handler) GetAge(object interface{}) (time.Duration, error) {
	switch val := object.(type) {
	case string:
		pv, err := h.Get(val)
		if err != nil {
			return time.Duration(0), err
		}
		return time.Now().Sub(pv.CreationTimestamp.Time), nil
	case *corev1.PersistentVolume:
		return time.Now().Sub(val.CreationTimestamp.Time), nil
	case corev1.PersistentVolume:
		return time.Now().Sub(val.CreationTimestamp.Time), nil
	default:
		return time.Duration(0), ErrInvalidToolsType
	}
}

// GetVolumeMode get volume mode of the persistentvolume.
// volumeMode defines what type of volume is required by the claim.
// Value of Filesystem is implied when not included in claim spec.
func (h *Handler) GetVolumeMode(object interface{}) (string, error) {
	switch val := object.(type) {
	case string:
		pv, err := h.Get(val)
		if err != nil {
			return "", err
		}
		return string(*pv.Spec.VolumeMode), nil
	case *corev1.PersistentVolume:
		return string(*val.Spec.VolumeMode), nil
	case corev1.PersistentVolume:
		return string(*val.Spec.VolumeMode), nil
	default:
		return "", ErrInvalidToolsType
	}
}
