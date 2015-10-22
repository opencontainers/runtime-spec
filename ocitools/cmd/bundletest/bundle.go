package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"os"
)

func parseBundle(context *cli.Context) {
	if len(context.Args()) > 0 {
		if msgs, ok := BundleValid(context.Args()[0]); !ok {
			fmt.Println("The bundle is not valid, details errors:")
			for index := 0; index < len(msgs); index++ {
				fmt.Println(msgs[index])
			}
		} else {
			fmt.Println("The bundle is valid!")
		}
	} else {
		cli.ShowCommandHelp(context, "bundle")
	}
}

func parseConfig(context *cli.Context) {
	if len(context.Args()) > 0 {
		if msgs, ok := ConfigValid(context.Args()[0]); !ok {
			fmt.Println("The config.json is not valid, details errors:")
			for index := 0; index < len(msgs); index++ {
				fmt.Println(msgs[index])
			}
		} else {
			fmt.Println("The config.json is valid!")
		}
	} else {
		cli.ShowCommandHelp(context, "config")
	}
}

func parseRuntime(context *cli.Context) {
	if len(context.Args()) > 0 {
		if msgs, ok := RuntimeValid(context.Args()[0], "linux", ""); !ok {
			fmt.Println("The runtime.json is not valid, details errors:")
			for index := 0; index < len(msgs); index++ {
				fmt.Println(msgs[index])
			}
		} else {
			fmt.Println("The runtime.json is valid!")
		}
	} else {
		cli.ShowCommandHelp(context, "runtime")
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "Bundle Validator"
	app.Usage = "Standard Container Validator: tool to validate if a `bundle` was a standand container"
	app.Version = "0.2.0"
	app.Commands = []cli.Command{
		{
			Name:    "bundle",
			Aliases: []string{"vb"},
			Usage:   "Validate all the config.json, runtime.json and files in the rootfs",
			Action:  parseBundle,
		},
		{
			Name:    "config",
			Aliases: []string{"vc"},
			Usage:   "Validate the config.json only",
			Action:  parseConfig,
		},
		{
			Name:    "runtime",
			Aliases: []string{"vr"},
			Usage:   "Validate the runtime.json only, runtime + os, default to 'linux'",
			Action:  parseRuntime,
		},
	}

	app.Run(os.Args)

	return
}
