package secret

import "fmt"

var (
	ERR_TYPE_TOOLS  = fmt.Errorf("type must be string *corev1.Secret, or corev1.Secret")
	ERR_TYPE_CREATE = fmt.Errorf("type must be string, []byte, *corev1.Secret, corev1.Secret, runtime.Object, *unstructured.Unstructured, unstructured.Unstructured or map[string]interface{}")
	ErrInvalidUpdateType = ERR_TYPE_CREATE
	ErrInvalidApplyType  = ERR_TYPE_CREATE
	ErrInvalidDeleteType = ERR_TYPE_CREATE
	ErrInvalidGetType    = ERR_TYPE_CREATE
)
