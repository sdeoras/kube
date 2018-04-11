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
	OrderForward Order = iota
	OrderBackward
)
