package specs

// State holds information about the runtime state of the container.
type State struct {
	// Version is the version of the specification that is supported.
	Version string `json:"version"`
	// ID is the container ID
	ID string `json:"id"`
	// Pid is the process id for the container's main process.
	Pid int `json:"pid"`
	// Root is the path to the container's bundle directory.
	Root string `json:"root"`
}
