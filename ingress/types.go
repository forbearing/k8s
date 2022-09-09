package ingress

import "errors"

var (
	ErrInvalidToolsType  = errors.New("type must be string, *networkingv1.Ingress, networkingv1.Ingress or runtime.Object")
	ErrInvalidCreateType = errors.New("type must be string, []byte, *networkingv1.Ingress, networkingv1.Ingress, runtime.Object, *unstructured.Unstructured, unstructured.Unstructured or map[string]interface{}")
	ErrInvalidUpdateType = ErrInvalidCreateType
	ErrInvalidApplyType  = ErrInvalidCreateType
	ErrInvalidDeleteType = ErrInvalidCreateType
	ErrInvalidGetType    = ErrInvalidCreateType
	ErrInvalidPathType   = errors.New("path data type must be string, []byte, *networkingv1.Ingress, networkingv1.Ingress, runtime.Object, *unstructured.Unstructured, unstructured.Unstructured or map[string]interface{}")
)
