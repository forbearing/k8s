package node

import "fmt"

var (
	ErrInvalidToolsType  = fmt.Errorf("type must be string, *corev1.Node, or corev1.Node")
	ErrInvalidCreateType = fmt.Errorf("type must be string, []byte, *corev1.Node, corev1.Node, runtime.Object, *unstructured.Unstructured, unstructured.Unstructured or map[string]interface{}")
	ErrInvalidUpdateType = ErrInvalidCreateType
	ErrInvalidApplyType  = ErrInvalidCreateType
	ErrInvalidDeleteType = ErrInvalidCreateType
	ErrInvalidGetType    = ErrInvalidCreateType
)
