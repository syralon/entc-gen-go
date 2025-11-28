package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/syralon/entc-gen-go/internal/command"
	"os"
)

func main() {
	cmd := &cobra.Command{
		Use:           "entc-gen",
		Long:          "A service generator base on ent(https://entgo.io/).\nSee also: https://github.com/syralon/entc-gen-go",
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	cmd.AddCommand(
		command.New(),
		command.Proto(),
		command.Service(),
	)
	if err := cmd.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s", err.Error())
	}
}
