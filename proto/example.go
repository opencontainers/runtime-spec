// +build ignore

package main

import (
	"encoding/json"
	"log"

	oci "./go/"
	"github.com/golang/protobuf/proto"
)

func main() {
	s := oci.LinuxSpec{
		Spec: &oci.Spec{
			Platform: &oci.Platform{Os: proto.String("linux"), Arch: proto.String("x86_64")},
			Process: &oci.Process{
				Cwd: proto.String("/"),
				Env: []string{"TERM=linux"},
			},
		},
	}

	buf, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	println(string(buf))
}
