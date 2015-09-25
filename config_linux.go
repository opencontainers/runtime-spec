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
}
