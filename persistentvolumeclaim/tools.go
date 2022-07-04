package persistentvolumeclaim

import corev1 "k8s.io/api/core/v1"

// GetPV get the pv name of the persistentvolumeclaim
func (h *Handler) GetPV(name string) (pv string, err error) {
	pvc, err := h.Get(name)
	if err != nil {
		return
	}
	pv = pvc.Spec.VolumeName
	return
}

// GetStorageClass get the storageclass name of the persistentvolumeclaim.
func (h *Handler) GetStorageClass(name string) (sc string, err error) {
	pvc, err := h.Get(name)
	if err != nil {
		return
	}
	sc = *pvc.Spec.StorageClassName
	return
}

// GetAccessModes get the access modes of the persistentvolumeclaim.
func (h *Handler) GetAccessModes(name string) (accessModes []string, err error) {
	pvc, err := h.Get(name)
	if err != nil {
		return
	}
	for _, accessMode := range pvc.Status.AccessModes {
		accessModes = append(accessModes, string(accessMode))
	}
	return
}

// GetCapacity get the storage capacity of the persistentvolumeclaim.
func (h *Handler) GetCapacity(name string) (capacity int64, err error) {
	pvc, err := h.Get(name)
	if err != nil {
		return
	}
	storage := pvc.Status.Capacity[corev1.ResourceName(corev1.ResourceStorage)]
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

// GetPhase get the status phase of the persistentvolumeclaim.
func (h *Handler) GetPhase(name string) (phase string, err error) {
	pvc, err := h.Get(name)
	if err != nil {
		return
	}
	phase = string(pvc.Status.Phase)
	return
}
