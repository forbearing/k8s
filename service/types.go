package service

import "fmt"

var (
	ErrInvalidToolsType  = fmt.Errorf("type must be string, *corev1.Service, or corev1.Service")
	ErrInvalidCreateType = fmt.Errorf("type must be string, []byte, *corev1.Service, corev1.Service, runtime.Object, *unstructured.Unstructured, unstructured.Unstructured or map[string]interface{}")
	ErrInvalidUpdateType = ErrInvalidCreateType
	ErrInvalidApplyType  = ErrInvalidCreateType
	ErrInvalidDeleteType = ErrInvalidCreateType
	ErrInvalidGetType    = ErrInvalidCreateType
)
