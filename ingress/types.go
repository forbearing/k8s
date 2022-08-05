package ingress

import "fmt"

var (
	ERR_TYPE_TOOLS  = fmt.Errorf("type must be string *networkingv1.Ingress, or networkingv1.Ingress")
	ERR_TYPE_CREATE = fmt.Errorf("type must be string, []byte, *networkingv1.Ingress, networkingv1.Ingress, runtime.Object, *unstructured.Unstructured, unstructured.Unstructured or map[string]interface{}")
	ERR_TYPE_UPDATE = ERR_TYPE_CREATE
	ERR_TYPE_APPLY  = ERR_TYPE_CREATE
	ERR_TYPE_DELETE = ERR_TYPE_CREATE
	ErrInvalidGetType    = ERR_TYPE_CREATE
)
