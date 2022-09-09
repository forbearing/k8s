package deployment

import (
	"errors"
)

var (
	ErrInvalidToolsType  = errors.New("type must be string, *appsv1.Deployment, appsv1.Deployment or runtime.Object")
	ErrInvalidCreateType = errors.New("type must be string, []byte, *appsv1.Deployment, appsv1.Deployment, runtime.Object, *unstructured.Unstructured, unstructured.Unstructured or map[string]interface{}")
	ErrInvalidUpdateType = ErrInvalidCreateType
	ErrInvalidApplyType  = ErrInvalidCreateType
	ErrInvalidDeleteType = ErrInvalidCreateType
	ErrInvalidGetType    = ErrInvalidCreateType
	ErrInvalidScaleType  = ErrInvalidCreateType
	ErrInvalidPathType   = errors.New("path data type must be string, []byte, *appsv1.Deployment, appsv1.Deployment, runtime.Object, *unstructured.Unstructured, unstructured.Unstructured or map[string]interface{}")
)
