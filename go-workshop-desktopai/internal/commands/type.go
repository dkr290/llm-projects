package commands

import (
	"fmt"
	"github.com/spf13/cobra"
)

func NewTypeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "type [text]",
		Short: "Type text at the current cursor position",
		Long:  "Type the provided text string at the current mouse cursor position.",
		Example: `
  desktop-automation type "Hello world!"
  desktop-automation type "Automate all the things."
`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement type logic
			fmt.Println("Type command is not yet implemented.")
			return nil
		},
	}
	return cmd
}
