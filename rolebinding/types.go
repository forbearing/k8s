package rolebinding

import "fmt"

var (
	ErrInvalidToolsType  = fmt.Errorf("type must be string, *rbacv1.RoleBinding, or rbacv1.RoleBinding")
	ErrInvalidCreateType = fmt.Errorf("type must be string, []byte, *rbacv1.RoleBinding, rbacv1.RoleBinding, runtime.Object, *unstructured.Unstructured, unstructured.Unstructured or map[string]interface{}")
	ErrInvalidUpdateType = ErrInvalidCreateType
	ErrInvalidApplyType  = ErrInvalidCreateType
	ErrInvalidDeleteType = ErrInvalidCreateType
	ErrInvalidGetType    = ErrInvalidCreateType
)
