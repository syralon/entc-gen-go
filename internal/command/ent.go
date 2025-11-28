package command

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func Ent() []*cobra.Command {
	var commands = []string{"describe", "generate", "new", "schema"}
	var cmds = make([]*cobra.Command, 0, len(commands))
	for _, cmd := range commands {
		cmds = append(cmds, ent(cmd))
	}
	return cmds
}

func ent(command string) *cobra.Command {
	return &cobra.Command{
		Use:   command,
		Short: fmt.Sprintf("same as 'ent %s'", command),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				args = append(args, "--help")
			}
			cc := exec.Command("ent", append([]string{command}, args...)...)
			cc.Stdout = os.Stdout
			cc.Stderr = os.Stderr
			return cc.Run()
		},
	}
}
