package persistentvolume

import corev1 "k8s.io/api/core/v1"

// GetPVC get the pvc name of the persistentvolume.
func (h *Handler) GetPVC(name string) (pvc string, err error) {
	pv, err := h.Get(name)
	if err != nil {
		return
	}
	if pv.Spec.ClaimRef != nil {
		if pv.Spec.ClaimRef.Kind == "PersistentVolumeClaim" {
			pvc = pv.Spec.ClaimRef.Name
		}
	}
	return
}

// GetStorageClass get the storageclass name of the persistentvolume.
func (h *Handler) GetStorageClass(name string) (sc string, err error) {
	pv, err := h.Get(name)
	if err != nil {
		return
	}
	sc = pv.Spec.StorageClassName
	return
}

// GetAccessModes get the accessModes of the persistentvolume.
func (h *Handler) GetAccessModes(name string) (accessModes []string, err error) {
	pv, err := h.Get(name)
	if err != nil {
		return
	}
	for _, accessMode := range pv.Spec.AccessModes {
		accessModes = append(accessModes, string(accessMode))
	}
	return
}

// GetCapacity get the the storage capacity of the persistentvolume.
func (h *Handler) GetCapacity(name string) (capacity int64, err error) {
	pv, err := h.Get(name)
	if err != nil {
		return
	}
	storage := pv.Spec.Capacity[corev1.ResourceName(corev1.ResourceStorage)]
	//capacity = storage.Value()
	//capacity = storage.MilliValue()
	//capacity = storage.ScaledValue(resource.Kilo)
	//capacity = storage.ScaledValue(resource.Mega)
	//capacity = storage.ScaledValue(resource.Giga)
	//capacity = storage.ScaledValue(resource.Tera)
	//capacity = storage.ScaledValue(resource.Peta)
	//capacity = storage.ScaledValue(resource.Exa)
	capacity = storage.Value()
	return
}

// GetPhase get the status phase of the persistentvolume.
func (h *Handler) GetPhase(name string) (phase string, err error) {
	pv, err := h.Get(name)
	if err != nil {
		return
	}
	phase = string(pv.Status.Phase)
	return
}

// GetReclaimPolicy get the reclaim policy of the persistentvolume.
func (h *Handler) GetReclaimPolicy(name string) (policy string, err error) {
	pv, err := h.Get(name)
	if err != nil {
		return
	}
	policy = string(pv.Spec.PersistentVolumeReclaimPolicy)
	return
}
