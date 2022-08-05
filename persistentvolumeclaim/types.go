package persistentvolumeclaim

import "fmt"

var (
	ERR_TYPE_TOOLS  = fmt.Errorf("type must be string *corev1.PersistentVolumeClaim, or corev1.PersistentVolumeClaim")
	ErrInvalidCreateType = fmt.Errorf("type must be string, []byte, *corev1.PersistentVolumeClaim, corev1.PersistentVolumeClaim, runtime.Object, *unstructured.Unstructured, unstructured.Unstructured or map[string]interface{}")
	ErrInvalidUpdateType = ErrInvalidCreateType
	ErrInvalidApplyType  = ErrInvalidCreateType
	ErrInvalidDeleteType = ErrInvalidCreateType
	ErrInvalidGetType    = ErrInvalidCreateType
)
