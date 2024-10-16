package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/alis-exchange/protoc-gen-fieldgetters/plugin"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
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

	generateMethodMessages := false
	if includeMsgMethods != nil {
		generateMethodMessages = *includeMsgMethods
	}
	options.Run(func(p *protogen.Plugin) error {
		p.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
		return plugin.Generate(p, generateMethodMessages)
	})
}
