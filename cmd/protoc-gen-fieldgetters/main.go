package main

import (
	"flag"

	"google.golang.org/protobuf/compiler/protogen"
)

var (
	includeMsgMethods *bool
)

func main() {
	var flags flag.FlagSet

	// Get the flags
	includeMsgMethods = flags.Bool("include_msg_methods", false, "Include getter methods on messages")

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
