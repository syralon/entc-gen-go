package command

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/syralon/entc-gen-go/internal/tools/text"
)

type options struct {
	target    string
	output    string
	module    string
	overwrite bool
}

func (o *options) register(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&o.target, "target", "./ent/schema", "The ent target directory.")
	cmd.PersistentFlags().StringVarP(&o.output, "output", "o", ".", "The output directory.")
	cmd.PersistentFlags().BoolVar(&o.overwrite, "overwrite", false, "overwrite generated files.")
}

func (o *options) parse() error {
	var err error
	o.module, err = text.Module(".")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("%w, use 'go mod init' to create a new mod", err)
		}
		return fmt.Errorf("parse mod file on error: %w", err)
	}
	return nil
}
