package specs

// RuntimeSpec is the generic runtime state information on a running container
type RuntimeSpec struct {
	// Mounts is a mapping of names to mount configurations.
	// Which mounts will be mounted and where should be chosen with MountPoints
	// in Spec.
	Mounts map[string]Mount `json:"mounts" mandatory:"required"`
	// Hooks are the commands run at various lifecycle events of the container.
	Hooks Hooks `json:"hooks" mandatory:"optional"`
}

// Hook specifies a command that is run at a particular event in the lifecycle of a container
type Hook struct {
	Path string   `json:"path" mandatory:"required"`
	Args []string `json:"args" mandatory:"optional"`
	Env  []string `json:"env" mandatory:"optional"`
}

// Hooks for container setup and teardown
type Hooks struct {
	// Prestart is a list of hooks to be run before the container process is executed.
	// On Linux, they are run after the container namespaces are created.
	Prestart []Hook `json:"prestart" mandatory:"optional"`
	// Poststart is a list of hooks to be run after the container process is started.
	Poststart []Hook `json:"poststart" mandatory:"optional"`
	// Poststop is a list of hooks to be run after the container process exits.
	Poststop []Hook `json:"poststop" mandatory:"optional"`
}

// Mount specifies a mount for a container
type Mount struct {
	// Type specifies the mount kind.
	Type string `json:"type" mandatory:"required"`
	// Source specifies the source path of the mount.  In the case of bind mounts on
	// linux based systems this would be the file on the host.
	Source string `json:"source" mandatory:"required"`
	// Options are fstab style mount options.
	Options []string `json:"options" mandatory:"optional"`
}
