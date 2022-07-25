package types

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
