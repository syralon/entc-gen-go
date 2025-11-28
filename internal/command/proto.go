package command

import (
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"github.com/spf13/cobra"
	"github.com/syralon/entc-gen-go/internal/generate/entproto"
	"os/exec"
)

func Proto() *cobra.Command {
	var target string
	var output string
	cmd := &cobra.Command{
		Use:   "proto",
		Short: "Generate protobuf and grpc proto files from ent schemas.",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			cfg := &gen.Config{}
			graph, err := entc.LoadGraph(target, cfg)
			if err != nil {
				return err
			}
			generator := entproto.New(output, "")
			if err = generator.Generate(ctx, graph); err != nil {
				return err
			}
			_ = exec.Command("buf", "format", "-w").Run()

			return nil
		},
	}
	cmd.PersistentFlags().StringVar(&target, "target", "./ent/schema", "ent target")
	cmd.PersistentFlags().StringVarP(&output, "output", "o", "", "proto files output directory")
	return cmd
}
