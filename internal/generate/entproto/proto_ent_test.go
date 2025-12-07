package entproto

import (
	"context"
	"fmt"
	"testing"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"github.com/jhump/protoreflect/v2/protoprint"
)

func TestEntBuilder(t *testing.T) {
	cfg := &gen.Config{}
	graph, err := entc.LoadGraph("./testdata/ent/schema", cfg)
	if err != nil {
		t.Fatal(err)
	}
	eb := NewEntBuilder(WithProtoPackage("example"), WithGoPackage("github.com/syralon/entc-gen-go/proto/syralon/example;example"))

	ctx := NewContext(context.Background())

	files, err := eb.Build(ctx, graph)
	if err != nil {
		t.Fatal(err)
	}

	printer := &protoprint.Printer{}
	for _, file := range files {
		descriptor, err := file.Build()
		if err != nil {
			t.Error(err)
			return
		}
		text, err := printer.PrintProtoToString(descriptor)
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Println(text)
		fmt.Println()
	}
}

func TestServiceBuilder(t *testing.T) {
	cfg := &gen.Config{}
	graph, err := entc.LoadGraph("./testdata/ent/schema", cfg)
	if err != nil {
		t.Fatal(err)
	}

	sb := NewServiceBuilder(WithProtoPackage("example"), WithGoPackage("github.com/syralon/entc-gen-go/proto/syralon/example;example"))
	eb := NewEntBuilder(WithProtoPackage("example"), WithGoPackage("github.com/syralon/entc-gen-go/proto/syralon/example;example"))

	ctx := NewContext(context.Background())

	printer := &protoprint.Printer{}

	if _, err = eb.Build(ctx, graph); err != nil {
		t.Fatal(err)
		return
	}
	files, err := sb.Build(ctx, graph)
	if err != nil {
		t.Fatal(err)
		return
	}

	for _, message := range ctx.messages {
		if message.ParentFile() == nil {
			t.Error(message.Name())
		}
	}

	for _, file := range files {
		descriptor, err := file.Build()
		if err != nil {
			t.Error(err)
			return
		}
		text, err := printer.PrintProtoToString(descriptor)
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Println("=============================")
		fmt.Println(file.Path())
		fmt.Println("=============================")
		fmt.Println(text)
		fmt.Println()
	}
}
