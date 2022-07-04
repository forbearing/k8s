package clusterrole

import (
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
)

// ApplyFromRaw apply clusterrole from map[string]interface{}.
func (h *Handler) ApplyFromRaw(raw map[string]interface{}) (*rbacv1.ClusterRole, error) {
	cr := &rbacv1.ClusterRole{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(raw, cr)
	if err != nil {
		return nil, err
	}

	cr, err = h.clientset.RbacV1().ClusterRoles().Create(h.ctx, cr, h.Options.CreateOptions)
	if k8serrors.IsAlreadyExists(err) {
		cr, err = h.clientset.RbacV1().ClusterRoles().Update(h.ctx, cr, h.Options.UpdateOptions)
	}
	return cr, err
}

// ApplyFromBytes apply clusterrole from bytes.
func (h *Handler) ApplyFromBytes(data []byte) (clusterrole *rbacv1.ClusterRole, err error) {
	clusterrole, err = h.CreateFromBytes(data)
	if errors.IsAlreadyExists(err) {
		clusterrole, err = h.UpdateFromBytes(data)
	}
	return
}

// ApplyFromFile apply clusterrole from yaml file.
func (h *Handler) ApplyFromFile(filename string) (clusterrole *rbacv1.ClusterRole, err error) {
	clusterrole, err = h.CreateFromFile(filename)
	if errors.IsAlreadyExists(err) {
		clusterrole, err = h.UpdateFromFile(filename)
	}
	return
}

// Apply apply clusterrole from yaml file, alias to "ApplyFromFile".
func (h *Handler) Apply(filename string) (*rbacv1.ClusterRole, error) {
	return h.ApplyFromFile(filename)
}
