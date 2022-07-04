package replicationcontroller

import (
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
)

// ApplyFromRaw apply replicationcontroller from map[string]interface{}
func (h *Handler) ApplyFromRaw(raw map[string]interface{}) (*corev1.ReplicationController, error) {
	rc := &corev1.ReplicationController{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(raw, rc)
	if err != nil {
		return nil, err
	}

	var namespace string
	if len(rc.Namespace) != 0 {
		namespace = rc.Namespace
	} else {
		namespace = h.namespace
	}

	rc, err = h.clientset.CoreV1().ReplicationControllers(namespace).Create(h.ctx, rc, h.Options.CreateOptions)
	if k8serrors.IsAlreadyExists(err) {
		rc, err = h.clientset.CoreV1().ReplicationControllers(namespace).Update(h.ctx, rc, h.Options.UpdateOptions)
	}
	return rc, err
}

// ApplyFromBytes apply replicationcontroller from bytes
func (h *Handler) ApplyFromBytes(data []byte) (rc *corev1.ReplicationController, err error) {
	rc, err = h.CreateFromBytes(data)
	if k8serrors.IsAlreadyExists(err) {
		log.Debug(err)
		rc, err = h.UpdateFromBytes(data)
	}
	return
}

// ApplyFromFile apply replicationcontroller from yaml file
func (h *Handler) ApplyFromFile(filename string) (rc *corev1.ReplicationController, err error) {
	rc, err = h.CreateFromFile(filename)
	if k8serrors.IsAlreadyExists(err) {
		rc, err = h.UpdateFromFile(filename)
	}
	return
}

// Apply apply replicationcontroller from yaml file, alias to "ApplyFromFile"
func (h *Handler) Apply(filename string) (*corev1.ReplicationController, error) {
	return h.ApplyFromFile(filename)
}
