package cronjob

import "fmt"

var (
	ERR_TYPE_TOOLS  = fmt.Errorf("type must be string *batchv1.CronJob, or batchv1.CronJob")
	ERR_TYPE_CREATE = fmt.Errorf("type must be string, []byte, *batchv1.CronJob, batchv1.CronJob, runtime.Object, *unstructured.Unstructured, unstructured.Unstructured or map[string]interface{}")
	ERR_TYPE_UPDATE = ERR_TYPE_CREATE
	ERR_TYPE_APPLY  = ERR_TYPE_CREATE
	ERR_TYPE_DELETE = ERR_TYPE_CREATE
	ErrInvalidGetType    = ERR_TYPE_CREATE
)
