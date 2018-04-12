package kube

const (
	KindOfPod Kind = "pod"
	KindOfPv  Kind = "persistent-volume"
	KindOfPvc Kind = "persistent-volume-claim"
	KindOfSvc Kind = "service"
	KindOfJob Kind = "job"
	KindOfDs  Kind = "daemonset"
)

const (
	DefaultNamespace = "default"
)

const (
	Forward Order = iota
	Backward
)

const (
	TypeAssertionError Error = "type assertion error"
)
