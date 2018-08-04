package kube

const (
	KindOfSleep      Kind = "sleep"      // sleeper, a no op
	KindOfPod        Kind = "pods"       // pods
	KindOfPv         Kind = "pv"         // persistent volumes
	KindOfPvc        Kind = "pvc"        // persistent volume claims
	KindOfSvc        Kind = "svc"        // services
	KindOfJob        Kind = "jobs"       // jobs
	KindOfDs         Kind = "ds"         // daemon set
	KindOfNs         Kind = "ns"         // namespace
	KindOfSc         Kind = "sc"         // storage class
	KindOfDeployment Kind = "deployment" // deployment controller
)

const (
	DefaultNamespace = "default"
	NamedNamespace   = "kube"
)

const (
	Sync SyncType = iota
	Async
)

const (
	TypeAssertionError   Error = "type assertion error"
	UnsupportedCoderKind Error = "coder kind not supported"
	UnsupportedSync      Error = "unsupported sync"
)
