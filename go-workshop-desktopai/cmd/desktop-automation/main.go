package main

import (
	"go-workshop-desktopai/internal/commands"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "desktop-automation",
		Short: "Beautiful Desktop Automation CLI",
		Long:  "Beautiful Desktop Automation CLI for scripting mouse and keyboard actions.",
	}

	rootCmd.AddCommand(
		commands.NewClickCommand(),

		commands.NewTypeCommand(),
		commands.NewMoveCommand(),
	)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
