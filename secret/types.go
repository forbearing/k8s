package secret

import "fmt"

var (
	ErrInvalidToolsType  = fmt.Errorf("type must be string, *corev1.Secret, or corev1.Secret")
	ErrInvalidCreateType = fmt.Errorf("type must be string, []byte, *corev1.Secret, corev1.Secret, runtime.Object, *unstructured.Unstructured, unstructured.Unstructured or map[string]interface{}")
	ErrInvalidUpdateType = ErrInvalidCreateType
	ErrInvalidApplyType  = ErrInvalidCreateType
	ErrInvalidDeleteType = ErrInvalidCreateType
	ErrInvalidGetType    = ErrInvalidCreateType
)
