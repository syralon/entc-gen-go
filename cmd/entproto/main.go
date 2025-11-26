package main

import (
	"context"
	"log"
	"os/exec"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"github.com/syralon/entc-gen-go/internal/generate/entproto"
)

func main() {
	cfg := &gen.Config{}
	graph, err := entc.LoadGraph("./example/ent/schema", cfg)
	if err != nil {
		log.Fatal(err)
	}

	generator := entproto.New(
		"example",
		"proto/example",
		"",
	)
	err = generator.Generate(context.Background(), graph)
	if err != nil {
		log.Fatal(err)
	}
	_ = exec.Command("buf", "format", "-w").Run()
}
