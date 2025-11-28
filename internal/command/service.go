package command

import (
	"path"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"github.com/spf13/cobra"
	"github.com/syralon/entc-gen-go/internal/generate/entservice"
)

func Service() *cobra.Command {
	var opt options
	cmd := &cobra.Command{
		Use:   "service",
		Short: "Generate api services from ent schemas.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opt.parse(); err != nil {
				return err
			}
			cfg := &gen.Config{}
			graph, err := entc.LoadGraph(opt.target, cfg)
			if err != nil {
				return err
			}
			module := path.Join(opt.module, opt.output)
			b := entservice.NewBuilder(
				entservice.WithModule(module),
				entservice.WithOutput(opt.output),
			)
			return b.Generate(cmd.Context(), graph)
		},
	}
	opt.register(cmd)
	return cmd
}
