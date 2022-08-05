package clusterrolebinding

import "fmt"

var (
	ERR_TYPE_TOOLS  = fmt.Errorf("type must be string *rbacv1.ClusterRoleBinding, or rbacv1.ClusterRoleBinding")
	ERR_TYPE_CREATE = fmt.Errorf("type must be string, []byte, *rbacv1.ClusterRoleBinding, rbacv1.ClusterRoleBinding, runtime.Object, *unstructured.Unstructured, unstructured.Unstructured or map[string]interface{}")
	ErrInvalidUpdateType = ERR_TYPE_CREATE
	ErrInvalidApplyType  = ERR_TYPE_CREATE
	ErrInvalidDeleteType = ERR_TYPE_CREATE
	ErrInvalidGetType    = ERR_TYPE_CREATE
)
