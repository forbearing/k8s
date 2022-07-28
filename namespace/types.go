package namespace

import "fmt"

var (
	ERR_TYPE_TOOLS  = fmt.Errorf("type must be string *corev1.Namespace, or corev1.Namespace")
	ERR_TYPE_CREATE = fmt.Errorf("type must be string, []byte, *corev1.Namespace, corev1.Namespace, runtime.Object, *unstructured.Unstructured, unstructured.Unstructured or map[string]interface{}")
	ERR_TYPE_UPDATE = ERR_TYPE_CREATE
	ERR_TYPE_APPLY  = ERR_TYPE_CREATE
	ERR_TYPE_DELETE = ERR_TYPE_CREATE
	ERR_TYPE_GET    = ERR_TYPE_CREATE
)
