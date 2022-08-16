package finalizer

import (
	"github.com/forbearing/k8s/types"
)

// AddFinalizer accepts an Object and adds the provided finalizer if not present.
// It returns an indication of whether it updated the object's list of finalizers.
func AddFinalizer(o types.Object, finalizer string) (finalizersUpdated bool) {
	f := o.GetFinalizers()
	for _, e := range f {
		if e == finalizer {
			return false
		}
	}
	o.SetFinalizers(append(f, finalizer))
	return true
}

// RemoveFinalizer accepts an Object and removes the provided finalizer if present.
// It returns an indication of whether it updated the object's list of finalizers.
func RemoveFinalizer(o types.Object, finalizer string) (finalizersUpdated bool) {
	f := o.GetFinalizers()
	for i := 0; i < len(f); i++ {
		if f[i] == finalizer {
			f = append(f[:i], f[i+1:]...)
			i--
			finalizersUpdated = true
		}
	}
	o.SetFinalizers(f)
	return
}

// ContainsFinalizer checks an Object that the provided finalizer is present.
func ContainsFinalizer(o types.Object, finalizer string) bool {
	f := o.GetFinalizers()
	for _, e := range f {
		if e == finalizer {
			return true
		}
	}
	return false
}
