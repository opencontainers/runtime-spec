package specs

// Spec is the base configuration for the container.  It specifies platform
// independent configuration.
type Spec struct {
	// Version is the version of the specification that is supported.
	Version string `json:"version" mandatory:"required"`
	// Platform is the host information for OS and Arch.
	Platform Platform `json:"platform" mandatory:"required"`
	// Process is the container's main process.
	Process Process `json:"process" mandatory:"required"`
	// Root is the root information for the container's filesystem.
	Root Root `json:"root" mandatory:"required"`
	// Hostname is the container's host name.
	Hostname string `json:"hostname" mandatory:"optional"`
	// Mounts profile configuration for adding mounts to the container's filesystem.
	Mounts []MountPoint `json:"mounts" mandatory:"optional"`
}

// Process contains information to start a specific application inside the container.
type Process struct {
	// Terminal creates an interactive terminal for the container.
	Terminal bool `json:"terminal" mandatory:"optional"`
	// User specifies user information for the process.
	User User `json:"user" mandatory:"required"`
	// Args specifies the binary and arguments for the application to execute.
	Args []string `json:"args" mandatory:"required"`
	// Env populates the process environment for the process.
	Env []string `json:"env" mandatory:"optional"`
	// Cwd is the current working directory for the process and must be
	// relative to the container's root.
	Cwd string `json:"cwd" mandatory:"optional"`
}

// Root contains information about the container's root filesystem on the host.
type Root struct {
	// Path is the absolute path to the container's root filesystem.
	Path string `json:"path" mandatory:"required"`
	// Readonly makes the root filesystem for the container readonly before the process is executed.
	Readonly bool `json:"readonly" mandatory:"optional"`
}

// Platform specifies OS and arch information for the host system that the container
// is created for.
type Platform struct {
	// OS is the operating system.
	OS string `json:"os" mandatory:"required"`
	// Arch is the architecture
	Arch string `json:"arch" mandatory:"required"`
}

// MountPoint describes a directory that may be fullfilled by a mount in the runtime.json.
type MountPoint struct {
	// Name is a unique descriptive identifier for this mount point.
	Name string `json:"name" mandatory:"required"`
	// Path specifies the path of the mount. The path and child directories MUST exist, a runtime MUST NOT create directories automatically to a mount point.
	Path string `json:"path" mandatory:"required"`
}

// State holds information about the runtime state of the container.
type State struct {
	// Version is the version of the specification that is supported.
	Version string `json:"version" mandatory:"required"`
	// ID is the container ID
	ID string `json:"id" mandatory:"required"`
	// Pid is the process id for the container's main process.
	Pid int `json:"pid" mandatory:"required"`
	// BundlePath is the path to the container's bundle directory.
	BundlePath string `json:"bundlePath" mandatory:"required"`
}
