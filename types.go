package kube

type Kind string
type Order int

type Error string

func (e Error) Error() string {
	return string(e)
}
