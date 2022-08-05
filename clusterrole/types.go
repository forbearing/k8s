package clusterrole

import "fmt"

var (
	ERR_TYPE_TOOLS  = fmt.Errorf("type must be string *rbacv1.ClusterRole, or rbacv1.ClusterRole")
	ErrInvalidCreateType = fmt.Errorf("type must be string, []byte, *rbacv1.ClusterRole, rbacv1.ClusterRole, runtime.Object, *unstructured.Unstructured, unstructured.Unstructured or map[string]interface{}")
	ErrInvalidUpdateType = ErrInvalidCreateType
	ErrInvalidApplyType  = ErrInvalidCreateType
	ErrInvalidDeleteType = ErrInvalidCreateType
	ErrInvalidGetType    = ErrInvalidCreateType
)
