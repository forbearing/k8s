package networkpolicy

import "fmt"

var (
	ERR_TYPE_TOOLS  = fmt.Errorf("type must be string *networkingv1.NetworkPolicy, or networkingv1.NetworkPolicy")
	ErrInvalidCreateType = fmt.Errorf("type must be string, []byte, *networkingv1.NetworkPolicy, networkingv1.NetworkPolicy, runtime.Object, *unstructured.Unstructured, unstructured.Unstructured or map[string]interface{}")
	ErrInvalidUpdateType = ErrInvalidCreateType
	ErrInvalidApplyType  = ErrInvalidCreateType
	ErrInvalidDeleteType = ErrInvalidCreateType
	ErrInvalidGetType    = ErrInvalidCreateType
)
