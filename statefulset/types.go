package statefulset

import "fmt"

var (
	ERR_TYPE_TOOLS  = fmt.Errorf("type must be string *appsv1.StatefulSet, or appsv1.StatefulSet")
	ERR_TYPE_CREATE = fmt.Errorf("type must be string, []byte, *appsv1.StatefulSet, appsv1.StatefulSet, runtime.Object, *unstructured.Unstructured, unstructured.Unstructured or map[string]interface{}")
	ERR_TYPE_UPDATE = ERR_TYPE_CREATE
	ERR_TYPE_APPLY  = ERR_TYPE_CREATE
	ERR_TYPE_DELETE = ERR_TYPE_CREATE
	ErrInvalidGetType    = ERR_TYPE_CREATE
)
