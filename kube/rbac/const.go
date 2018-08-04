package rbac

const (
	PackageName       = "kube/sc"
	DefaultConfigDir  = ".config/" + PackageName
	DefaultConfigFile = "config.json"
	Kind              = "ClusterRoleBinding"
	APIVersion        = "rbac.authorization.k8s.io/v1beta1"
)
