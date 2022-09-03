package k8s

type Options int

const (
	IgnoreAlreadyExists Options = iota
	IgnoreNotFound
	IgnoreInvalid
	IgnoreTimeout
)
