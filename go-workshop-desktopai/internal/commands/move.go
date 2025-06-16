package commands

import (
	"fmt"
	"github.com/spf13/cobra"
)

func NewMoveCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "move [x] [y]",
		Short: "Move the mouse cursor to coordinates",
		Long:  "Move the mouse cursor to the specified (x, y) screen coordinate.",
		Example: `
  desktop-automation move 300 400
  desktop-automation move 800 600
`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement move logic
			fmt.Println("Move command is not yet implemented.")
			return nil
		},
	}
	return cmd
}
