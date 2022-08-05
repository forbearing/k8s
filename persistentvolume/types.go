package persistentvolume

import "fmt"

var (
	ERR_TYPE_TOOLS  = fmt.Errorf("type must be string *corev1.PersistentVolume, or corev1.PersistentVolume")
	ERR_TYPE_CREATE = fmt.Errorf("type must be string, []byte, *corev1.PersistentVolume, corev1.PersistentVolume, runtime.Object, *unstructured.Unstructured, unstructured.Unstructured or map[string]interface{}")
	ErrInvalidUpdateType = ERR_TYPE_CREATE
	ErrInvalidApplyType  = ERR_TYPE_CREATE
	ErrInvalidDeleteType = ERR_TYPE_CREATE
	ErrInvalidGetType    = ERR_TYPE_CREATE
)
