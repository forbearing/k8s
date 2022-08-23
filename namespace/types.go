package namespace

import "fmt"

var (
	ErrInvalidToolsType  = fmt.Errorf("type must be string, *corev1.Namespace, or corev1.Namespace")
	ErrInvalidCreateType = fmt.Errorf("type must be string, []byte, *corev1.Namespace, corev1.Namespace, runtime.Object, *unstructured.Unstructured, unstructured.Unstructured or map[string]interface{}")
	ErrInvalidUpdateType = ErrInvalidCreateType
	ErrInvalidApplyType  = ErrInvalidCreateType
	ErrInvalidDeleteType = ErrInvalidCreateType
	ErrInvalidGetType    = ErrInvalidCreateType
)
