package replicaset

import "fmt"

var (
	ERR_TYPE_TOOLS  = fmt.Errorf("type must be string *appsv1.ReplicaSet, or appsv1.ReplicaSet")
	ERR_TYPE_CREATE = fmt.Errorf("type must be string, []byte, *appsv1.ReplicaSet, appsv1.ReplicaSet, runtime.Object, *unstructured.Unstructured, unstructured.Unstructured or map[string]interface{}")
	ERR_TYPE_UPDATE = ERR_TYPE_CREATE
	ERR_TYPE_APPLY  = ERR_TYPE_CREATE
	ERR_TYPE_DELETE = ERR_TYPE_CREATE
	ErrInvalidGetType    = ERR_TYPE_CREATE
)
