package role

import (
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
)

// ApplyFromRaw apply role from map[string]interface{}.
func (h *Handler) ApplyFromRaw(raw map[string]interface{}) (*rbacv1.Role, error) {
	role := &rbacv1.Role{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(raw, role)
	if err != nil {
		return nil, err
	}

	var namespace string
	if len(role.Namespace) != 0 {
		namespace = role.Namespace
	} else {
		namespace = h.namespace
	}

	role, err = h.clientset.RbacV1().Roles(namespace).Create(h.ctx, role, h.Options.CreateOptions)
	if k8serrors.IsAlreadyExists(err) {
		role, err = h.clientset.RbacV1().Roles(namespace).Update(h.ctx, role, h.Options.UpdateOptions)
	}
	return role, err
}

// ApplyFromBytes apply role from bytes.
func (h *Handler) ApplyFromBytes(data []byte) (role *rbacv1.Role, err error) {
	role, err = h.CreateFromBytes(data)
	if errors.IsAlreadyExists(err) {
		role, err = h.UpdateFromBytes(data)
	}
	return
}

// ApplyFromFile apply role from yaml file.
func (h *Handler) ApplyFromFile(filename string) (role *rbacv1.Role, err error) {
	role, err = h.CreateFromFile(filename)
	if errors.IsAlreadyExists(err) {
		role, err = h.UpdateFromFile(filename)
	}
	return
}

// Apply apply role from yaml file, alias to "ApplyFromFile".
func (h *Handler) Apply(filename string) (*rbacv1.Role, error) {
	return h.ApplyFromFile(filename)
}
