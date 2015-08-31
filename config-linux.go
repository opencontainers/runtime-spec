// +build linux

package specs

// LinuxSpec is the full specification for linux containers.
type LinuxSpec struct {
	Spec
	// Linux is platform specific configuration for linux based containers.
	Linux Linux `json:"linux"`
}

// Linux contains platform specific configuration for linux based containers.
type Linux struct {
	// Capabilities are linux capabilities that are kept for the container.
	Capabilities []string `json:"capabilities"`
	// RootfsPropagation is the rootfs mount propagation mode for the container.
	RootfsPropagation string `json:"rootfsPropagation"`
}

// User specifies linux specific user and group information for the container's
// main process.
type User struct {
	// Uid is the user id.
	Uid int32 `json:"uid"`
	// Gid is the group id.
	Gid int32 `json:"gid"`
	// AdditionalGids are additional group ids set for the container's process.
	AdditionalGids []int32 `json:"additionalGids"`
}
