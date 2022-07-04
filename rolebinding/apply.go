package rolebinding

import (
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
)

// ApplyFromRaw apply rolebinding from map[string]interface{}.
func (h *Handler) ApplyFromRaw(raw map[string]interface{}) (*rbacv1.RoleBinding, error) {
	rolebinding := &rbacv1.RoleBinding{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(raw, rolebinding)
	if err != nil {
		return nil, err
	}

	var namespace string
	if len(rolebinding.Namespace) != 0 {
		namespace = rolebinding.Namespace
	} else {
		namespace = h.namespace
	}

	rolebinding, err = h.clientset.RbacV1().RoleBindings(namespace).Create(h.ctx, rolebinding, h.Options.CreateOptions)
	if k8serrors.IsAlreadyExists(err) {
		rolebinding, err = h.clientset.RbacV1().RoleBindings(namespace).Update(h.ctx, rolebinding, h.Options.UpdateOptions)
	}
	return rolebinding, err
}

// ApplyFromBytes apply rolebinding from bytes.
func (h *Handler) ApplyFromBytes(data []byte) (rolebinding *rbacv1.RoleBinding, err error) {
	rolebinding, err = h.CreateFromBytes(data)
	if errors.IsAlreadyExists(err) {
		rolebinding, err = h.UpdateFromBytes(data)
	}
	return
}

// ApplyFromFile apply rolebinding from yaml file.
func (h *Handler) ApplyFromFile(filename string) (rolebinding *rbacv1.RoleBinding, err error) {
	rolebinding, err = h.CreateFromFile(filename)
	if errors.IsAlreadyExists(err) {
		rolebinding, err = h.UpdateFromFile(filename)
	}
	return
}

// Apply apply rolebinding from yaml file, alias to "ApplyFromFile".
func (h *Handler) Apply(filename string) (*rbacv1.RoleBinding, error) {
	return h.ApplyFromFile(filename)
}
