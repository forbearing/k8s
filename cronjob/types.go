package cronjob

import "errors"

var (
	ErrInvalidToolsType  = errors.New("type must be string, *batchv1.CronJob, batchv1.CronJob or runtime.Object")
	ErrInvalidCreateType = errors.New("type must be string, []byte, *batchv1.CronJob, batchv1.CronJob, runtime.Object, *unstructured.Unstructured, unstructured.Unstructured or map[string]interface{}")
	ErrInvalidUpdateType = ErrInvalidCreateType
	ErrInvalidApplyType  = ErrInvalidCreateType
	ErrInvalidDeleteType = ErrInvalidCreateType
	ErrInvalidGetType    = ErrInvalidCreateType
)
