package main

import (
	"flag"
	"fmt"
	"os"

	"google.golang.org/protobuf/compiler/protogen"
)

var (
	includeMsgMethods *bool
	version           string // This will be set at build time
)

func main() {
	var flags flag.FlagSet

	// Get the flags
	includeMsgMethods = flags.Bool("include_msg_methods", false, "Include getter methods on messages")
	showVersion := flag.Bool("version", false, "Print the version of protoc-gen-go-fieldgetters")
	flag.Parse()

	if *showVersion {
		if version == "" {
			version = "development" // Default version if not provided at build time
		}
		fmt.Printf("%s\n", version)
		os.Exit(0)
	}

	options := protogen.Options{
		ParamFunc: flags.Set,
	}
	options.Run(func(gen *protogen.Plugin) error {
		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}

			if _, err := generateFile(gen, f); err != nil {
				gen.Error(err)
				return err
			}
		}

		return nil
	})
}
