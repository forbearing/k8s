package daemonset

import "errors"

var (
	ErrInvalidToolsType  = errors.New("type must be string, *appsv1.DaemonSet, appsv1.DaemonSet or runtime.Object")
	ErrInvalidCreateType = errors.New("type must be string, []byte, *appsv1.DaemonSet, appsv1.DaemonSet, runtime.Object, *unstructured.Unstructured, unstructured.Unstructured or map[string]interface{}")
	ErrInvalidUpdateType = ErrInvalidCreateType
	ErrInvalidApplyType  = ErrInvalidCreateType
	ErrInvalidDeleteType = ErrInvalidCreateType
	ErrInvalidGetType    = ErrInvalidCreateType
	ErrInvalidPatchType  = errors.New("patch data type must be string, []byte, *appsv1.DaemonSet, appsv1.DaemonSet, runtime.Object, *unstructured.Unstructured, unstructured.Unstructured or map[string]interface{}")
)
