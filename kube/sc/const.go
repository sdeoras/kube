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

// GCP persistent disks volume creation params.
const (
	CGP_PD_Type          = "type"
	GCP_PD_Type_Standard = "pd-standard"
	GCP_PD_Type_SSD      = "pd-ssd"
)

// Portworx volume creation constants for storage class
const (
	// Filesystem to be laid out: none, xfs, ext4
	PX_FS      = "fs"
	PX_FS_None = "none"
	PX_FS_Xfs  = "xfs"
	PX_FS_Ext4 = "ext4"

	// Block size
	PX_BlockSize    = "block_size"
	PX_BlockSize_4K = "4k"

	// Replication factor for the volume: 1, 2, 3
	PX_Repl   = "repl"
	PX_Repl_1 = "1"
	PX_Repl_2 = "2"
	PX_Repl_3 = "3"

	// Shared volume allowing RWMany
	PX_Shared       = "shared"
	PX_Shared_True  = "true"
	PX_Shared_False = "false"

	// IO Priority: low, medium, high
	PX_PriorityIO        = "priority_io"
	PX_PriorityIO_Low    = "low"
	PX_PriorityIO_Medium = "medium"
	PX_PriorityIO_High   = "high"

	// IO Profile can be used to override the
	// I/O algorithm Portworx uses for the volumes.
	// Supported values are db, sequential, random, cms
	PX_IOProfile            = "io_profile"
	PX_IOProfile_Db         = "db"
	PX_IOProfile_Sequential = "sequential"
	PX_IOProfile_Random     = "random"
	PX_IOProfile_CMS        = "cms"

	// The group a volume should belong too.
	// Portworx will restrict replication sets
	// of volumes of the same group on different nodes.
	// If the force group option ‘fg’ is set to true,
	// the volume group rule will be strictly enforced.
	// By default, it’s not strictly enforced.
	PX_Group = "group"

	// This option enforces volume group policy.
	// If a volume belonging to a group cannot find nodes
	// for it’s replication sets which don’t have other
	// volumes of same group, the volume creation will fail.
	PX_Fg       = "fg"
	PX_Fg_True  = "true"
	PX_Fg_False = "false"

	// List of comma-separated name=value pairs to apply
	// to the Portworx volume
	PX_Label = "label"

	// Comma-separated Portworx Node ID’s to use for
	// replication sets of the volume
	PX_Nodes = "nodes"

	// Specifies the number of replication sets the
	// volume can be aggregated from
	PX_AggregationLevel = "aggregation_level"

	// Snapshot schedule (PX 1.3 and higher). Following are the accepted formats:
	//
	//periodic=mins,snaps-to-keep
	//daily=hh:mm,snaps-to-keep
	//weekly=weekday@hh:mm,snaps-to-keep
	//monthly=day@hh:mm,snaps-to-keep
	//
	//snaps-to-keep is optional.
	// Periodic, Daily, Weekly and Monthly keep
	// last 5, 7, 5 and 12 snapshots by default
	// respectively.
	PX_SnapSchedule = "snap_schedule"

	// Snapshot interval in minutes.
	// 0 disables snaps. Minimum value: 60.
	// It is recommended to use snap_schedule above.
	PX_SnapInterval = "snap_interval"

	// Flag to create sticky volumes that cannot be
	// deleted until the flag is disabled
	PX_Sticky       = "sticky"
	PX_Sticky_True  = "true"
	PX_Sticky_False = "false"

	// (PX 1.3 and higher) Flag to indicate if you want
	// to use journal device for the volume’s metadata.
	// This will use the journal device that you used
	// when installing Portworx. As of PX version 1.3,
	// it is recommended to use a journal device to absorb
	// PX metadata writes. Default: false
	PX_Journal       = "journal"
	PX_Journal_True  = "true"
	PX_Journal_False = "false"
)
