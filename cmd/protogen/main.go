// Command protogen generates Go code from Protobuf contracts using Go Modules.
package main

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/bufbuild/protocompile"
	gogoproto "github.com/gogo/protobuf/proto"
	gogodescriptor "github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	gogoplugin "github.com/gogo/protobuf/protoc-gen-gogo/plugin"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

const gogoImportPrefix = "github.com/gogo/protobuf/"

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string) error {
	flags := flag.NewFlagSet("protogen", flag.ContinueOnError)
	importPath := flags.String("import-path", "", "Go import path for generated messages")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 1 {
		return errors.New("expected exactly one .proto source file")
	}
	source := flags.Arg(0)
	if filepath.Base(source) != source || filepath.Ext(source) != ".proto" {
		return fmt.Errorf("invalid .proto source file %q", source)
	}

	modulePath, err := exec.Command("go", "list", "-m", "-f", "{{.Dir}}", "github.com/gogo/protobuf").Output()
	if err != nil {
		return fmt.Errorf("locate gogo/protobuf Go module: %w", err)
	}
	gogoDir := strings.TrimSpace(string(modulePath))
	if !filepath.IsAbs(gogoDir) {
		return fmt.Errorf("invalid gogo/protobuf module directory %q", gogoDir)
	}

	compiler := protocompile.Compiler{
		Resolver: protocompile.WithStandardImports(&protocompile.SourceResolver{
			Accessor: func(name string) (io.ReadCloser, error) {
				if strings.HasPrefix(name, gogoImportPrefix) {
					name = filepath.Join(gogoDir, strings.TrimPrefix(name, gogoImportPrefix))
				}
				return os.Open(name)
			},
		}),
		SourceInfoMode: protocompile.SourceInfoStandard,
	}
	files, err := compiler.Compile(context.Background(), source)
	if err != nil {
		return fmt.Errorf("compile %s: %w", source, err)
	}

	var descriptors []*descriptorpb.FileDescriptorProto
	visited := make(map[string]bool)
	var visit func(protoreflect.FileDescriptor)
	visit = func(file protoreflect.FileDescriptor) {
		if visited[file.Path()] {
			return
		}
		visited[file.Path()] = true
		for i := 0; i < file.Imports().Len(); i++ {
			visit(file.Imports().Get(i).FileDescriptor)
		}
		descriptors = append(descriptors, protodesc.ToFileDescriptorProto(file))
	}
	for _, file := range files {
		visit(file)
	}
	legacyDescriptors := make([]*gogodescriptor.FileDescriptorProto, 0, len(descriptors))
	for _, descriptor := range descriptors {
		data, err := proto.Marshal(descriptor)
		if err != nil {
			return fmt.Errorf("marshal descriptor: %w", err)
		}
		converted := new(gogodescriptor.FileDescriptorProto)
		if err := gogoproto.Unmarshal(data, converted); err != nil {
			return fmt.Errorf("convert descriptor for gogofaster: %w", err)
		}
		legacyDescriptors = append(legacyDescriptors, converted)
	}

	request := &gogoplugin.CodeGeneratorRequest{
		FileToGenerate: []string{source},
		ProtoFile:      legacyDescriptors,
	}
	if *importPath != "" {
		parameter := "import_path=" + *importPath
		request.Parameter = &parameter
	}
	payload, err := gogoproto.Marshal(request)
	if err != nil {
		return fmt.Errorf("encode generator request: %w", err)
	}
	command := exec.Command("go", "tool", "protoc-gen-gogofaster")
	command.Stdin = bytes.NewReader(payload)
	output, err := command.Output()
	if err != nil {
		return fmt.Errorf("run Go module generator: %w", err)
	}
	response := new(gogoplugin.CodeGeneratorResponse)
	if err := gogoproto.Unmarshal(output, response); err != nil {
		return fmt.Errorf("decode generator response: %w", err)
	}
	if response.GetError() != "" {
		return fmt.Errorf("generate %s: %s", source, response.GetError())
	}
	if len(response.File) == 0 {
		return fmt.Errorf("generator produced no files for %s", source)
	}
	for _, file := range response.File {
		name := filepath.Clean(file.GetName())
		if filepath.IsAbs(name) || name == "." || name == ".." || strings.HasPrefix(name, ".."+string(filepath.Separator)) {
			return fmt.Errorf("invalid generated file name %q", name)
		}
		if err := os.MkdirAll(filepath.Dir(name), 0755); err != nil {
			return fmt.Errorf("prepare generated output %s: %w", name, err)
		}
		if err := os.WriteFile(name, []byte(file.GetContent()), 0644); err != nil {
			return fmt.Errorf("write generated output %s: %w", name, err)
		}
	}
	return nil
}
