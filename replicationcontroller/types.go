package replicationcontroller

import "fmt"

var (
	ErrInvalidToolsType  = fmt.Errorf("type must be string, *corev1.ReplicationController, or corev1.ReplicationController")
	ErrInvalidCreateType = fmt.Errorf("type must be string, []byte, *corev1.ReplicationController, corev1.ReplicationController, runtime.Object, *unstructured.Unstructured, unstructured.Unstructured or map[string]interface{}")
	ErrInvalidUpdateType = ErrInvalidCreateType
	ErrInvalidApplyType  = ErrInvalidCreateType
	ErrInvalidDeleteType = ErrInvalidCreateType
	ErrInvalidGetType    = ErrInvalidCreateType
	ErrInvalidScaleType  = ErrInvalidCreateType
)
