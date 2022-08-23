package clusterrolebinding

import "fmt"

var (
	ErrInvalidToolsType  = fmt.Errorf("type must be string, *rbacv1.ClusterRoleBinding, or rbacv1.ClusterRoleBinding")
	ErrInvalidCreateType = fmt.Errorf("type must be string, []byte, *rbacv1.ClusterRoleBinding, rbacv1.ClusterRoleBinding, runtime.Object, *unstructured.Unstructured, unstructured.Unstructured or map[string]interface{}")
	ErrInvalidUpdateType = ErrInvalidCreateType
	ErrInvalidApplyType  = ErrInvalidCreateType
	ErrInvalidDeleteType = ErrInvalidCreateType
	ErrInvalidGetType    = ErrInvalidCreateType
)
