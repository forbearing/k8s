package job

import "errors"

var (
	ErrInvalidToolsType  = errors.New("type must be string, *batchv1.Job, batchv1.Job or runtime.Object")
	ErrInvalidCreateType = errors.New("type must be string, []byte, *batchv1.Job, batchv1.Job, runtime.Object, *unstructured.Unstructured, unstructured.Unstructured or map[string]interface{}")
	ErrInvalidUpdateType = ErrInvalidCreateType
	ErrInvalidApplyType  = ErrInvalidCreateType
	ErrInvalidDeleteType = ErrInvalidCreateType
	ErrInvalidGetType    = ErrInvalidCreateType
)
