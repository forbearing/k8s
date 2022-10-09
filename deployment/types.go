package deployment

import (
	"errors"
)

var (
	ErrInvalidToolsType  = errors.New("type must be string, *appsv1.Deployment, appsv1.Deployment, metav1.Object or runtime.Object")
	ErrInvalidCreateType = errors.New("type must be string, []byte, *appsv1.Deployment, appsv1.Deployment, metav1.Object, runtime.Object, *unstructured.Unstructured, unstructured.Unstructured or map[string]interface{}")
	ErrInvalidUpdateType = ErrInvalidCreateType
	ErrInvalidApplyType  = ErrInvalidCreateType
	ErrInvalidDeleteType = ErrInvalidCreateType
	ErrInvalidGetType    = ErrInvalidCreateType
	ErrInvalidScaleType  = ErrInvalidCreateType
	ErrInvalidPatchType  = errors.New("patch data type must be string, []byte, *appsv1.Deployment, appsv1.Deployment, metav1.Object, runtime.Object, *unstructured.Unstructured, unstructured.Unstructured or map[string]interface{}")
)
