package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/alis-exchange/protoc-gen-fieldgetters/cmd/protoc-gen-go-fieldgetters/utils"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type FieldsGetter struct {
	gen    *protogen.GeneratedFile
	msg    *protogen.Message
	kinds  []protoreflect.Kind
	isList bool
}

func generateFile(gen *protogen.Plugin, file *protogen.File) (*protogen.GeneratedFile, error) {
	filename := file.GeneratedFilenamePrefix + "_getters.pb.go"
	g := gen.NewGeneratedFile(filename, file.GoImportPath)
	g.P("// Code generated by protoc-gen-go-fieldgetters. DO NOT EDIT.")
	g.P("// ")
	g.P(fmt.Sprintf("// Generated on: %s UTC", time.Now().UTC().Format("2006-01-02 15:04:05")))
	g.P(fmt.Sprintf("// Source: %s", file.Desc.Path()))
	g.P(fmt.Sprintf("// protoc-gen-go-fieldgetters --version: %s", version))

	g.P()

	// Add package name
	g.P(fmt.Sprintf("package %s", file.GoPackageName))

	// Add imports
	g.P("import (")
	g.P("\"google.golang.org/protobuf/reflect/protoreflect\"")
	g.P(")")
	g.P()

	// Generate getter functions for messages
	for _, msg := range file.Messages {

		var generateResourceMethods bool
		if includeMsgMethods != nil {
			generateResourceMethods = *includeMsgMethods
		}
		err := generateMessageGetters(g, msg, generateResourceMethods)
		if err != nil {
			return nil, err
		}
	}

	return g, nil
}

