package plugin

import "google.golang.org/protobuf/compiler/protogen"

func Generate(plugin *protogen.Plugin, includeMsgMethods bool) error {
	for _, f := range plugin.Files {
		if !f.Generate {
			continue
		}

		generator := Generator{
			includeMsgMethods: includeMsgMethods,
		}

		if _, err := generator.generateFile(plugin, f); err != nil {
			plugin.Error(err)
			return err
		}
	}

	return nil
}
