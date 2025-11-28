package command

import (
	"os/exec"

	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	return &cobra.Command{
		Use:   "new",
		Short: "Create new entities, same as 'ent new'.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return exec.Command("ent", append([]string{"new"}, args...)...).Run()
		},
	}
}
