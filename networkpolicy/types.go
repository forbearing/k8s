package networkpolicy

import "errors"

var (
	ErrInvalidToolsType  = errors.New("type must be string, *networkingv1.NetworkPolicy, networkingv1.NetworkPolicy or runtime.Object")
	ErrInvalidCreateType = errors.New("type must be string, []byte, *networkingv1.NetworkPolicy, networkingv1.NetworkPolicy, runtime.Object, *unstructured.Unstructured, unstructured.Unstructured or map[string]interface{}")
	ErrInvalidUpdateType = ErrInvalidCreateType
	ErrInvalidApplyType  = ErrInvalidCreateType
	ErrInvalidDeleteType = ErrInvalidCreateType
	ErrInvalidGetType    = ErrInvalidCreateType
)
