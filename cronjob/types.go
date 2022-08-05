package cronjob

import "fmt"

var (
	ERR_TYPE_TOOLS  = fmt.Errorf("type must be string *batchv1.CronJob, or batchv1.CronJob")
	ErrInvalidCreateType = fmt.Errorf("type must be string, []byte, *batchv1.CronJob, batchv1.CronJob, runtime.Object, *unstructured.Unstructured, unstructured.Unstructured or map[string]interface{}")
	ErrInvalidUpdateType = ErrInvalidCreateType
	ErrInvalidApplyType  = ErrInvalidCreateType
	ErrInvalidDeleteType = ErrInvalidCreateType
	ErrInvalidGetType    = ErrInvalidCreateType
)
