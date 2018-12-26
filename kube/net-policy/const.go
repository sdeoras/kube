package netpolicy

const (
	PackageName       = "kube/netpolicy"
	DefaultConfigDir  = ".config/" + PackageName
	DefaultConfigFile = "config.json"
	Kind              = "NetworkPolicy"
	ApiVersion        = "networking.k8s.io/v1"
)
