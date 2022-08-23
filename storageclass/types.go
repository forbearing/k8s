package storageclass

import "fmt"

var (
	ErrInvalidToolsType  = fmt.Errorf("type must be string, *storagev1.StorageClass, or storagev1.StorageClass")
	ErrInvalidCreateType = fmt.Errorf("type must be string, []byte, *storagev1.StorageClass, storagev1.StorageClass, runtime.Object, *unstructured.Unstructured, unstructured.Unstructured or map[string]interface{}")
	ErrInvalidUpdateType = ErrInvalidCreateType
	ErrInvalidApplyType  = ErrInvalidCreateType
	ErrInvalidDeleteType = ErrInvalidCreateType
	ErrInvalidGetType    = ErrInvalidCreateType
)
