package persistentvolumeclaim

import "errors"

var (
	ErrInvalidToolsType  = errors.New("type must be string, *corev1.PersistentVolumeClaim, corev1.PersistentVolumeClaim or runtime.Object")
	ErrInvalidCreateType = errors.New("type must be string, []byte, *corev1.PersistentVolumeClaim, corev1.PersistentVolumeClaim, runtime.Object, *unstructured.Unstructured, unstructured.Unstructured or map[string]interface{}")
	ErrInvalidUpdateType = ErrInvalidCreateType
	ErrInvalidApplyType  = ErrInvalidCreateType
	ErrInvalidDeleteType = ErrInvalidCreateType
	ErrInvalidGetType    = ErrInvalidCreateType
)
