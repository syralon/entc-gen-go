package command

import (
	"os/exec"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"github.com/spf13/cobra"
	"github.com/syralon/entc-gen-go/internal/generate/entproto"
)

func Proto() *cobra.Command {
	var opt options
	cmd := &cobra.Command{
		Use:   "proto",
		Short: "Generate protobuf and grpc proto files from ent schemas.",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			if err := opt.parse(); err != nil {
				return err
			}
			cfg := &gen.Config{}
			graph, err := entc.LoadGraph(opt.target, cfg)
			if err != nil {
				return err
			}
			generator := entproto.New(opt.module, opt.output)
			if err = generator.Generate(ctx, graph); err != nil {
				return err
			}
			_ = exec.Command("buf", "format", "-w").Run()
			return nil
		},
	}
	opt.register(cmd)
	return cmd
}
