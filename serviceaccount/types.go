package serviceaccount

import "fmt"

var (
	ERR_TYPE_TOOLS  = fmt.Errorf("type must be string *corev1.ServiceAccount, or corev1.ServiceAccount")
	ErrInvalidCreateType = fmt.Errorf("type must be string, []byte, *corev1.ServiceAccount, corev1.ServiceAccount, runtime.Object, *unstructured.Unstructured, unstructured.Unstructured or map[string]interface{}")
	ErrInvalidUpdateType = ErrInvalidCreateType
	ErrInvalidApplyType  = ErrInvalidCreateType
	ErrInvalidDeleteType = ErrInvalidCreateType
	ErrInvalidGetType    = ErrInvalidCreateType
)
