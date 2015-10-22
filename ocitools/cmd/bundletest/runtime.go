package main

import (
	"fmt"
	"github.com/opencontainers/specs"
	"os"
	"path"
)

func RuntimeSpecValid(rs specs.RuntimeSpec, rootfs string) ([]string, bool) {
	return nil, true
}

func LinuxRuntimeSpecValid(lrs specs.LinuxRuntimeSpec, rootfs string) ([]string, bool) {
	msgs, valid := RuntimeSpecValid(lrs.RuntimeSpec, rootfs)
	lr := lrs.Linux

	if len(lr.UIDMappings) > 5 {
		msgs = append(msgs, "The UID mapping is limited to 5")
		valid = false
	}
	if len(lr.GIDMappings) > 5 {
		msgs = append(msgs, "The GID mapping is limited to 5")
		valid = false
	}

	for index := 0; index < len(lr.Rlimits); index++ {
		if ms, ok := RlimitValid(lr.Rlimits[index]); !ok {
			msgs = append(msgs, ms...)
			valid = false
		}
	}

	for index := 0; index < len(lr.Namespaces); index++ {
		if ms, ok := NamespaceValid(lr.Namespaces[index]); !ok {
			msgs = append(msgs, ms...)
			valid = false
		}
	}

	//minimum devices
	devices := requiredDevices()
	for index := 0; index < len(devices); index++ {
		found := false
		for dIndex := 0; dIndex < len(lr.Devices); dIndex++ {
			if lr.Devices[dIndex].Path == devices[index] {
				found = true
				break
			}
		}
		if found == false {
			msgs = append(msgs, fmt.Sprintf("The required device %s is missing", devices[index]))
			valid = false
		}
	}

	for index := 0; index < len(lr.Devices); index++ {
		if ms, ok := DeviceValid(lr.Devices[index]); !ok {
			msgs = append(msgs, ms...)
			valid = false
		}
	}

	if len(lr.ApparmorProfile) > 0 && len(rootfs) > 0 {
		profilePath := path.Join(rootfs, "/etc/apparmor.d", lr.ApparmorProfile)
		_, err := os.Stat(profilePath)
		if err != nil {
			msgs = append(msgs, fmt.Sprintf("ApparmorProfile %s is not exist", lr.ApparmorProfile))
			valid = false
		}
	}

	switch lr.RootfsPropagation {
	case "":
	case "slave":
	case "private":
	case "shared":
		break
	default:
		valid = false
		msgs = append(msgs, "RootfsPropagation should limited to 'slave', 'private', or 'shared'")
		break

	}

	return msgs, valid
}

func NamespaceValid(ns specs.Namespace) (msgs []string, valid bool) {
	valid = true
	switch ns.Type {
	case "":
		valid = false
		msgs = append(msgs, "The type of the namespace should not be empty")
		break
	case "pid":
	case "network":
	case "mount":
	case "ipc":
	case "uts":
	case "user":
		break
	default:
		valid = false
		msgs = append(msgs, "The type of the namespace should limited to 'pid/network/mount/ipc/nts/user'")
		break
	}
	return msgs, valid
}

func RlimitValid(r specs.Rlimit) (msgs []string, valid bool) {
	if !rlimitValid(r.Type) {
		msgs = append(msgs, "Rlimit is invalid")
		return msgs, false
	}
	return msgs, true
}

func DeviceValid(d specs.Device) (msgs []string, valid bool) {
	valid = true
	switch d.Type {
	case 'b':
	case 'c':
	case 'u':
		if d.Major <= 0 {
			msgs = append(msgs, fmt.Sprintf("Device %s type is `b/c/u`, please set the major number", d.Path))
			valid = false
		}
		if d.Minor <= 0 {
			msgs = append(msgs, fmt.Sprintf("Device %s type is `b/c/u`, please set the minor number", d.Path))
			valid = false
		}
	case 'p':
		if d.Major > 0 || d.Minor > 0 {
			msgs = append(msgs, fmt.Sprintf("Device %s type is `p`, no need to set major/minor number", d.Path))
			valid = false
		}
	default:
		msgs = append(msgs, fmt.Sprintf("Device %s type should limited to `b/c/u/p`", d.Path))
		valid = false
	}
	return msgs, valid
}

func SeccompValid(s specs.Seccomp) (msgs []string, valid bool) {
	valid = true
	if !seccompValid(string(s.DefaultAction)) {
		msgs = append(msgs, "Seccomp.DefaultAction is invalid")
		valid = false
	}
	for index := 0; index < len(s.Syscalls); index++ {
		if s.Syscalls[index] != nil {
			if ms, ok := SyscallValid(*(s.Syscalls[index])); !ok {
				msgs = append(msgs, ms...)
				valid = false
			}
		}
	}
	return msgs, valid
}

func SyscallValid(s specs.Syscall) (msgs []string, valid bool) {
	valid = true
	if !seccompValid(string(s.Action)) {
		msgs = append(msgs, fmt.Sprintf("Syscall.Action %s is invalid", s.Action))
		valid = false
	}

	return msgs, valid
}
