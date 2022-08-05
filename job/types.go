package job

import "fmt"

var (
	ERR_TYPE_TOOLS  = fmt.Errorf("type must be string *batchv1.Job, or batchv1.Job")
	ERR_TYPE_CREATE = fmt.Errorf("type must be string, []byte, *batchv1.Job, batchv1.Job, runtime.Object, *unstructured.Unstructured, unstructured.Unstructured or map[string]interface{}")
	ERR_TYPE_UPDATE = ERR_TYPE_CREATE
	ERR_TYPE_APPLY  = ERR_TYPE_CREATE
	ERR_TYPE_DELETE = ERR_TYPE_CREATE
	ErrInvalidGetType    = ERR_TYPE_CREATE
)
