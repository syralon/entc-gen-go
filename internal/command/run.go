package command

import (
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"github.com/spf13/cobra"
	"github.com/syralon/entc-gen-go/internal/generate/entproto"
	"github.com/syralon/entc-gen-go/internal/generate/entservice"
	"os/exec"
)

func Run() *cobra.Command {
	var opt options
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Generate api services from ent schemas.",
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

			b := entservice.NewBuilder(
				entservice.WithModule(opt.module),
				entservice.WithOutput(opt.output),
				entservice.WithOverwrite(opt.overwrite),
			)
			_ = exec.Command("buf", "format", "-w").Run()
			_ = exec.Command("ent", "generate", "--target", opt.target).Run()
			_ = exec.Command("go", "generate", "./...", opt.target).Run()
			return b.Generate(cmd.Context(), graph)
		},
	}
	opt.register(cmd)
	return cmd
}
