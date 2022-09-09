package storageclass

import "errors"

var (
	ErrInvalidToolsType  = errors.New("type must be string, *storagev1.StorageClass, storagev1.StorageClass or runtime.Object")
	ErrInvalidCreateType = errors.New("type must be string, []byte, *storagev1.StorageClass, storagev1.StorageClass, runtime.Object, *unstructured.Unstructured, unstructured.Unstructured or map[string]interface{}")
	ErrInvalidUpdateType = ErrInvalidCreateType
	ErrInvalidApplyType  = ErrInvalidCreateType
	ErrInvalidDeleteType = ErrInvalidCreateType
	ErrInvalidGetType    = ErrInvalidCreateType
	ErrInvalidPathType   = errors.New("path data type must be string, []byte, *storagev1.StorageClass, storagev1.StorageClass, runtime.Object, *unstructured.Unstructured, unstructured.Unstructured or map[string]interface{}")
)
