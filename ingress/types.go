package ingress

import "fmt"

var (
	ErrInvalidToolsType  = fmt.Errorf("type must be string, *networkingv1.Ingress, or networkingv1.Ingress")
	ErrInvalidCreateType = fmt.Errorf("type must be string, []byte, *networkingv1.Ingress, networkingv1.Ingress, runtime.Object, *unstructured.Unstructured, unstructured.Unstructured or map[string]interface{}")
	ErrInvalidUpdateType = ErrInvalidCreateType
	ErrInvalidApplyType  = ErrInvalidCreateType
	ErrInvalidDeleteType = ErrInvalidCreateType
	ErrInvalidGetType    = ErrInvalidCreateType
)
