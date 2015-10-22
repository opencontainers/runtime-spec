package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/opencontainers/specs"
	"reflect"
)

func main() {

	var configFile = flag.String("f", "", "input the config file, or default to config.json")
	flag.Parse()

	if len(*configFile) == 0 {
		*configFile = "./config.json"
	}

	var sp specs.Spec
	content, err := ReadFile(*configFile)
	if err != nil {
		return
	}
	json.Unmarshal([]byte(content), &sp)
	var secret interface{} = sp
	value := reflect.ValueOf(secret)

	var err_msg []string
	err_msg, ok := TagStructValid(value, reflect.TypeOf(secret).Name())

	if ok == false {
		fmt.Println("The configuration is incomplete, see the details: ")
		for index := 0; index < len(err_msg); index++ {
			fmt.Println(err_msg[index])
		}
	} else {
		fmt.Println("The configuration is Good")

	}
}
