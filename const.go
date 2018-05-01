package kube

const (
	KindOfSleep Kind = "sleep" // sleeper, a no op
	KindOfPod   Kind = "pods"  // pods
	KindOfPv    Kind = "pv"    // persistent volumes
	KindOfPvc   Kind = "pvc"   // persistent volume claims
	KindOfSvc   Kind = "svc"   // services
	KindOfJob   Kind = "jobs"  // jobs
	KindOfDs    Kind = "ds"    // daemon set
	KindOfNs    Kind = "ns"    // namespace
)

const (
	DefaultNamespace = "default"
	NamedNamespace   = "kube"
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
