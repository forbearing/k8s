package clusterrolebinding

import (
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
)

// ApplyFromRaw apply clusterrolebinding from map[string]interface{}.
func (h *Handler) ApplyFromRaw(raw map[string]interface{}) (*rbacv1.ClusterRoleBinding, error) {
	crb := &rbacv1.ClusterRoleBinding{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(raw, crb)
	if err != nil {
		return nil, err
	}

	crb, err = h.clientset.RbacV1().ClusterRoleBindings().Create(h.ctx, crb, h.Options.CreateOptions)
	if k8serrors.IsAlreadyExists(err) {
		crb, err = h.clientset.RbacV1().ClusterRoleBindings().Update(h.ctx, crb, h.Options.UpdateOptions)
	}
	return crb, err
}

// ApplyFromBytes apply clusterrolebinding from bytes.
func (h *Handler) ApplyFromBytes(data []byte) (crb *rbacv1.ClusterRoleBinding, err error) {
	crb, err = h.CreateFromBytes(data)
	if errors.IsAlreadyExists(err) {
		crb, err = h.UpdateFromBytes(data)
	}
	return
}

// ApplyFromFile apply clusterrolebinding from yaml file.
func (h *Handler) ApplyFromFile(filename string) (crb *rbacv1.ClusterRoleBinding, err error) {
	crb, err = h.CreateFromFile(filename)
	if errors.IsAlreadyExists(err) {
		crb, err = h.UpdateFromFile(filename)
	}
	return
}

// Apply apply clusterrolebinding from file, alias to "ApplyFromFile".
func (h *Handler) Apply(filename string) (*rbacv1.ClusterRoleBinding, error) {
	return h.ApplyFromFile(filename)
}
