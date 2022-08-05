package persistentvolume

import "fmt"

var (
	ERR_TYPE_TOOLS  = fmt.Errorf("type must be string *corev1.PersistentVolume, or corev1.PersistentVolume")
	ERR_TYPE_CREATE = fmt.Errorf("type must be string, []byte, *corev1.PersistentVolume, corev1.PersistentVolume, runtime.Object, *unstructured.Unstructured, unstructured.Unstructured or map[string]interface{}")
	ERR_TYPE_UPDATE = ERR_TYPE_CREATE
	ERR_TYPE_APPLY  = ERR_TYPE_CREATE
	ERR_TYPE_DELETE = ERR_TYPE_CREATE
	ErrInvalidGetType    = ERR_TYPE_CREATE
)
