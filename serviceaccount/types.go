package serviceaccount

import "fmt"

var (
	ERR_TYPE_TOOLS  = fmt.Errorf("type must be string *corev1.ServiceAccount, or corev1.ServiceAccount")
	ERR_TYPE_CREATE = fmt.Errorf("type must be string, []byte, *corev1.ServiceAccount, corev1.ServiceAccount, runtime.Object or map[string]interface{}")
	ERR_TYPE_UPDATE = ERR_TYPE_CREATE
	ERR_TYPE_APPLY  = ERR_TYPE_CREATE
	ERR_TYPE_DELETE = ERR_TYPE_CREATE
	ERR_TYPE_GET    = ERR_TYPE_CREATE
)