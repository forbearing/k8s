package cronjob

import "errors"

var (
	ErrInvalidToolsType  = errors.New("type must be string, *batchv1.CronJob, batchv1.CronJob, metav1.Object or runtime.Object")
	ErrInvalidCreateType = errors.New("type must be string, []byte, *batchv1.CronJob, batchv1.CronJob, metav1.Object, runtime.Object, *unstructured.Unstructured, unstructured.Unstructured or map[string]interface{}")
	ErrInvalidUpdateType = ErrInvalidCreateType
	ErrInvalidApplyType  = ErrInvalidCreateType
	ErrInvalidDeleteType = ErrInvalidCreateType
	ErrInvalidGetType    = ErrInvalidCreateType
	ErrInvalidPatchType  = errors.New("patch data type must be string, []byte, *batchv1.CronJob, batchv1.CronJob, metav1.Object, runtime.Object, *unstructured.Unstructured, unstructured.Unstructured or map[string]interface{}")
)
