package dynamic

import "errors"

var (
	ErrGroupEmpty    = errors.New("group must not be empty")
	ErrVersionEmpty  = errors.New("version must not be empty")
	ErrResourceEmpty = errors.New("resource must not be empty")
	ErrInvalidType   = errors.New("type must be string, []byte, runtime.Object, *unstructured.Unstructured, unstructured.Unstructured or map[string]interface{}")
)
