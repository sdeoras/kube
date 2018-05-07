package kube

type Kind string
type SyncType int

type Error string

func (e Error) Error() string {
	return string(e)
}
