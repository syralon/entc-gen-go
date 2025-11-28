package command

import (
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"errors"
	"github.com/spf13/cobra"
	"github.com/syralon/entc-gen-go/internal/generate/entservice"
)

func Service() *cobra.Command {
	var target string
	var module string
	var output string
	cmd := &cobra.Command{
		Use:   "service",
		Short: "Generate api services from ent schemas.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if module == "" {
				return errors.New("module name is required")
			}
			cfg := &gen.Config{}
			graph, err := entc.LoadGraph(target, cfg)
			if err != nil {
				return err
			}

			b := entservice.NewBuilder(
				entservice.WithModule(module),
				entservice.WithOutput(output),
			)
			return b.Generate(cmd.Context(), graph)
		},
	}
	cmd.PersistentFlags().StringVar(&target, "target", "./ent/schema", "ent target")
	cmd.PersistentFlags().StringVarP(&output, "output", "o", "", "proto files output directory")
	cmd.PersistentFlags().StringVarP(&module, "module", "m", "", "module name")
	return cmd
}
