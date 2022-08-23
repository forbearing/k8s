package persistentvolume

import "fmt"

var (
	ErrInvalidToolsType  = fmt.Errorf("type must be string, *corev1.PersistentVolume, or corev1.PersistentVolume")
	ErrInvalidCreateType = fmt.Errorf("type must be string, []byte, *corev1.PersistentVolume, corev1.PersistentVolume, runtime.Object, *unstructured.Unstructured, unstructured.Unstructured or map[string]interface{}")
	ErrInvalidUpdateType = ErrInvalidCreateType
	ErrInvalidApplyType  = ErrInvalidCreateType
	ErrInvalidDeleteType = ErrInvalidCreateType
	ErrInvalidGetType    = ErrInvalidCreateType
)
