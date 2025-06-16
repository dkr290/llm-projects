package commands

import (
	"github.com/spf13/cobra"
)

// AddCommands adds all subcommands to the root command
func AddCommands(rootCmd *cobra.Command) {
	rootCmd.AddCommand(
		NewClickCommand(),
		NewTypeCommand(),
		NewMoveCommand(),
	)
}
