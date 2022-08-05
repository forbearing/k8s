package clusterrole

import "fmt"

var (
	ERR_TYPE_TOOLS  = fmt.Errorf("type must be string *rbacv1.ClusterRole, or rbacv1.ClusterRole")
	ERR_TYPE_CREATE = fmt.Errorf("type must be string, []byte, *rbacv1.ClusterRole, rbacv1.ClusterRole, runtime.Object, *unstructured.Unstructured, unstructured.Unstructured or map[string]interface{}")
	ERR_TYPE_UPDATE = ERR_TYPE_CREATE
	ERR_TYPE_APPLY  = ERR_TYPE_CREATE
	ERR_TYPE_DELETE = ERR_TYPE_CREATE
	ErrInvalidGetType    = ERR_TYPE_CREATE
)
