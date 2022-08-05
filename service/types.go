package service

import "fmt"

var (
	ERR_TYPE_TOOLS  = fmt.Errorf("type must be string *corev1.Service, or corev1.Service")
	ERR_TYPE_CREATE = fmt.Errorf("type must be string, []byte, *corev1.Service, corev1.Service, runtime.Object, *unstructured.Unstructured, unstructured.Unstructured or map[string]interface{}")
	ErrInvalidUpdateType = ERR_TYPE_CREATE
	ErrInvalidApplyType  = ERR_TYPE_CREATE
	ErrInvalidDeleteType = ERR_TYPE_CREATE
	ErrInvalidGetType    = ERR_TYPE_CREATE
)
