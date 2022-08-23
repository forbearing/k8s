package role

import "fmt"

var (
	ErrInvalidToolsType  = fmt.Errorf("type must be string, *rbacv1.Role, or rbacv1.Role")
	ErrInvalidCreateType = fmt.Errorf("type must be string, []byte, *rbacv1.Role, rbacv1.Role, runtime.Object, *unstructured.Unstructured, unstructured.Unstructured or map[string]interface{}")
	ErrInvalidUpdateType = ErrInvalidCreateType
	ErrInvalidApplyType  = ErrInvalidCreateType
	ErrInvalidDeleteType = ErrInvalidCreateType
	ErrInvalidGetType    = ErrInvalidCreateType
)
