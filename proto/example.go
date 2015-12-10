// +build ignore

package main

import (
	"encoding/hex"
	"encoding/json"
	"log"

	oci "./go/"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

func main() {
	s := &oci.LinuxSpec{
		Spec: &oci.Spec{
			Version:  proto.String("0.3.0"),
			Hostname: proto.String("darkstar"),
			Platform: &oci.Platform{Os: proto.String("linux"), Arch: proto.String("x86_64")},
			Process: &oci.Process{
				Terminal: proto.Bool(true),
				User:     &oci.User{},
				Cwd:      proto.String("/"),
				Args:     []string{"/bin/sh"},
				Env:      []string{"TERM=linux"},
			},
			Root: &oci.Root{
				Path:     proto.String("/"),
				Readonly: proto.Bool(false),
			},
			Mounts: []*oci.MountPoint{
				&oci.MountPoint{
					Name: proto.String("proc"),
					Path: proto.String("/proc"),
				},
				&oci.MountPoint{
					Name: proto.String("dev"),
					Path: proto.String("/dev"),
				},
				&oci.MountPoint{
					Name: proto.String("devpts"),
					Path: proto.String("/dev/pts"),
				},
				&oci.MountPoint{
					Name: proto.String("shm"),
					Path: proto.String("/dev/shm"),
				},
				&oci.MountPoint{
					Name: proto.String("mqueue"),
					Path: proto.String("/dev/mqueue"),
				},
				&oci.MountPoint{
					Name: proto.String("sysfs"),
					Path: proto.String("/sys"),
				},
				&oci.MountPoint{
					Name: proto.String("cgroup"),
					Path: proto.String("/sys/fs/cgroup"),
				},
			},
		},
		LinuxConfig: &oci.LinuxConfig{
			Capabilities: []string{
				"CAP_AUDIT_WRITE",
				"CAP_KILL",
				"CAP_NET_BIND_SERVICE",
			},
		},
	}

	//proto.SetExtension(s.Spec, oci.E_Uid, 0)

	println("## Using github.com/golang/protobuf/jsonpb to marshal")
	m := jsonpb.Marshaler{}
	jsonStr, err := m.MarshalToString(s)
	if err != nil {
		log.Fatal(err)
	}
	println(jsonStr)
	print("## len: ")
	println(len(jsonStr))
	println("")

	println("## Using encoding/json to marshal")
	buf, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	println(string(buf))
	print("## len: ")
	println(len(buf))
	println("")

	println("## Marshaling to protobuf binary message")
	data, err := proto.Marshal(s)
	if err != nil {
		log.Fatal(err)
	}
	println(hex.Dump(data))
	print("## len: ")
	println(len(data))
}
