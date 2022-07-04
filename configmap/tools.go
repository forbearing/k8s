package configmap

// GetData get configmap.spec.data
func (h *Handler) GetData(name string) (map[string]string, error) {
	data := make(map[string]string)
	configmap, err := h.Get(name)
	if err != nil {
		return data, err
	}
	data = configmap.Data
	return data, nil
}
