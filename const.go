package kube

const (
	KindOfPod Kind = "pods"
	KindOfPv  Kind = "pv"
	KindOfPvc Kind = "pvc"
	KindOfSvc Kind = "svc"
	KindOfJob Kind = "jobs"
	KindOfDs  Kind = "ds"
)

const (
	DefaultNamespace = "default"
)

const (
	Forward Order = iota
	Backward
)

const (
	TypeAssertionError   Error = "type assertion error"
	UnsupportedCoderKind Error = "coder kind not supported"
)
