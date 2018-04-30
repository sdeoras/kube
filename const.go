package kube

const (
	KindOfNoop Kind = "noop"
	KindOfPod  Kind = "pods"
	KindOfPv   Kind = "pv"
	KindOfPvc  Kind = "pvc"
	KindOfSvc  Kind = "svc"
	KindOfJob  Kind = "jobs"
	KindOfDs   Kind = "ds"
)

const (
	DefaultNamespace = "default"
)

const (
	Forward Order = iota
	Backward
	Async
)

const (
	TypeAssertionError   Error = "type assertion error"
	UnsupportedCoderKind Error = "coder kind not supported"
	UnsupportedOrder     Error = "unsupported exec order"
)