// processFields processes fields of a message and generates getter functions for them
func (sg *FieldsGetter) processFields(fields []*protogen.Field) error {

	for _, field := range fields {

		if field.Desc.IsList() && sg.isList && utils.Contains(sg.kinds, field.Desc.Kind()) {
			sg.gen.P(fmt.Sprintf("case \"%s\":", field.Desc.Name()))

			switch field.Desc.Kind() {
			case protoreflect.EnumKind:
				{
					sg.gen.P(fmt.Sprintf("%s_list := make([]protoreflect.EnumNumber, len(%sMsg.Get%s()))", field.Desc.Name(), strings.ToLower(sg.msg.GoIdent.GoName), field.GoName))
					sg.gen.P(fmt.Sprintf("for i, item_%s := range %sMsg.Get%s() {", strings.ToLower(field.GoName), strings.ToLower(sg.msg.GoIdent.GoName), field.GoName))
					sg.gen.P(fmt.Sprintf("%s_list[i] = item_%s.Number()", field.Desc.Name(), strings.ToLower(field.GoName)))
					sg.gen.P("}")
					sg.gen.P(fmt.Sprintf("return %s_list, nil", field.Desc.Name()))
					continue
				}
			case protoreflect.Int32Kind:
				{
					sg.gen.P(fmt.Sprintf("%s_list := make([]int64, len(%sMsg.Get%s()))", field.Desc.Name(), strings.ToLower(sg.msg.GoIdent.GoName), field.GoName))
					sg.gen.P(fmt.Sprintf("for i, item_%s := range %sMsg.Get%s() {", strings.ToLower(field.GoName), strings.ToLower(sg.msg.GoIdent.GoName), field.GoName))
					sg.gen.P(fmt.Sprintf("%s_list[i] = int64(item_%s)", field.Desc.Name(), strings.ToLower(field.GoName)))
					sg.gen.P("}")
					sg.gen.P(fmt.Sprintf("return %s_list, nil", field.Desc.Name()))
					continue
				}
			case protoreflect.FloatKind:
				{
					sg.gen.P(fmt.Sprintf("%s_list := make([]float64, len(%sMsg.Get%s()))", field.Desc.Name(), strings.ToLower(sg.msg.GoIdent.GoName), field.GoName))
					sg.gen.P(fmt.Sprintf("for i, item_%s := range %sMsg.Get%s() {", strings.ToLower(field.GoName), strings.ToLower(sg.msg.GoIdent.GoName), field.GoName))
					sg.gen.P(fmt.Sprintf("%s_list[i] = float64(item_%s)", field.Desc.Name(), strings.ToLower(field.GoName)))
					sg.gen.P("}")
					sg.gen.P(fmt.Sprintf("return %s_list, nil", field.Desc.Name()))
					continue
				}
			case protoreflect.MessageKind:
				{
					sg.gen.P(fmt.Sprintf("%s_list := make([]protoreflect.ProtoMessage, len(%sMsg.Get%s()))", field.Desc.Name(), strings.ToLower(sg.msg.GoIdent.GoName), field.GoName))
					sg.gen.P(fmt.Sprintf("for i, item_%s := range %sMsg.Get%s() {", strings.ToLower(field.Message.GoIdent.GoName), strings.ToLower(sg.msg.GoIdent.GoName), field.GoName))
					sg.gen.P(fmt.Sprintf("%s_list[i] = item_%s", field.Desc.Name(), strings.ToLower(field.Message.GoIdent.GoName)))
					sg.gen.P("}")
					sg.gen.P(fmt.Sprintf("return %s_list, nil", field.Desc.Name()))
					continue
				}
			default:
				{
					sg.gen.P(fmt.Sprintf("return %sMsg.Get%s(), nil", strings.ToLower(sg.msg.GoIdent.GoName), field.GoName))
					continue
				}
			}
		}

		if field.Desc.IsList() || sg.isList {
			continue
		}

		if field.Desc.IsMap() {
			continue
		}

		if field.Desc.IsExtension() {
			continue
		}

		if field.Desc.Kind() == protoreflect.MessageKind && !utils.Contains(sg.kinds, protoreflect.MessageKind) {
			err := sg.processMessage([]*protogen.Field{field}, field.Message)
			if err != nil {
				return err
			}
		}

		if !utils.Contains(sg.kinds, field.Desc.Kind()) {
			continue
		}

		sg.gen.P(fmt.Sprintf("case \"%s\":", field.Desc.Name()))

		returnStr := fmt.Sprintf("%sMsg.Get%s()", strings.ToLower(sg.msg.GoIdent.GoName), field.GoName)
		if field.Desc.Kind() == protoreflect.Int32Kind {
			returnStr = fmt.Sprintf("int64(%s)", returnStr)
		}
		if field.Desc.Kind() == protoreflect.FloatKind {
			returnStr = fmt.Sprintf("float64(%s)", returnStr)
		}
		if field.Desc.Kind() == protoreflect.EnumKind {
			returnStr = fmt.Sprintf("%s.Number()", returnStr)
		}
		sg.gen.P(fmt.Sprintf("return %s, nil", returnStr))

	}

	return nil
}

// processMessage processes fields of a message and generates getter functions for them
func (sg *FieldsGetter) processMessage(parentFields []*protogen.Field, message *protogen.Message) error {
	for _, field := range message.Fields {

		if field.Desc.IsList() {
			continue
		}

		if field.Desc.IsMap() {
			continue
		}

		if field.Desc.IsExtension() {
			continue
		}

		if field.Desc.Kind() == protoreflect.MessageKind && !utils.Contains(sg.kinds, protoreflect.MessageKind) {
			err := sg.processMessage(append(parentFields, field), field.Message)
			if err != nil {
				return err
			}
		}

		if !utils.Contains(sg.kinds, field.Desc.Kind()) {
			continue
		}

		caseStr := strings.Join(utils.Transform(parentFields, func(f *protogen.Field) string {
			return fmt.Sprintf("%s", f.Desc.Name())
		}), ".")
		sg.gen.P(fmt.Sprintf("case \"%s.%s\":", caseStr, field.Desc.Name()))

		traverseStr := ""
		for _, m := range utils.Transform(parentFields, func(f *protogen.Field) string {
			return fmt.Sprintf("%s", f.GoName)

		}) {
			traverseStr += fmt.Sprintf(".Get%s()", m)
		}

		returnStr := fmt.Sprintf("%sMsg%s.Get%s()", strings.ToLower(sg.msg.GoIdent.GoName), traverseStr, field.GoName)
		if field.Desc.Kind() == protoreflect.Int32Kind {
			returnStr = fmt.Sprintf("int64(%s)", returnStr)
		}
		if field.Desc.Kind() == protoreflect.FloatKind {
			returnStr = fmt.Sprintf("float64(%s)", returnStr)
		}
		if field.Desc.Kind() == protoreflect.EnumKind {
			returnStr = fmt.Sprintf("%s.Number()", returnStr)
		}
		sg.gen.P(fmt.Sprintf("return %s, nil", returnStr))

	}

	return nil
}

