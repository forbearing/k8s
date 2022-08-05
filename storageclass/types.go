package storageclass

import "fmt"

var (
	ERR_TYPE_TOOLS  = fmt.Errorf("type must be string *storagev1.StorageClass, or storagev1.StorageClass")
	ERR_TYPE_CREATE = fmt.Errorf("type must be string, []byte, *storagev1.StorageClass, storagev1.StorageClass, runtime.Object, *unstructured.Unstructured, unstructured.Unstructured or map[string]interface{}")
	ERR_TYPE_UPDATE = ERR_TYPE_CREATE
	ERR_TYPE_APPLY  = ERR_TYPE_CREATE
	ERR_TYPE_DELETE = ERR_TYPE_CREATE
	ErrInvalidGetType    = ERR_TYPE_CREATE
)
