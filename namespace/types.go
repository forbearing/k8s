package namespace

import "fmt"

var (
	ERR_TYPE_TOOLS  = fmt.Errorf("type must be string *corev1.Namespace, or corev1.Namespace")
	ERR_TYPE_CREATE = fmt.Errorf("type must be string, []byte, *corev1.Namespace, corev1.Namespace, runtime.Object, *unstructured.Unstructured, unstructured.Unstructured or map[string]interface{}")
	ErrInvalidUpdateType = ERR_TYPE_CREATE
	ErrInvalidApplyType  = ERR_TYPE_CREATE
	ErrInvalidDeleteType = ERR_TYPE_CREATE
	ErrInvalidGetType    = ERR_TYPE_CREATE
)
