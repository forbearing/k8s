package secret

import "errors"

var (
	ErrInvalidToolsType  = errors.New("type must be string, *corev1.Secret, corev1.Secret or runtime.Object")
	ErrInvalidCreateType = errors.New("type must be string, []byte, *corev1.Secret, corev1.Secret, runtime.Object, *unstructured.Unstructured, unstructured.Unstructured or map[string]interface{}")
	ErrInvalidUpdateType = ErrInvalidCreateType
	ErrInvalidApplyType  = ErrInvalidCreateType
	ErrInvalidDeleteType = ErrInvalidCreateType
	ErrInvalidGetType    = ErrInvalidCreateType
)
