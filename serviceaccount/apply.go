package serviceaccount

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
)

// ApplyFromRaw apply serviceaccount from map[string]interface{}.
func (h *Handler) ApplyFromRaw(raw map[string]interface{}) (*corev1.ServiceAccount, error) {
	sa := &corev1.ServiceAccount{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(raw, sa)
	if err != nil {
		return nil, err
	}

	var namespace string
	if len(sa.Namespace) != 0 {
		namespace = sa.Namespace
	} else {
		namespace = h.namespace
	}

	_, err = h.clientset.CoreV1().ServiceAccounts(namespace).Create(h.ctx, sa, h.Options.CreateOptions)
	if k8serrors.IsAlreadyExists(err) {
		sa, err = h.clientset.CoreV1().ServiceAccounts(namespace).Update(h.ctx, sa, h.Options.UpdateOptions)
	}
	return sa, err
}

// ApplyFromBytes apply serviceaccount from file.
func (h *Handler) ApplyFromBytes(data []byte) (sa *corev1.ServiceAccount, err error) {
	sa, err = h.CreateFromBytes(data)
	if errors.IsAlreadyExists(err) {
		sa, err = h.UpdateFromBytes(data)
	}
	return
}

// ApplyFromFile apply serviceaccount from yaml file.
func (h *Handler) ApplyFromFile(filename string) (sa *corev1.ServiceAccount, err error) {
	sa, err = h.CreateFromFile(filename)
	if errors.IsAlreadyExists(err) {
		sa, err = h.UpdateFromFile(filename)
	}
	return
}

// Apply apply serviceaccount from yaml file, alias to "ApplyFromFile".
func (h *Handler) Apply(filename string) (*corev1.ServiceAccount, error) {
	return h.ApplyFromFile(filename)
}
