package serviceaccount

import "errors"

var (
	ErrInvalidToolsType  = errors.New("type must be string, *corev1.ServiceAccount, corev1.ServiceAccount or runtime.Object")
	ErrInvalidCreateType = errors.New("type must be string, []byte, *corev1.ServiceAccount, corev1.ServiceAccount, runtime.Object, *unstructured.Unstructured, unstructured.Unstructured or map[string]interface{}")
	ErrInvalidUpdateType = ErrInvalidCreateType
	ErrInvalidApplyType  = ErrInvalidCreateType
	ErrInvalidDeleteType = ErrInvalidCreateType
	ErrInvalidGetType    = ErrInvalidCreateType
	ErrInvalidPathType   = errors.New("path data type must be string, []byte, *corev1.ServiceAccount, corev1.ServiceAccount, runtime.Object, *unstructured.Unstructured, unstructured.Unstructured or map[string]interface{}")
)
