package main

import (
	"context"
	"log"

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

	b := entservice.NewBuilder(
		entservice.WithModule("github.com/syralon/entc-gen-go/example"),
		entservice.WithOutput("example"),
	)
	if err = b.Generate(ctx, graph); err != nil {
		log.Fatal(err)
	}
}
