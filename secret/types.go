package secret

import "fmt"

var (
	ERR_TYPE_TOOLS  = fmt.Errorf("type must be string *corev1.Secret, or corev1.Secret")
	ERR_TYPE_CREATE = fmt.Errorf("type must be string, []byte, *corev1.Secret, corev1.Secret, runtime.Object or map[string]interface{}")
	ERR_TYPE_UPDATE = ERR_TYPE_CREATE
	ERR_TYPE_APPLY  = ERR_TYPE_CREATE
	ERR_TYPE_DELETE = ERR_TYPE_CREATE
	ERR_TYPE_GET    = ERR_TYPE_CREATE
)
