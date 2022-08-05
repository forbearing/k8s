package deployment

import "fmt"

var (
	ERR_TYPE_TOOLS  = fmt.Errorf("type must be string *appsv1.Deployment, or appsv1.Deployment")
	ERR_TYPE_CREATE = fmt.Errorf("type must be string, []byte, *appsv1.Deployment, appsv1.Deployment, runtime.Object, *unstructured.Unstructured, unstructured.Unstructured or map[string]interface{}")
	ErrInvalidUpdateType = ERR_TYPE_CREATE
	ErrInvalidApplyType  = ERR_TYPE_CREATE
	ErrInvalidDeleteType = ERR_TYPE_CREATE
	ErrInvalidGetType    = ERR_TYPE_CREATE
)
