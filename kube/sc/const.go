package sc

const (
	PackageName       = "kube/sc"
	DefaultConfigDir  = ".config/" + PackageName
	DefaultConfigFile = "config.json"
	Kind              = "StorageClass"
	APIVersion        = "storage.k8s.io/v1beta1"

	AWSElasticBlockStore = "kubernetes.io/aws-ebs"
	AzureFile            = "kubernetes.io/azure-file"
	AzureDisk            = "kubernetes.io/azure-disk"
	CephFS               = ""
	Cinder               = "kubernetes.io/cinder"
	FC                   = ""
	FlexVolume           = ""
	Flocker              = ""
	GCEPersistentDisk    = "kubernetes.io/gce-pd"
	Glusterfs            = ""
	ISCSI                = ""
	Quobyte              = ""
	NFS                  = ""
	RBD                  = ""
	VsphereVolume        = ""
	PortworxVolume       = "kubernetes.io/portworx-volume"
	ScaleIO              = ""
	StorageOS            = ""
	Local                = "kubernetes.io/no-provisioner"
)
