package dynamic

import "errors"

var (
	ErrGroupEmpty        = errors.New("group must not be empty")
	ErrVersionEmpty      = errors.New("version must not be empty")
	ErrResourceEmpty     = errors.New("resource must not be empty")
	ErrInvalidCreateType = errors.New("type must be string, []byte, metav1.Object, runtime.Object, *unstructured.Unstructured, unstructured.Unstructured or map[string]interface{}")
	ErrInvalidUpdateType = ErrInvalidCreateType
	ErrInvalidApplyType  = ErrInvalidCreateType
	ErrInvalidDeleteType = ErrInvalidCreateType
	ErrInvalidGetType    = ErrInvalidCreateType
	ErrInvalidPatchType  = errors.New("patch type must be string, []byte, metav1.Object, runtime.Object, *unstructured.Unstructured, unstructured.Unstructured or map[string]interface{}")
)
