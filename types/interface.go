package types

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type Handler interface {
	Creater
	Updater
	Applyer
	Deleter
	Geter
	Lister
	Watcher
}

type Creater interface {
	//CreateFromRaw(raw map[string]interface{}) (interface{}, error)
	//CreateFromBytes(data []byte) (interface{}, error)
	//CreateFromFile(path string) (interface{}, error)
	Create(name string) (interface{}, error)
}

type Updater interface {
	//UpdateFromRaw(raw map[string]interface{}) (interface{}, error)
	//UpdateFromBytes(data []byte) (interface{}, error)
	//UpdateFromFile(path string) (interface{}, error)
	Update(name string) (interface{}, error)
}

type Applyer interface {
	//ApplyFromRaw(raw map[string]interface{}) (interface{}, error)
	//ApplyFromBytes(data []byte) (interface{}, error)
	//ApplyFromFile(path string) (interface{}, error)
	Apply(name string) (interface{}, error)
}
type Deleter interface {
	//DeleteByName(data []byte) error
	//DeleteFromBytes(data []byte) error
	//DeleteFromFile(path string) error
	Delete(name string) error
}

type Geter interface {
	//GetByName(name string) (interface{}, error)
	//GetFromBytes(name string) (interface{}, error)
	//GetFromFile(path string) (interface{}, error)
	Get(name string) (interface{}, error)
}

type Lister interface {
	//ListByLabel(label string) (interface{}, error)
	//ListAll() (interface{}, error)
	List(label string) (interface{}, error)
}

type Watcher interface {
	//WatchByName(name string, addFunc, modifyFunc, deleteFunc func(x interface{}), x interface{}) error
	//WatchByLabel(label string, addFunc, modifyFunc, deleteFunc func(x interface{}), x interface{}) error
	Watch(name string, addFunc, modifyFunc, deleteFunc func(x interface{}), x interface{}) error
}

// Object is a Kubernetes object, allows functions to work indistinctly with
// any resource that implements both Object interfaces.
//
// Semantically, these are objects which are both serializable (runtime.Object)
// and identifiable (metav1.Object) -- think any object which you could write
// as YAML or JSON, and then `kubectl create`.
//
// Code-wise, this means that any object which embeds both ObjectMeta (which
// provides metav1.Object) and TypeMeta (which provides half of runtime.Object)
// and has a `DeepCopyObject` implementation (the other half of runtime.Object)
// will implement this by default.
//
// For example, nearly all the built-in types are Objects, as well as all
// KubeBuilder-generated CRDs (unless you do something real funky to them).
//
// By and large, most things that implement runtime.Object also implement
// Object -- it's very rare to have *just* a runtime.Object implementation (the
// cases tend to be funky built-in types like Webhook payloads that don't have
// a `metadata` field).
//
// Notice that XYZList types are distinct: they implement ObjectList instead.
type Object interface {
	metav1.Object
	runtime.Object
}

// ObjectList is a Kubernetes object list, allows functions to work
// indistinctly with any resource that implements both runtime.Object and
// metav1.ListInterface interfaces.
//
// Semantically, this is any object which may be serialized (ObjectMeta), and
// is a kubernetes list wrapper (has items, pagination fields, etc) -- think
// the wrapper used in a response from a `kubectl list --output yaml` call.
//
// Code-wise, this means that any object which embedds both ListMeta (which
// provides metav1.ListInterface) and TypeMeta (which provides half of
// runtime.Object) and has a `DeepCopyObject` implementation (the other half of
// runtime.Object) will implement this by default.
//
// For example, nearly all the built-in XYZList types are ObjectLists, as well
// as the XYZList types for all KubeBuilder-generated CRDs (unless you do
// something real funky to them).
//
// By and large, most things that are XYZList and implement runtime.Object also
// implement ObjectList -- it's very rare to have *just* a runtime.Object
// implementation (the cases tend to be funky built-in types like Webhook
// payloads that don't have a `metadata` field).
//
// This is similar to Object, which is almost always implemented by the items
// in the list themselves.
type ObjectList interface {
	metav1.ListInterface
	runtime.Object
}
