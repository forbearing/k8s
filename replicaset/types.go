package replicaset

import "fmt"

var (
	ERR_TYPE_TOOLS  = fmt.Errorf("type must be string *appsv1.ReplicaSet, or appsv1.ReplicaSet")
	ErrInvalidCreateType = fmt.Errorf("type must be string, []byte, *appsv1.ReplicaSet, appsv1.ReplicaSet, runtime.Object, *unstructured.Unstructured, unstructured.Unstructured or map[string]interface{}")
	ErrInvalidUpdateType = ErrInvalidCreateType
	ErrInvalidApplyType  = ErrInvalidCreateType
	ErrInvalidDeleteType = ErrInvalidCreateType
	ErrInvalidGetType    = ErrInvalidCreateType
)
