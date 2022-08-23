package ingressclass

import "fmt"

var (
	ErrInvalidToolsType  = fmt.Errorf("type must be string, *networkingv1.IngressClass, or networkingv1.IngressClass")
	ErrInvalidCreateType = fmt.Errorf("type must be string, []byte, *networkingv1.IngressClass, networkingv1.IngressClass, runtime.Object, *unstructured.Unstructured, unstructured.Unstructured or map[string]interface{}")
	ErrInvalidUpdateType = ErrInvalidCreateType
	ErrInvalidApplyType  = ErrInvalidCreateType
	ErrInvalidDeleteType = ErrInvalidCreateType
	ErrInvalidGetType    = ErrInvalidCreateType
)
