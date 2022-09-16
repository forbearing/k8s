package service

import "errors"

var (
	ErrInvalidToolsType  = errors.New("type must be string, *corev1.Service, corev1.Service or runtime.Object")
	ErrInvalidCreateType = errors.New("type must be string, []byte, *corev1.Service, corev1.Service, runtime.Object, *unstructured.Unstructured, unstructured.Unstructured or map[string]interface{}")
	ErrInvalidUpdateType = ErrInvalidCreateType
	ErrInvalidApplyType  = ErrInvalidCreateType
	ErrInvalidDeleteType = ErrInvalidCreateType
	ErrInvalidGetType    = ErrInvalidCreateType
	ErrInvalidPatchType  = errors.New("patch data type must be string, []byte, *corev1.Service, corev1.Service, runtime.Object, *unstructured.Unstructured, unstructured.Unstructured or map[string]interface{}")
)
