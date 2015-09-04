package specs

// Bundle is the metadata of the files in the container.  It specifies 
// the attributes of files.
type Bundle struct {
	// Files metadata of the files in the container.
	// Use the 'Device' in runtime_config_linux.go
	Files []Device `json:"files"`
}
