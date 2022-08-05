package statefulset

import "fmt"

var (
	ERR_TYPE_TOOLS  = fmt.Errorf("type must be string *appsv1.StatefulSet, or appsv1.StatefulSet")
	ErrInvalidCreateType = fmt.Errorf("type must be string, []byte, *appsv1.StatefulSet, appsv1.StatefulSet, runtime.Object, *unstructured.Unstructured, unstructured.Unstructured or map[string]interface{}")
	ErrInvalidUpdateType = ErrInvalidCreateType
	ErrInvalidApplyType  = ErrInvalidCreateType
	ErrInvalidDeleteType = ErrInvalidCreateType
	ErrInvalidGetType    = ErrInvalidCreateType
)
