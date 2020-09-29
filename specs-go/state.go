package specs

// ContainerState represents the state of a container.
type ContainerState string

const (
	// StateCreating indicates that the container is being created
	StateCreating ContainerState = "creating"

	// StateCreated indicates that the runtime has finished the create operation
	StateCreated ContainerState = "created"

	// StateRunning indicates that the container process has executed the
	// user-specified program but has not exited
	StateRunning ContainerState = "running"

	// StateStopped indicates that the container process has exited
	StateStopped ContainerState = "stopped"
)

// State holds information about the runtime state of the container. The State
// can be displayed when requested (query state operation); it is also passed
// via stdin to many hooks.
type State struct {
	// Version is the version of the specification that is supported.
	Version string `json:"ociVersion"`
	// ID is the container ID
	ID string `json:"id"`
	// Status is the runtime status of the container.
	Status ContainerState `json:"status"`
	// Pid is the process ID for the container process.
	Pid int `json:"pid,omitempty"`
	// Bundle is the path to the container's bundle directory.
	Bundle string `json:"bundle"`
	// Annotations are key values associated with the container.
	Annotations map[string]string `json:"annotations,omitempty"`
}

type SeccompState struct {
	// Version is the version of the specification that is supported.
	Version string `json:"ociVersion"`
	// SeccompFd is the file descriptor for Seccomp User Notification
	SeccompFd int `json:"seccompFd"`
	// Pid is the process ID on which the seccomp filter is applied
	Pid int `json:"pid"`
	// PidFd is a pidfd for the process on which the seccomp filter is
	// applied
	PidFd int `json:"pidFd,omitempty"`
	// State of the container
	State State `json:"state"`
}
