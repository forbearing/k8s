package job

import "fmt"

var (
	ErrInvalidToolsType  = fmt.Errorf("type must be string, *batchv1.Job, or batchv1.Job")
	ErrInvalidCreateType = fmt.Errorf("type must be string, []byte, *batchv1.Job, batchv1.Job, runtime.Object, *unstructured.Unstructured, unstructured.Unstructured or map[string]interface{}")
	ErrInvalidUpdateType = ErrInvalidCreateType
	ErrInvalidApplyType  = ErrInvalidCreateType
	ErrInvalidDeleteType = ErrInvalidCreateType
	ErrInvalidGetType    = ErrInvalidCreateType
)
