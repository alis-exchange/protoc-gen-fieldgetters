# protoc-gen-fieldgetters
Custom protoc plugin that generates proto field getters. This is primarily useful for the https://pkg.go.dev/go.alis.build/validator package.

> [!IMPORTANT]   
> This plugin is designed to be used alongside the [protoc-gen-go](https://pkg.go.dev/google.golang.org/protobuf/cmd/protoc-gen-go) and [protoc-gen-go-grpc](https://pkg.go.dev/google.golang.org/grpc/cmd/protoc-gen-go-grpc) plugins. It is not designed to be used standalone.

This plugin generates getters for each field in a proto message. It adds a `{Message}_FieldGetters` struct which contains the different getters.
A new instance of `{Message}_FieldGetters` can be created by calling the `New{Message}FieldGetters()` function.

> [!NOTE]  
> Note: {Message} is the name of the proto message.

You can the use it as so

```go
// Create a new instance of the field getters
fieldGetters := New{Message}FieldGetters()

fieldGetters.StringGetter(&pb.Message{}, "field_name")
```

## Installation

```shell
go install github.com/alis-exchange/protoc-gen-fieldgetters/cmd/protoc-gen-fieldgetters@latest
```

## Usage

```shell
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative --fieldgetters_out=. --fieldgetters_opt=paths=source_relative,include_msg_methods=true path/to/your.proto
```

### Supported getters

1. `StringGetter` - Returns the string value of the field. `func(msg protoreflect.ProtoMessage, path string) (string, error)`
2. `StringListGetter` - Returns the list of strings of the field. `func(msg protoreflect.ProtoMessage, path string) ([]string, error)`
3. `BoolGetter` - Returns the bool value of the field. `func(msg protoreflect.ProtoMessage, path string) (bool, error)`
4. `BoolListGetter` - Returns the list of bools of the field. `func(msg protoreflect.ProtoMessage, path string) ([]bool, error)`
5. `IntGetter` - Returns the int value of the field. `func(msg protoreflect.ProtoMessage, path string) (int64, error)`
6. `IntListGetter` - Returns the list of ints of the field. `func(msg protoreflect.ProtoMessage, path string) ([]int64, error)`
7. `FloatGetter` - Returns the float value of the field. `func(msg protoreflect.ProtoMessage, path string) (float64, error)`
8. `FloatListGetter` - Returns the list of floats of the field. `func(msg protoreflect.ProtoMessage, path string) ([]float64, error)`
9. `EnumGetter` - Returns the enum value of the field. `func(msg protoreflect.ProtoMessage, path string) (protoreflect.EnumNumber, error)`
10. `EnumListGetter` - Returns the list of enums of the field. `func(msg protoreflect.ProtoMessage, path string) ([]protoreflect.EnumNumber, error)`
11. `SubMessageGetter` - Returns the submessage of the field. `func(msg protoreflect.ProtoMessage, path string) (protoreflect.ProtoMessage, error)`
12. `SubMessageListGetter` - Returns the list of submessages of the field. `func(msg protoreflect.ProtoMessage, path string) ([]protoreflect.ProtoMessage, error)`

### Flags

- `--fieldgetters_out` - Output directory for the generated files.
- `--fieldgetters_opt` - Comma-separated list of options to pass to the fieldgetters plugin.
  - `paths` - {string} -
  - `include_msg_methods` - {bool} - Include the methods on the message itself. Default: false
  
   ```go
   func (m *Message) StringGetter(path string) (string, error) {
       ...
   }
   ```
   Which can be used as so
   ```go
   myMessage := &pb.Message{}
  
   value, err := myMessage.StringGetter("field_name")
   ```


## Contributing

