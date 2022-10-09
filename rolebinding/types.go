package rolebinding

import "errors"

var (
	ErrInvalidToolsType  = errors.New("type must be string, *rbacv1.RoleBinding, rbacv1.RoleBinding, metav1.Object or runtime.Object")
	ErrInvalidCreateType = errors.New("type must be string, []byte, *rbacv1.RoleBinding, rbacv1.RoleBinding, metav1.Object, runtime.Object, *unstructured.Unstructured, unstructured.Unstructured or map[string]interface{}")
	ErrInvalidUpdateType = ErrInvalidCreateType
	ErrInvalidApplyType  = ErrInvalidCreateType
	ErrInvalidDeleteType = ErrInvalidCreateType
	ErrInvalidGetType    = ErrInvalidCreateType
	ErrInvalidPatchType  = errors.New("patch data type must be string, []byte, *rbacv1.RoleBinding, rbacv1.RoleBinding, metav1.Object, runtime.Object, *unstructured.Unstructured, unstructured.Unstructured or map[string]interface{}")
)
