package networkpolicy

import "fmt"

var (
	ERR_TYPE_TOOLS  = fmt.Errorf("type must be string *networkingv1.NetworkPolicy, or networkingv1.NetworkPolicy")
	ERR_TYPE_CREATE = fmt.Errorf("type must be string, []byte, *networkingv1.NetworkPolicy, networkingv1.NetworkPolicy, runtime.Object, *unstructured.Unstructured, unstructured.Unstructured or map[string]interface{}")
	ErrInvalidUpdateType = ERR_TYPE_CREATE
	ErrInvalidApplyType  = ERR_TYPE_CREATE
	ErrInvalidDeleteType = ERR_TYPE_CREATE
	ErrInvalidGetType    = ERR_TYPE_CREATE
)
