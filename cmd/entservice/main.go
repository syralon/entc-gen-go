package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"github.com/syralon/entc-gen-go/internal/generate/entservice"
)

func main() {
	cfg := &gen.Config{}
	graph, err := entc.LoadGraph("./example/ent/schema", cfg)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	if err = build(ctx, nil, &entservice.HelperBuilder{}, "example/internal/service/helper.go"); err != nil {
		fmt.Println(err)
	}
	s := entservice.NewServiceBuilder(
		entservice.WithEntPackage("github.com/syralon/entc-gen-go/example/ent"),
		entservice.WithProtoPackage("github.com/syralon/entc-gen-go/example/proto/example"),
	)
	for _, node := range graph.Nodes {
		name := strings.ToLower(node.Name)
		filename := fmt.Sprintf("example/internal/service/%s.go", name)
		if err = build(ctx, node, s, filename); err != nil {
			fmt.Println(err)
		}
	}
}

func build(ctx context.Context, node *gen.Type, b entservice.Builder, filename string) error {
	file, err := b.Build(ctx, node)
	if err != nil {
		return err
	}
	fmt.Println(filename)
	_ = os.MkdirAll(path.Dir(filename), 0700)
	return file.Save(filename)
}
