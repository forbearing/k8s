package clusterrolebinding

import (
	"time"

	rbacv1 "k8s.io/api/rbac/v1"
)

type Role struct {
	Kind string
	Name string
}

type Subject struct {
	Kind      string
	Name      string
	Namespace string
}

// GetRole get the clusterrole binding by this clusterrolebinding.
func (h *Handler) GetRole(object interface{}) (*Role, error) {
	switch val := object.(type) {
	case string:
		crb, err := h.Get(val)
		if err != nil {
			return nil, err
		}
		return &Role{Kind: crb.RoleRef.Kind, Name: crb.RoleRef.Name}, nil
	case *rbacv1.ClusterRoleBinding:
		return &Role{Kind: val.RoleRef.Kind, Name: val.RoleRef.Name}, nil
	case rbacv1.ClusterRoleBinding:
		return &Role{Kind: val.RoleRef.Kind, Name: val.RoleRef.Name}, nil
	default:
		return nil, ErrInvalidToolsType
	}
}

// GetAge returns clusterrolebinding age.
func (h *Handler) GetAge(object interface{}) (time.Duration, error) {
	switch val := object.(type) {
	case string:
		crb, err := h.Get(val)
		if err != nil {
			return time.Duration(0), err
		}
		return time.Now().Sub(crb.CreationTimestamp.Time), nil
	case *rbacv1.ClusterRoleBinding:
		return time.Now().Sub(val.CreationTimestamp.Time), nil
	case rbacv1.ClusterRoleBinding:
		return time.Now().Sub(val.CreationTimestamp.Time), nil
	default:
		return time.Duration(0), ErrInvalidToolsType
	}
}

// GetSubjects get the subjects to which the clusterrole applies.
// All supported subject kinds are: User, Group, ServiceAccount.
func (h *Handler) GetSubjects(object interface{}) ([]Subject, error) {
	switch val := object.(type) {
	case string:
		crb, err := h.Get(val)
		if err != nil {
			return nil, err
		}
		return h.getSubjects(crb), nil
	case *rbacv1.ClusterRoleBinding:
		return h.getSubjects(val), nil
	case rbacv1.ClusterRoleBinding:
		return h.getSubjects(&val), nil
	default:
		return nil, ErrInvalidToolsType
	}
}
func (h *Handler) getSubjects(crb *rbacv1.ClusterRoleBinding) []Subject {
	var sl []Subject
	for _, subject := range crb.Subjects {
		s := Subject{
			Kind:      subject.Kind,
			Name:      subject.Name,
			Namespace: subject.Namespace,
		}
		sl = append(sl, s)
	}
	return sl
}
