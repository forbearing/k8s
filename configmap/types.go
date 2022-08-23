package configmap

import "fmt"

var (
	ErrInvalidToolsType  = fmt.Errorf("type must be string, *corev1.ConfigMap, or corev1.ConfigMap")
	ErrInvalidCreateType = fmt.Errorf("type must be string, []byte, *corev1.ConfigMap, corev1.ConfigMap, runtime.Object, *unstructured.Unstructured, unstructured.Unstructured or map[string]interface{}")
	ErrInvalidUpdateType = ErrInvalidCreateType
	ErrInvalidApplyType  = ErrInvalidCreateType
	ErrInvalidDeleteType = ErrInvalidCreateType
	ErrInvalidGetType    = ErrInvalidCreateType
)
