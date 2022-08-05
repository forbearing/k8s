package clusterrolebinding

import "fmt"

var (
	ERR_TYPE_TOOLS  = fmt.Errorf("type must be string *rbacv1.ClusterRoleBinding, or rbacv1.ClusterRoleBinding")
	ErrInvalidCreateType = fmt.Errorf("type must be string, []byte, *rbacv1.ClusterRoleBinding, rbacv1.ClusterRoleBinding, runtime.Object, *unstructured.Unstructured, unstructured.Unstructured or map[string]interface{}")
	ErrInvalidUpdateType = ErrInvalidCreateType
	ErrInvalidApplyType  = ErrInvalidCreateType
	ErrInvalidDeleteType = ErrInvalidCreateType
	ErrInvalidGetType    = ErrInvalidCreateType
)
