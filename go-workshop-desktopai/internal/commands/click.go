package commands

import (
	"fmt"
	"github.com/spf13/cobra"
)

func NewClickCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "click [x] [y]",
		Short: "Click at a specific screen coordinate",
		Long:  "Click at the given (x, y) screen coordinate with the mouse.",
		Example: `
  desktop-automation click 100 200
  desktop-automation click 500 600
`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement click logic
			fmt.Println("Click command is not yet implemented.")
			return nil
		},
	}
	return cmd
}
