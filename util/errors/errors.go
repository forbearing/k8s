package errors

import (
	"context"
	"errors"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
)

// IgnoreNotFound returns nil on NotFound errors.
// All other values that are not NotFound errors or nil are returned unmodified.
func IgnoreNotFound(err error) error {
	if apierrors.IsNotFound(err) {
		return nil
	}
	return err
}

// IgnoreAlreadyExists returns nil on AlreadyExists errors.
// All other values that are not AlreadyExists errors or nil are returned unmodified.
func IgnoreAlreadyExists(err error) error {
	if apierrors.IsAlreadyExists(err) {
		return nil
	}
	return err
}

// IgnoreDeadlineExceeded returns nil on context.DeadlineExceeded errors.
// All other values that are not context.DeadlineExceeded or nil are returned unmodified.
func IgnoreDeadlineExceeded(err error) error {
	if errors.Is(err, context.DeadlineExceeded) {
		return nil
	}
	return err
}
