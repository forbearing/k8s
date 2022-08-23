package deployment

import "fmt"

var (
	ErrInvalidToolsType  = fmt.Errorf("type must be string *appsv1.Deployment, or appsv1.Deployment")
	ErrInvalidCreateType = fmt.Errorf("type must be string, []byte, *appsv1.Deployment, appsv1.Deployment, runtime.Object, *unstructured.Unstructured, unstructured.Unstructured or map[string]interface{}")
	ErrInvalidUpdateType = ErrInvalidCreateType
	ErrInvalidApplyType  = ErrInvalidCreateType
	ErrInvalidDeleteType = ErrInvalidCreateType
	ErrInvalidGetType    = ErrInvalidCreateType
	ErrInvalidScaleType  = ErrInvalidCreateType
)