// fieldsLen returns the number of fields that match the given kinds
func (sg *FieldsGetter) fieldsLen(fields []*protogen.Field) int {
	count := 0
	for _, field := range fields {
		if field.Desc.IsList() && sg.isList && utils.Contains(sg.kinds, field.Desc.Kind()) {
			count++
			continue
		}

		if field.Desc.IsList() || sg.isList {
			continue
		}

		if field.Desc.IsMap() {
			continue
		}

		if field.Desc.IsExtension() {
			continue
		}

		if utils.Contains(sg.kinds, field.Desc.Kind()) {
			count++
			continue
		}

		if field.Desc.Kind() == protoreflect.MessageKind {
			count += sg.fieldsLen(field.Message.Fields)
			continue
		}
	}
	return count
}

// generateMessageGetters generates getter functions for fields of a message
func generateMessageGetters(g *protogen.GeneratedFile, msg *protogen.Message, generateResourceMethods bool) error {
	g.P(fmt.Sprintf("// %s_FieldGetters is a struct that contains getter functions for fields of %s", msg.GoIdent.GoName, msg.Desc.Name()))
	g.P(fmt.Sprintf("type %s_FieldGetters struct {", msg.GoIdent.GoName))
	g.P(fmt.Sprintf("StringGetter      func(msg protoreflect.ProtoMessage, path string) (string, error)"))
	g.P(fmt.Sprintf("StringListGetter  func(msg protoreflect.ProtoMessage, path string) ([]string, error)"))
	g.P(fmt.Sprintf("BoolGetter        func(msg protoreflect.ProtoMessage, path string) (bool, error)"))
	g.P(fmt.Sprintf("BoolListGetter    func(msg protoreflect.ProtoMessage, path string) ([]bool, error)"))
	g.P(fmt.Sprintf("IntGetter         func(msg protoreflect.ProtoMessage, path string) (int64, error)"))
	g.P(fmt.Sprintf("IntListGetter     func(msg protoreflect.ProtoMessage, path string) ([]int64, error)"))
	g.P(fmt.Sprintf("FloatGetter       func(msg protoreflect.ProtoMessage, path string) (float64, error)"))
	g.P(fmt.Sprintf("FloatListGetter   func(msg protoreflect.ProtoMessage, path string) ([]float64, error)"))
	g.P(fmt.Sprintf("EnumGetter        func(msg protoreflect.ProtoMessage, path string) (protoreflect.EnumNumber, error)"))
	g.P(fmt.Sprintf("EnumListGetter    func(msg protoreflect.ProtoMessage, path string) ([]protoreflect.EnumNumber, error)"))
	g.P(fmt.Sprintf("SubMessageGetter  func(msg protoreflect.ProtoMessage, path string) (protoreflect.ProtoMessage, error)"))
	g.P(fmt.Sprintf("SubMessageListGetter func(msg protoreflect.ProtoMessage, path string) ([]protoreflect.ProtoMessage, error)"))
	g.P("}")

	g.P()

	g.P(fmt.Sprintf("// New%sFieldGetters creates a new instance of %s_FieldGetters", msg.GoIdent.GoName, msg.GoIdent.GoName))
	g.P(fmt.Sprintf("func New%sFieldGetters() *%s_FieldGetters {", msg.GoIdent.GoName, msg.GoIdent.GoName))
	g.P(fmt.Sprintf("fieldGetters := &%s_FieldGetters{}", msg.GoIdent.GoName))
	g.P()
	// StringGetter
	{
		fieldsGetter := &FieldsGetter{
			gen:   g,
			msg:   msg,
			kinds: []protoreflect.Kind{protoreflect.StringKind},
		}

		fieldsGetter.gen.P("fieldGetters.StringGetter = func(msg protoreflect.ProtoMessage, path string) (string, error) {")
		if fieldsGetter.fieldsLen(msg.Fields) > 0 {
			fieldsGetter.gen.P(fmt.Sprintf("%sMsg := msg.(*%s)", strings.ToLower(msg.GoIdent.GoName), msg.GoIdent.GoName))
			fieldsGetter.gen.P()
		}
		fieldsGetter.gen.P("switch path {")
		err := fieldsGetter.processFields(msg.Fields)
		if err != nil {
			return err
		}
		fieldsGetter.gen.P("default:")
		fieldsGetter.gen.P("return \"\", nil")
		fieldsGetter.gen.P("}")
		fieldsGetter.gen.P("}")
		fieldsGetter.gen.P()
	}

	// StringListGetter
	{
		fieldsGetter := &FieldsGetter{
			gen:    g,
			msg:    msg,
			kinds:  []protoreflect.Kind{protoreflect.StringKind},
			isList: true,
		}

		fieldsGetter.gen.P("fieldGetters.StringListGetter = func(msg protoreflect.ProtoMessage, path string) ([]string, error) {")
		if fieldsGetter.fieldsLen(msg.Fields) > 0 {
			fieldsGetter.gen.P(fmt.Sprintf("%sMsg := msg.(*%s)", strings.ToLower(msg.GoIdent.GoName), msg.GoIdent.GoName))
			fieldsGetter.gen.P()
		}
		fieldsGetter.gen.P("switch path {")
		err := fieldsGetter.processFields(msg.Fields)
		if err != nil {
			return err
		}
		fieldsGetter.gen.P("default:")
		fieldsGetter.gen.P("return nil, nil")
		fieldsGetter.gen.P("}")
		fieldsGetter.gen.P("}")
		fieldsGetter.gen.P()
	}

	// BoolGetter
	{
		fieldsGetter := &FieldsGetter{
			gen:   g,
			msg:   msg,
			kinds: []protoreflect.Kind{protoreflect.BoolKind},
		}

		fieldsGetter.gen.P("fieldGetters.BoolGetter = func(msg protoreflect.ProtoMessage, path string) (bool, error) {")
		if fieldsGetter.fieldsLen(msg.Fields) > 0 {
			fieldsGetter.gen.P(fmt.Sprintf("%sMsg := msg.(*%s)", strings.ToLower(msg.GoIdent.GoName), msg.GoIdent.GoName))
			fieldsGetter.gen.P()
		}
		fieldsGetter.gen.P("switch path {")
		err := fieldsGetter.processFields(msg.Fields)
		if err != nil {
			return err
		}
		fieldsGetter.gen.P("default:")
		fieldsGetter.gen.P("return false, nil")
		fieldsGetter.gen.P("}")
		fieldsGetter.gen.P("}")
		fieldsGetter.gen.P()
	}

	// BoolListGetter
	{
		fieldsGetter := &FieldsGetter{
			gen:    g,
			msg:    msg,
			kinds:  []protoreflect.Kind{protoreflect.BoolKind},
			isList: true,
		}

		fieldsGetter.gen.P("fieldGetters.BoolListGetter = func(msg protoreflect.ProtoMessage, path string) ([]bool, error) {")
		if fieldsGetter.fieldsLen(msg.Fields) > 0 {
			fieldsGetter.gen.P(fmt.Sprintf("%sMsg := msg.(*%s)", strings.ToLower(msg.GoIdent.GoName), msg.GoIdent.GoName))
			fieldsGetter.gen.P()
		}
		fieldsGetter.gen.P("switch path {")
		err := fieldsGetter.processFields(msg.Fields)
		if err != nil {
			return err
		}
		fieldsGetter.gen.P("default:")
		fieldsGetter.gen.P("return nil, nil")
		fieldsGetter.gen.P("}")
		fieldsGetter.gen.P("}")
		fieldsGetter.gen.P()
	}

	// IntGetter
	{
		fieldsGetter := &FieldsGetter{
			gen:   g,
			msg:   msg,
			kinds: []protoreflect.Kind{protoreflect.Int32Kind, protoreflect.Int64Kind},
		}

		fieldsGetter.gen.P("fieldGetters.IntGetter = func(msg protoreflect.ProtoMessage, path string) (int64, error) {")
		if fieldsGetter.fieldsLen(msg.Fields) > 0 {
			fieldsGetter.gen.P(fmt.Sprintf("%sMsg := msg.(*%s)", strings.ToLower(msg.GoIdent.GoName), msg.GoIdent.GoName))
			fieldsGetter.gen.P()
		}
		fieldsGetter.gen.P("switch path {")
		err := fieldsGetter.processFields(msg.Fields)
		if err != nil {
			return err
		}
		fieldsGetter.gen.P("default:")
		fieldsGetter.gen.P("return 0, nil")
		fieldsGetter.gen.P("}")
		fieldsGetter.gen.P("}")
		fieldsGetter.gen.P()
	}

	// IntListGetter
	{
		fieldsGetter := &FieldsGetter{
			gen:    g,
			msg:    msg,
			kinds:  []protoreflect.Kind{protoreflect.Int32Kind, protoreflect.Int64Kind},
			isList: true,
		}

		fieldsGetter.gen.P("fieldGetters.IntListGetter = func(msg protoreflect.ProtoMessage, path string) ([]int64, error) {")
		if fieldsGetter.fieldsLen(msg.Fields) > 0 {
			fieldsGetter.gen.P(fmt.Sprintf("%sMsg := msg.(*%s)", strings.ToLower(msg.GoIdent.GoName), msg.GoIdent.GoName))
			fieldsGetter.gen.P()
		}
		fieldsGetter.gen.P("switch path {")
		err := fieldsGetter.processFields(msg.Fields)
		if err != nil {
			return err
		}
		fieldsGetter.gen.P("default:")
		fieldsGetter.gen.P("return nil, nil")
		fieldsGetter.gen.P("}")
		fieldsGetter.gen.P("}")
		fieldsGetter.gen.P()
	}

	// FloatGetter
	{
		fieldsGetter := &FieldsGetter{
			gen:   g,
			msg:   msg,
			kinds: []protoreflect.Kind{protoreflect.FloatKind, protoreflect.DoubleKind},
		}

		fieldsGetter.gen.P("fieldGetters.FloatGetter = func(msg protoreflect.ProtoMessage, path string) (float64, error) {")
		if fieldsGetter.fieldsLen(msg.Fields) > 0 {
			fieldsGetter.gen.P(fmt.Sprintf("%sMsg := msg.(*%s)", strings.ToLower(msg.GoIdent.GoName), msg.GoIdent.GoName))
			fieldsGetter.gen.P()
		}
		fieldsGetter.gen.P("switch path {")
		err := fieldsGetter.processFields(msg.Fields)
		if err != nil {
			return err
		}
		fieldsGetter.gen.P("default:")
		fieldsGetter.gen.P("return 0.0, nil")
		fieldsGetter.gen.P("}")
		fieldsGetter.gen.P("}")
		fieldsGetter.gen.P()
	}

	// FloatGetter
	{
		fieldsGetter := &FieldsGetter{
			gen:    g,
			msg:    msg,
			kinds:  []protoreflect.Kind{protoreflect.FloatKind, protoreflect.DoubleKind},
			isList: true,
		}

		fieldsGetter.gen.P("fieldGetters.FloatListGetter = func(msg protoreflect.ProtoMessage, path string) ([]float64, error) {")
		if fieldsGetter.fieldsLen(msg.Fields) > 0 {
			fieldsGetter.gen.P(fmt.Sprintf("%sMsg := msg.(*%s)", strings.ToLower(msg.GoIdent.GoName), msg.GoIdent.GoName))
			fieldsGetter.gen.P()
		}
		fieldsGetter.gen.P("switch path {")
		err := fieldsGetter.processFields(msg.Fields)
		if err != nil {
			return err
		}
		fieldsGetter.gen.P("default:")
		fieldsGetter.gen.P("return nil, nil")
		fieldsGetter.gen.P("}")
		fieldsGetter.gen.P("}")
		fieldsGetter.gen.P()
	}

	// EnumGetter
	{
		fieldsGetter := &FieldsGetter{
			gen:   g,
			msg:   msg,
			kinds: []protoreflect.Kind{protoreflect.EnumKind},
		}

		fieldsGetter.gen.P("fieldGetters.EnumGetter = func(msg protoreflect.ProtoMessage, path string) (protoreflect.EnumNumber, error) {")
		if fieldsGetter.fieldsLen(msg.Fields) > 0 {
			fieldsGetter.gen.P(fmt.Sprintf("%sMsg := msg.(*%s)", strings.ToLower(msg.GoIdent.GoName), msg.GoIdent.GoName))
			fieldsGetter.gen.P()
		}
		fieldsGetter.gen.P("switch path {")
		err := fieldsGetter.processFields(msg.Fields)
		if err != nil {
			return err
		}
		fieldsGetter.gen.P("default:")
		fieldsGetter.gen.P("return 0, nil")
		fieldsGetter.gen.P("}")
		fieldsGetter.gen.P("}")
		fieldsGetter.gen.P()
	}

	// EnumListGetter
	{
		fieldsGetter := &FieldsGetter{
			gen:    g,
			msg:    msg,
			kinds:  []protoreflect.Kind{protoreflect.EnumKind},
			isList: true,
		}

		fieldsGetter.gen.P("fieldGetters.EnumListGetter = func(msg protoreflect.ProtoMessage, path string) ([]protoreflect.EnumNumber, error) {")
		if fieldsGetter.fieldsLen(msg.Fields) > 0 {
			fieldsGetter.gen.P(fmt.Sprintf("%sMsg := msg.(*%s)", strings.ToLower(msg.GoIdent.GoName), msg.GoIdent.GoName))
			fieldsGetter.gen.P()
		}
		fieldsGetter.gen.P("switch path {")
		err := fieldsGetter.processFields(msg.Fields)
		if err != nil {
			return err
		}
		fieldsGetter.gen.P("default:")
		fieldsGetter.gen.P("return nil, nil")
		fieldsGetter.gen.P("}")
		fieldsGetter.gen.P("}")
		fieldsGetter.gen.P()
	}

	// SubMessageGetter
	{
		fieldsGetter := &FieldsGetter{
			gen:   g,
			msg:   msg,
			kinds: []protoreflect.Kind{protoreflect.MessageKind},
		}

		fieldsGetter.gen.P("fieldGetters.SubMessageGetter = func(msg protoreflect.ProtoMessage, path string) (protoreflect.ProtoMessage, error) {")
		if fieldsGetter.fieldsLen(msg.Fields) > 0 {
			fieldsGetter.gen.P(fmt.Sprintf("%sMsg := msg.(*%s)", strings.ToLower(msg.GoIdent.GoName), msg.GoIdent.GoName))
			fieldsGetter.gen.P()
		}
		fieldsGetter.gen.P("switch path {")
		err := fieldsGetter.processFields(msg.Fields)
		if err != nil {
			return err
		}
		fieldsGetter.gen.P("default:")
		fieldsGetter.gen.P("return nil, nil")
		fieldsGetter.gen.P("}")
		fieldsGetter.gen.P("}")
		fieldsGetter.gen.P()
	}

	// SubMessageListGetter
	{
		fieldsGetter := &FieldsGetter{
			gen:    g,
			msg:    msg,
			kinds:  []protoreflect.Kind{protoreflect.MessageKind},
			isList: true,
		}

		fieldsGetter.gen.P("fieldGetters.SubMessageListGetter = func(msg protoreflect.ProtoMessage, path string) ([]protoreflect.ProtoMessage, error) {")
		if fieldsGetter.fieldsLen(msg.Fields) > 0 {
			fieldsGetter.gen.P(fmt.Sprintf("%sMsg := msg.(*%s)", strings.ToLower(msg.GoIdent.GoName), msg.GoIdent.GoName))
			fieldsGetter.gen.P()
		}
		fieldsGetter.gen.P("switch path {")
		err := fieldsGetter.processFields(msg.Fields)
		if err != nil {
			return err
		}
		fieldsGetter.gen.P("default:")
		fieldsGetter.gen.P("return nil, nil")
		fieldsGetter.gen.P("}")
		fieldsGetter.gen.P("}")
		fieldsGetter.gen.P()
	}

	g.P("return fieldGetters")
	g.P("}")

	if generateResourceMethods {
		// StringGetter
		{
			g.P(fmt.Sprintf("// StringGetter is a getter function for string fields of %s", msg.Desc.Name()))
			g.P(fmt.Sprintf("func (x *%s) StringGetter(path string) (string, error) {", msg.GoIdent.GoName))
			g.P(fmt.Sprintf("fieldGetters := New%sFieldGetters()", msg.GoIdent.GoName))
			g.P("return fieldGetters.StringGetter(x, path)")
			g.P("}")
			g.P()
		}

		// StringListGetter
		{

			g.P(fmt.Sprintf("// StringListGetter is a getter function for string list fields of %s", msg.Desc.Name()))
			g.P(fmt.Sprintf("func (x *%s) StringListGetter(path string) ([]string, error) {", msg.GoIdent.GoName))
			g.P(fmt.Sprintf("fieldGetters := New%sFieldGetters()", msg.GoIdent.GoName))
			g.P("return fieldGetters.StringListGetter(x, path)")
			g.P("}")
			g.P()
		}

		// BoolGetter
		{

			g.P(fmt.Sprintf("// BoolGetter is a getter function for bool fields of %s", msg.Desc.Name()))
			g.P(fmt.Sprintf("func (x *%s) BoolGetter(path string) (bool, error) {", msg.GoIdent.GoName))
			g.P(fmt.Sprintf("fieldGetters := New%sFieldGetters()", msg.GoIdent.GoName))
			g.P("return fieldGetters.BoolGetter(x, path)")
			g.P("}")
			g.P()
		}

		// BoolListGetter
		{

			g.P(fmt.Sprintf("// BoolListGetter is a getter function for bool list fields of %s", msg.Desc.Name()))
			g.P(fmt.Sprintf("func (x *%s) BoolListGetter(path string) ([]bool, error) {", msg.GoIdent.GoName))
			g.P(fmt.Sprintf("fieldGetters := New%sFieldGetters()", msg.GoIdent.GoName))
			g.P("return fieldGetters.BoolListGetter(x, path)")
			g.P("}")
			g.P()
		}

		// IntGetter
		{

			g.P(fmt.Sprintf("// IntGetter is a getter function for int32 and int64 fields of %s", msg.Desc.Name()))
			g.P(fmt.Sprintf("func (x *%s) IntGetter(path string) (int64, error) {", msg.GoIdent.GoName))
			g.P(fmt.Sprintf("fieldGetters := New%sFieldGetters()", msg.GoIdent.GoName))
			g.P("return fieldGetters.IntGetter(x, path)")
			g.P("}")
			g.P()
		}

		// IntListGetter
		{

			g.P(fmt.Sprintf("// IntListGetter is a getter function for int32 and int64 list fields of %s", msg.Desc.Name()))
			g.P(fmt.Sprintf("func (x *%s) IntListGetter(path string) ([]int64, error) {", msg.GoIdent.GoName))
			g.P(fmt.Sprintf("fieldGetters := New%sFieldGetters()", msg.GoIdent.GoName))
			g.P("return fieldGetters.IntListGetter(x, path)")
			g.P("}")
			g.P()
		}

		// FloatGetter
		{

			g.P(fmt.Sprintf("// FloatGetter is a getter function for float32 and float64 fields of %s", msg.Desc.Name()))
			g.P(fmt.Sprintf("func (x *%s) FloatGetter(path string) (float64, error) {", msg.GoIdent.GoName))
			g.P(fmt.Sprintf("fieldGetters := New%sFieldGetters()", msg.GoIdent.GoName))
			g.P("return fieldGetters.FloatGetter(x, path)")
			g.P("}")
			g.P()
		}

		// FloatListGetter
		{

			g.P(fmt.Sprintf("// FloatListGetter is a getter function for float32 and float64 list fields of %s", msg.Desc.Name()))
			g.P(fmt.Sprintf("func (x *%s) FloatListGetter(path string) ([]float64, error) {", msg.GoIdent.GoName))
			g.P(fmt.Sprintf("fieldGetters := New%sFieldGetters()", msg.GoIdent.GoName))
			g.P("return fieldGetters.FloatListGetter(x, path)")
			g.P("}")
			g.P()
		}

		// EnumGetter
		{

			g.P(fmt.Sprintf("// EnumGetter is a getter function for enum fields of %s", msg.Desc.Name()))
			g.P(fmt.Sprintf("func (x *%s) EnumGetter(path string) (protoreflect.EnumNumber, error) {", msg.GoIdent.GoName))
			g.P(fmt.Sprintf("fieldGetters := New%sFieldGetters()", msg.GoIdent.GoName))
			g.P("return fieldGetters.EnumGetter(x, path)")
			g.P("}")
			g.P()
		}

		// EnumListGetter
		{

			g.P(fmt.Sprintf("// EnumListGetter is a getter function for enum list fields of %s", msg.Desc.Name()))
			g.P(fmt.Sprintf("func (x *%s) EnumListGetter(path string) ([]protoreflect.EnumNumber, error) {", msg.GoIdent.GoName))
			g.P(fmt.Sprintf("fieldGetters := New%sFieldGetters()", msg.GoIdent.GoName))
			g.P("return fieldGetters.EnumListGetter(x, path)")
			g.P("}")
			g.P()
		}

		// SubMessageGetter
		{

			g.P(fmt.Sprintf("// SubMessageGetter is a getter function for submessage fields of %s", msg.Desc.Name()))
			g.P(fmt.Sprintf("func (x *%s) SubMessageGetter(path string) (protoreflect.ProtoMessage, error) {", msg.GoIdent.GoName))
			g.P(fmt.Sprintf("fieldGetters := New%sFieldGetters()", msg.GoIdent.GoName))
			g.P("return fieldGetters.SubMessageGetter(x, path)")
			g.P("}")
			g.P()
		}

		// SubMessageListGetter
		{

			g.P(fmt.Sprintf("// SubMessageListGetter is a getter function for submessage list fields of %s", msg.Desc.Name()))
			g.P(fmt.Sprintf("func (x *%s) SubMessageListGetter(path string) ([]protoreflect.ProtoMessage, error) {", msg.GoIdent.GoName))
			g.P(fmt.Sprintf("fieldGetters := New%sFieldGetters()", msg.GoIdent.GoName))
			g.P("return fieldGetters.SubMessageListGetter(x, path)")
			g.P("}")
			g.P()
		}
	}

	return nil
}
