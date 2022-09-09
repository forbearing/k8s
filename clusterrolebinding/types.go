package clusterrolebinding

import "errors"

var (
	ErrInvalidToolsType  = errors.New("type must be string, *rbacv1.ClusterRoleBinding, rbacv1.ClusterRoleBinding or runtime.Object")
	ErrInvalidCreateType = errors.New("type must be string, []byte, *rbacv1.ClusterRoleBinding, rbacv1.ClusterRoleBinding, runtime.Object, *unstructured.Unstructured, unstructured.Unstructured or map[string]interface{}")
	ErrInvalidUpdateType = ErrInvalidCreateType
	ErrInvalidApplyType  = ErrInvalidCreateType
	ErrInvalidDeleteType = ErrInvalidCreateType
	ErrInvalidGetType    = ErrInvalidCreateType
	ErrInvalidPathType   = errors.New("path data type must be string, []byte, *rbacv1.ClusterRoleBinding, rbacv1.ClusterRoleBinding, runtime.Object, *unstructured.Unstructured, unstructured.Unstructured or map[string]interface{}")
)
