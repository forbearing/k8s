package dynamic

import "fmt"

var (
	ErrGroupEmpty    = fmt.Errorf("group must not be empty")
	ErrVersionEmpty  = fmt.Errorf("version must not be empty")
	ErrResourceEmpty = fmt.Errorf("resource must not be empty")
	ErrInvalidType   = fmt.Errorf("type must be string, []byte, runtime.Object, *unstructured.Unstructured, unstructured.Unstructured or map[string]interface{}")
)
