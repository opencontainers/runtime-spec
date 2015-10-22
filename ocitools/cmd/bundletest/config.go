package main

import (
	"fmt"
	"github.com/opencontainers/specs"
	"os"
	"path"
	"regexp"
)

func CheckSemVer(version string) (string, bool) {
	re, _ := regexp.Compile("^(\\d+)?\\.(\\d+)?\\.(\\d+)?$")
	if ok := re.Match([]byte(version)); !ok {
		return fmt.Sprintf("%s is not a valid version format, please read 'SemVer v2.0.0'", version), false
	}
	return "", true
}

func SpecValid(s specs.Spec, runtime specs.RuntimeSpec, rootfs string) (msgs []string, valid bool) {
	valid = true
	if len(s.Version) > 0 {
		if m, ok := CheckSemVer(s.Version); !ok {
			msgs = append(msgs, m)
			valid = false
		}
	}

	if len(rootfs) > 0 {
		if ms, ok := MountPointsValid(s.Mounts, runtime.Mounts, rootfs); !ok {
			msgs = append(msgs, ms...)
			valid = false
		}
	}

	return msgs, valid
}

func MountPointsValid(mps []specs.MountPoint, rmps map[string]specs.Mount, rootfs string) (msgs []string, valid bool) {
	valid = true
	for index := 0; index < len(mps); index++ {
		if m, ok := MountPointValid(mps[index], rootfs); !ok {
			msgs = append(msgs, m)
			valid = false
			continue
		}

		if _, ok := rmps[mps[index].Name]; !ok {
			msgs = append(msgs, fmt.Sprintf("%s in config/mount is not exist in runtime/mount", mps[index].Name))
			valid = false
			continue
		}
		//Check if there were duplicated mount name
		for dIndex := index + 1; dIndex < len(mps); dIndex++ {
			if mps[index].Name == mps[dIndex].Name {
				msgs = append(msgs, fmt.Sprintf("%s in config/mount is duplicated", mps[index].Name))
				valid = false
			}
		}
	}
	return msgs, valid
}

func MountPointValid(mp specs.MountPoint, rootfs string) (string, bool) {
	mountPath := path.Join(rootfs, mp.Path)
	fi, err := os.Stat(mountPath)
	if err != nil {
		return fmt.Sprintf("The mountPoint %s %s is not exist in rootfs", mp.Name, mp.Path), false
	} else {
		if !fi.IsDir() {
			return fmt.Sprintf("The mountPoint %s %s is not a valid directory", mp.Name, mp.Path), false
		}
	}
	return "", true
}

func LinuxSpecValid(ls specs.LinuxSpec, runtime specs.LinuxRuntimeSpec, rootfs string) ([]string, bool) {
	msgs, valid := SpecValid(ls.Spec, runtime.RuntimeSpec, rootfs)

	paths := requiredPaths()
	for p_index := 0; p_index < len(paths); p_index++ {
		found := false
		for m_index := 0; m_index < len(ls.Spec.Mounts); m_index++ {
			mp := ls.Spec.Mounts[m_index]
			if paths[p_index] == mp.Path {
				found = true
				break
			}
		}
		if !found {
			msgs = append(msgs, fmt.Sprintf("The mount %s is missing", paths[p_index]))
			valid = false
		}
	}

	for index := 0; index < len(ls.Linux.Capabilities); index++ {
		capability := ls.Linux.Capabilities[index]
		if !capValid(capability) {
			msgs = append(msgs, fmt.Sprintf("%s is not valid, please `man capabilities`", ls.Linux.Capabilities[index]))
			valid = false
		}
	}
	return msgs, valid
}
