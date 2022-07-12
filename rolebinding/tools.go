package rolebinding

import (
	"fmt"
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

var ERR_TYPE = fmt.Errorf("type must be *rbacv1.RoleBinding, rbacv1.RoleBinding or string")

// GetRole get the role/clusterrole binding by this rolebinding.
func (h *Handler) GetRole(object interface{}) (*Role, error) {
	switch val := object.(type) {
	case string:
		rb, err := h.Get(val)
		if err != nil {
			return nil, err
		}
		return &Role{Kind: rb.RoleRef.Kind, Name: rb.RoleRef.Name}, nil
	case *rbacv1.RoleBinding:
		return &Role{Kind: val.RoleRef.Kind, Name: val.RoleRef.Name}, nil
	case rbacv1.RoleBinding:
		return &Role{Kind: val.RoleRef.Kind, Name: val.RoleRef.Name}, nil
	default:
		return nil, ERR_TYPE
	}

}

// GetAge returns rolebinding age.
func (h *Handler) GetAge(object interface{}) (time.Duration, error) {
	switch val := object.(type) {
	case string:
		rb, err := h.Get(val)
		if err != nil {
			return time.Duration(int64(0)), err
		}
		return time.Now().Sub(rb.CreationTimestamp.Time), nil
	case *rbacv1.RoleBinding:
		return time.Now().Sub(val.CreationTimestamp.Time), nil
	case rbacv1.RoleBinding:
		return time.Now().Sub(val.CreationTimestamp.Time), nil
	default:
		return time.Duration(int64(0)), ERR_TYPE
	}
}

// GetSubjects get the subjects to which the role/clusterrole applies.
// All supported subject kinds are: User, Group, ServiceAccount.
func (h *Handler) GetSubjects(object interface{}) ([]Subject, error) {
	switch val := object.(type) {
	case string:
		rb, err := h.Get(val)
		if err != nil {
			return nil, err
		}
		return h.getSubjects(rb), nil
	case *rbacv1.RoleBinding:
		return h.getSubjects(val), nil
	case rbacv1.RoleBinding:
		return h.getSubjects(&val), nil
	default:
		return nil, ERR_TYPE
	}
}
func (h *Handler) getSubjects(rb *rbacv1.RoleBinding) []Subject {
	var sl []Subject
	for _, subject := range rb.Subjects {
		s := Subject{
			Kind:      subject.Kind,
			Name:      subject.Name,
			Namespace: subject.Namespace,
		}
		sl = append(sl, s)
	}
	return sl
}
