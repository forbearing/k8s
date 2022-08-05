package rolebinding

import "fmt"

var (
	ERR_TYPE_TOOLS  = fmt.Errorf("type must be string *rbacv1.RoleBinding, or rbacv1.RoleBinding")
	ERR_TYPE_CREATE = fmt.Errorf("type must be string, []byte, *rbacv1.RoleBinding, rbacv1.RoleBinding, runtime.Object, *unstructured.Unstructured, unstructured.Unstructured or map[string]interface{}")
	ERR_TYPE_UPDATE = ERR_TYPE_CREATE
	ERR_TYPE_APPLY  = ERR_TYPE_CREATE
	ERR_TYPE_DELETE = ERR_TYPE_CREATE
	ErrInvalidGetType    = ERR_TYPE_CREATE
)
