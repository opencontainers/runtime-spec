package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/opencontainers/specs"
	"os"
	"path"
)

const (
	// Path to config file inside the bundle
	ConfigFile  = "config.json"
	RuntimeFile = "runtime.json"
	// Path to rootfs directory inside the bundle
	RootfsDir = "rootfs"
)

type Bundle struct {
	Config  specs.Spec
	Runtime specs.RuntimeSpec
	Rootfs  string
}

type LinuxBundle struct {
	Config  specs.LinuxSpec
	Runtime specs.LinuxRuntimeSpec
	Rootfs  string
}

func ReadFile(file_url string) (content string, err error) {
	_, err = os.Stat(file_url)
	if err != nil {
		fmt.Println("cannot find the file ", file_url)
		return content, err
	}
	file, err := os.Open(file_url)
	defer file.Close()
	if err != nil {
		fmt.Println("cannot open the file ", file_url)
		return content, err
	}
	buf := bytes.NewBufferString("")
	buf.ReadFrom(file)
	content = buf.String()

	return content, nil
}

func OSDetect(inputPath string) string {
	var configURL string
	fi, err := os.Stat(inputPath)
	if err == nil {
		if fi.IsDir() {
			configURL = path.Join(inputPath, ConfigFile)
		} else {
			configURL = inputPath
		}
	} else {
		return ""
	}
	content, err := ReadFile(configURL)
	if err == nil {
		var s specs.Spec
		err = json.Unmarshal([]byte(content), &s)
		if err == nil {
			return s.Platform.OS
		}
	}
	return ""
}

func FilesValid(bundlePath string) (msgs []string, valid bool) {
	valid = true
	fi, err := os.Stat(bundlePath)
	if err != nil {
		msgs = append(msgs, fmt.Sprintf("Error accessing bundle: %v", err))
		return msgs, false
	} else {
		if !fi.IsDir() {
			msgs = append(msgs, fmt.Sprintf("Given path %s is not a directory", bundlePath))
			return msgs, false
		}
	}

	configPath := path.Join(bundlePath, ConfigFile)
	_, err = os.Stat(configPath)
	if err != nil {
		msgs = append(msgs, fmt.Sprintf("Error accessing %s: %v", ConfigFile, err))
		valid = false
	}

	runtimePath := path.Join(bundlePath, RuntimeFile)
	_, err = os.Stat(runtimePath)
	if err != nil {
		msgs = append(msgs, fmt.Sprintf("Error accessing %s: %v", RuntimeFile, err))
		valid = false
	}

	rootfsPath := path.Join(bundlePath, RootfsDir)
	fi, err = os.Stat(rootfsPath)
	if err != nil {
		msgs = append(msgs, fmt.Sprintf("Error accessing %s: %v", RootfsDir, err))
		valid = false
	} else {
		if !fi.IsDir() {
			msgs = append(msgs, fmt.Sprintf("Given path %s is not a directory", rootfsPath))
			valid = false
		}
	}

	return msgs, valid
}

func ConfigValid(configPath string) (msgs []string, valid bool) {
	valid = true
	os := OSDetect(configPath)
	if len(os) == 0 {
		msgs = append(msgs, "Cannot detect OS in the config.json under the bundle, or maybe miss `config.json`.")
		return msgs, false
	}
	if os == "linux" {
		var ls specs.LinuxSpec
		var rt specs.LinuxRuntimeSpec
		content, _ := ReadFile(configPath)
		json.Unmarshal([]byte(content), &ls)

		var secret interface{} = ls
		if ms, ok := TagValid(secret); !ok {
			msgs = append(msgs, ms...)
			valid = false
		}
		if ms, ok := LinuxSpecValid(ls, rt, ""); !ok {
			msgs = append(msgs, ms...)
			valid = false
		}
	} else {
		var s specs.Spec
		var rt specs.RuntimeSpec
		content, _ := ReadFile(configPath)
		json.Unmarshal([]byte(content), &s)

		var secret interface{} = s
		if ms, ok := TagValid(secret); !ok {
			msgs = append(msgs, ms...)
			valid = false
		}
		if ms, ok := SpecValid(s, rt, ""); !ok {
			msgs = append(msgs, ms...)
			valid = false
		}
	}
	return msgs, valid
}

func RuntimeValid(runtimePath string, os string, rootfs string) (msgs []string, valid bool) {
	valid = true
	content, err := ReadFile(runtimePath)
	if err != nil {
		msgs = append(msgs, fmt.Sprintf("Cannot read %s", runtimePath))
		return msgs, false
	}
	if os == "linux" {
		var lrt specs.LinuxRuntimeSpec
		err = json.Unmarshal([]byte(content), &lrt)
		if err != nil {
			msgs = append(msgs, fmt.Sprintf("Cannot parse %s", runtimePath))
			valid = false
		} else {
			var secret interface{} = lrt
			if ms, ok := TagValid(secret); !ok {
				msgs = append(msgs, ms...)
				valid = false
			}
			if ms, ok := LinuxRuntimeSpecValid(lrt, rootfs); !ok {
				msgs = append(msgs, ms...)
				valid = false
			}
		}
	} else {
		var rt specs.RuntimeSpec
		err = json.Unmarshal([]byte(content), &rt)
		if err != nil {
			msgs = append(msgs, fmt.Sprintf("Cannot parse %s", runtimePath))
			valid = false
		} else {
			var secret interface{} = rt
			if ms, ok := TagValid(secret); !ok {
				msgs = append(msgs, ms...)
				valid = false
			}
			if ms, ok := RuntimeSpecValid(rt, rootfs); !ok {
				msgs = append(msgs, ms...)
				valid = false
			}
		}
	}
	return msgs, valid
}

func BundleValid(bundlePath string) ([]string, bool) {
	msgs, valid := FilesValid(bundlePath)
	if valid == false {
		return msgs, false
	}

	os := OSDetect(bundlePath)
	if len(os) == 0 {
		msgs = append(msgs, "Cannot detect OS in the config.json under the bundle, or maybe miss `config.json`.")
		return msgs, false
	}

	if ms, ok := ConfigValid(path.Join(bundlePath, ConfigFile)); !ok {
		msgs = append(msgs, ms...)
		valid = false
	}
	if ms, ok := RuntimeValid(path.Join(bundlePath, RuntimeFile), os, path.Join(bundlePath, RootfsDir)); !ok {
		msgs = append(msgs, ms...)
		valid = false
	}

	return msgs, valid
}
