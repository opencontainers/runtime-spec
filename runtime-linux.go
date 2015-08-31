package specs

import "os"

// LinuxRuntimeSpec is the full specification for linux containers.
type LinuxRuntimeSpec struct {
	RuntimeSpec
	// Linux is platform specific configuration for linux based containers.
	LinuxRuntime Linux `json:"linux"`
}

type LinuxRuntime struct {
	// UidMapping specifies user mappings for supporting user namespaces on linux.
	UidMappings []IDMapping `json:"uidMappings"`
	// UidMapping specifies group mappings for supporting user namespaces on linux.
	GidMappings []IDMapping `json:"gidMappings"`
	// Rlimits specifies rlimit options to apply to the container's process.
	Rlimits []Rlimit `json:"rlimits"`
	// Sysctl are a set of key value pairs that are set for the container on start
	Sysctl map[string]string `json:"sysctl"`
	// Resources contain cgroup information for handling resource constraints
	// for the container
	Resources Resources `json:"resources"`
	// Namespaces contains the namespaces that are created and/or joined by the container
	Namespaces []Namespace `json:"namespaces"`
	// Devices are a list of device nodes that are created and enabled for the container
	Devices []Device `json:"devices"`
	// ApparmorProfile specified the apparmor profile for the container.
	ApparmorProfile string `json:"apparmorProfile"`
	// SelinuxProcessLabel specifies the selinux context that the container process is run as.
	SelinuxProcessLabel string `json:"selinuxProcessLabel"`
	// Seccomp specifies the seccomp security settings for the container.
	Seccomp Seccomp `json:"seccomp"`
	// RootfsPropagation is the rootfs mount propagation mode for the container
	RootfsPropagation string `json:"rootfsPropagation"`
}

// Namespace is the configuration for a linux namespace.
type Namespace struct {
	// Type is the type of Linux namespace
	Type string `json:"type"`
	// Path is a path to an existing namespace persisted on disk that can be joined
	// and is of the same type
	Path string `json:"path"`
}

// IDMapping specifies UID/GID mappings
type IDMapping struct {
	// HostID is the UID/GID of the host user or group
	HostID int32 `json:"hostID"`
	// ContainerID is the UID/GID of the container's user or group
	ContainerID int32 `json:"containerID"`
	// Size is the length of the range of IDs mapped between the two namespaces
	Size int32 `json:"size"`
}

// Rlimit type and restrictions
type Rlimit struct {
	// Type of the rlimit to set
	Type int `json:"type"`
	// Hard is the hard limit for the specified type
	Hard uint64 `json:"hard"`
	// Soft is the soft limit for the specified type
	Soft uint64 `json:"soft"`
}

// HugepageLimit structure corresponds to limiting kernel hugepages
type HugepageLimit struct {
	Pagesize string `json:"pageSize"`
	Limit    int    `json:"limit"`
}

// InterfacePriority for network interfaces
type InterfacePriority struct {
	// Name is the name of the network interface
	Name string `json:"name"`
	// Priority for the interface
	Priority int64 `json:"priority"`
}

// BlockIO for Linux cgroup 'blockio' resource management
type BlockIO struct {
	// Specifies per cgroup weight, range is from 10 to 1000
	Weight int64 `json:"blkioWeight"`
	// Weight per cgroup per device, can override BlkioWeight
	WeightDevice string `json:"blkioWeightDevice"`
	// IO read rate limit per cgroup per device, bytes per second
	ThrottleReadBpsDevice string `json:"blkioThrottleReadBpsDevice"`
	// IO write rate limit per cgroup per divice, bytes per second
	ThrottleWriteBpsDevice string `json:"blkioThrottleWriteBpsDevice"`
	// IO read rate limit per cgroup per device, IO per second
	ThrottleReadIOpsDevice string `json:"blkioThrottleReadIopsDevice"`
	// IO write rate limit per cgroup per device, IO per second
	ThrottleWriteIOpsDevice string `json:"blkioThrottleWriteIopsDevice"`
}

// Memory for Linux cgroup 'memory' resource management
type Memory struct {
	// Memory limit (in bytes)
	Limit int64 `json:"limit"`
	// Memory reservation or soft_limit (in bytes)
	Reservation int64 `json:"reservation"`
	// Total memory usage (memory + swap); set `-1' to disable swap
	Swap int64 `json:"swap"`
	// Kernel memory limit (in bytes)
	Kernel int64 `json:"kernel"`
	// How aggressive the kernel will swap memory pages. Range from 0 to 100. Set -1 to use system default
	Swappiness int64 `json:"swappiness"`
}

// CPU for Linux cgroup 'cpu' resource management
type CPU struct {
	// CPU shares (relative weight vs. other cgroups with cpu shares)
	Shares int64 `json:"shares"`
	// CPU hardcap limit (in usecs). Allowed cpu time in a given period
	Quota int64 `json:"quota"`
	// CPU period to be used for hardcapping (in usecs). 0 to use system default
	Period int64 `json:"period"`
	// How many time CPU will use in realtime scheduling (in usecs)
	RealtimeRuntime int64 `json:"realtimeRuntime"`
	// CPU period to be used for realtime scheduling (in usecs)
	RealtimePeriod int64 `json:"realtimePeriod"`
	// CPU to use within the cpuset
	Cpus string `json:"cpus"`
	// MEM to use within the cpuset
	Mems string `json:"mems"`
}

// Network identification and priority configuration
type Network struct {
	// Set class identifier for container's network packets
	ClassID string `json:"classId"`
	// Set priority of network traffic for container
	Priorities []InterfacePriority `json:"priorities"`
}

// Resources has container runtime resource constraints
type Resources struct {
	// DisableOOMKiller disables the OOM killer for out of memory conditions
	DisableOOMKiller bool `json:"disableOOMKiller"`
	// Memory restriction configuration
	Memory Memory `json:"memory"`
	// CPU resource restriction configuration
	CPU CPU `json:"cpu"`
	// BlockIO restriction configuration
	BlockIO BlockIO `json:"blockIO"`
	// Hugetlb limit (in bytes)
	HugepageLimits []HugepageLimit `json:"hugepageLimits"`
	// Network restriction configuration
	Network Network `json:"network"`
}

type Device struct {
	// Path to the device.
	Path string `json:"path"`
	// Device type, block, char, etc.
	Type rune `json:"type"`
	// Major is the device's major number.
	Major int64 `json:"major"`
	// Minor is the device's minor number.
	Minor int64 `json:"minor"`
	// Cgroup permissions format, rwm.
	Permissions string `json:"permissions"`
	// FileMode permission bits for the device.
	FileMode os.FileMode `json:"fileMode"`
	// UID of the device.
	UID uint32 `json:"uid"`
	// Gid of the device.
	GID uint32 `json:"gid"`
}

// Seccomp represents syscall restrictions
type Seccomp struct {
	DefaultAction Action     `json:"defaultAction"`
	Syscalls      []*Syscall `json:"syscalls"`
}

// Action taken upon Seccomp rule match
type Action string

// Operator used to match syscall arguments in Seccomp
type Operator string

// Arg used for matching specific syscall arguments in Seccomp
type Arg struct {
	Index    uint     `json:"index"`
	Value    uint64   `json:"value"`
	ValueTwo uint64   `json:"valueTwo"`
	Op       Operator `json:"op"`
}

// Syscall is used to match a syscall in Seccomp
type Syscall struct {
	Name   string `json:"name"`
	Action Action `json:"action"`
	Args   []*Arg `json:"args"`
}
