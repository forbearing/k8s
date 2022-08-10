package ownerreference

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/utils/pointer"
)

// AlreadyOwnedError is an error returned if the object you are trying to assign
// a controller reference is already owned by another controller Object is the
// subject and Owner is the reference for the current owner.
type AlreadyOwnedError struct {
	Object metav1.Object
	Owner  metav1.OwnerReference
}

func (e *AlreadyOwnedError) Error() string {
	return fmt.Sprintf("Object %s/%s is already owned by another %s controller %s", e.Object.GetNamespace(), e.Object.GetName(), e.Owner.Kind, e.Owner.Name)
}

func newAlreadyOwnedError(obj metav1.Object, owner metav1.OwnerReference) *AlreadyOwnedError {
	return &AlreadyOwnedError{
		Object: obj,
		Owner:  owner,
	}
}

// SetControllerReference sets owner as a Controller OwnerReference on controlled.
// This is used for garbage collection of the controlled object and for
// reconciling the owner object on changes to controlled (with a Watch + EnqueueRequestForOwner).
// Since only one OwnerReference can be a controller, it returns an error if
// there is another OwnerReference with Controller flag set.
func SetControllerReference(owner, controlled metav1.Object, scheme *runtime.Scheme) error {
	// Validate the owner.
	ro, ok := owner.(runtime.Object)
	if !ok {
		return fmt.Errorf("%T is not a runtime.Object, cannot call SetControllerReference", owner)
	}
	if err := validateOwner(owner, controlled); err != nil {
		return err
	}

	// Create a new controller ref.
	gvk, err := GVKForObject(ro, scheme)
	if err != nil {
		return err
	}
	ref := metav1.OwnerReference{
		APIVersion:         gvk.GroupVersion().String(),
		Kind:               gvk.Kind,
		Name:               owner.GetName(),
		UID:                owner.GetUID(),
		BlockOwnerDeletion: pointer.BoolPtr(true),
		Controller:         pointer.BoolPtr(true),
	}

	// Return early with an error if the object is already controlled.
	if existing := metav1.GetControllerOf(controlled); existing != nil && !referSameObject(*existing, ref) {
		return newAlreadyOwnedError(controlled, *existing)
	}

	// Update owner references and return.
	upsertOwnerRef(ref, controlled)
	return nil
}

// SetOwnerReference is a helper method to make sure the given object contains an object reference to the object provided.
// This allows you to declare that owner has a dependency on the object without specifying it as a controller.
// If a reference to the same object already exists, it'll be overwritten with the newly provided version.
func SetOwnerReference(owner, object metav1.Object, scheme *runtime.Scheme) error {
	// Validate the owner.
	ro, ok := owner.(runtime.Object)
	if !ok {
		return fmt.Errorf("%T is not a runtime.Object, cannot call SetOwnerReference", owner)
	}
	if err := validateOwner(owner, object); err != nil {
		return err
	}

	// Create a new owner ref.
	gvk, err := GVKForObject(ro, scheme)
	if err != nil {
		return err
	}
	ref := metav1.OwnerReference{
		APIVersion: gvk.GroupVersion().String(),
		Kind:       gvk.Kind,
		UID:        owner.GetUID(),
		Name:       owner.GetName(),
	}

	// Update owner references and return.
	upsertOwnerRef(ref, object)
	return nil
}

func upsertOwnerRef(ref metav1.OwnerReference, object metav1.Object) {
	owners := object.GetOwnerReferences()
	if idx := indexOwnerRef(owners, ref); idx == -1 {
		owners = append(owners, ref)
	} else {
		owners[idx] = ref
	}
	object.SetOwnerReferences(owners)
}

// indexOwnerRef returns the index of the owner reference in the slice if found, or -1.
func indexOwnerRef(ownerReferences []metav1.OwnerReference, ref metav1.OwnerReference) int {
	for index, r := range ownerReferences {
		if referSameObject(r, ref) {
			return index
		}
	}
	return -1
}

func validateOwner(owner, object metav1.Object) error {
	ownerNs := owner.GetNamespace()
	if ownerNs != "" {
		objNs := object.GetNamespace()
		if objNs == "" {
			return fmt.Errorf("cluster-scoped resource must not have a namespace-scoped owner, owner's namespace %s", ownerNs)
		}
		if ownerNs != objNs {
			return fmt.Errorf("cross-namespace owner references are disallowed, owner's namespace %s, obj's namespace %s", owner.GetNamespace(), object.GetNamespace())
		}
	}
	return nil
}

// Returns true if a and b point to the same object.
func referSameObject(a, b metav1.OwnerReference) bool {
	aGV, err := schema.ParseGroupVersion(a.APIVersion)
	if err != nil {
		return false
	}

	bGV, err := schema.ParseGroupVersion(b.APIVersion)
	if err != nil {
		return false
	}

	return aGV.Group == bGV.Group && a.Kind == b.Kind && a.Name == b.Name
}

// GVKForObject finds the GroupVersionKind associated with the given object, if there is only a single such GVK.
func GVKForObject(obj runtime.Object, scheme *runtime.Scheme) (schema.GroupVersionKind, error) {
	// TODO(directxman12): do we want to generalize this to arbitrary container types?
	// I think we'd need a generalized form of scheme or something.  It's a
	// shame there's not a reliable "GetGVK" interface that works by default
	// for unpopulated static types and populated "dynamic" types
	// (unstructured, partial, etc)

	// check for PartialObjectMetadata, which is analogous to unstructured, but isn't handled by ObjectKinds
	_, isPartial := obj.(*metav1.PartialObjectMetadata) //nolint:ifshort
	_, isPartialList := obj.(*metav1.PartialObjectMetadataList)
	if isPartial || isPartialList {
		// we require that the GVK be populated in order to recognize the object
		gvk := obj.GetObjectKind().GroupVersionKind()
		if len(gvk.Kind) == 0 {
			return schema.GroupVersionKind{}, runtime.NewMissingKindErr("unstructured object has no kind")
		}
		if len(gvk.Version) == 0 {
			return schema.GroupVersionKind{}, runtime.NewMissingVersionErr("unstructured object has no version")
		}
		return gvk, nil
	}

	gvks, isUnversioned, err := scheme.ObjectKinds(obj)
	if err != nil {
		return schema.GroupVersionKind{}, err
	}
	if isUnversioned {
		return schema.GroupVersionKind{}, fmt.Errorf("cannot create group-version-kind for unversioned type %T", obj)
	}

	if len(gvks) < 1 {
		return schema.GroupVersionKind{}, fmt.Errorf("no group-version-kinds associated with type %T", obj)
	}
	if len(gvks) > 1 {
		// this should only trigger for things like metav1.XYZ --
		// normal versioned types should be fine
		return schema.GroupVersionKind{}, fmt.Errorf(
			"multiple group-version-kinds associated with type %T, refusing to guess at one", obj)
	}
	return gvks[0], nil
}
