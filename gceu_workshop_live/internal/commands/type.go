package commands

import (
	"fmt"
	"strings"

	"github.com/dmahlow/desktop-automation/internal/automation"
	"github.com/spf13/cobra"
)

// NewTypeCommand creates the type command
func NewTypeCommand() *cobra.Command {
	var delayMs int

	cmd := &cobra.Command{
		Use:   "type <text>",
		Short: "Type text at current cursor position",
		Long: `Type text at the current cursor position.

This command simulates keyboard input by typing the specified text at the current
cursor location. The text will be typed as if you were physically typing on the keyboard.

Use quotes to handle multi-word text or text containing special characters.`,
		Example: `  # Type a simple message
  desktop-automation type "Hello, World!"

  # Type a sentence with spaces
  desktop-automation type "This is a test message"

  # Type with delay between characters (50ms)
  desktop-automation type --delay=50 "Slow typing!"

  # Type special characters
  desktop-automation type "user@example.com"

  # Type numbers and symbols
  desktop-automation type "Password123!"`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTypeCommand(cmd, args, delayMs)
		},
	}

	// Add delay flag
	cmd.Flags().IntVar(&delayMs, "delay", 0, "Delay in milliseconds between each character (default: 0)")

	return cmd
}

// runTypeCommand handles the type command execution
func runTypeCommand(cmd *cobra.Command, args []string, delayMs int) error {
	text := args[0]

	// Validate that text is not empty or only whitespace
	if strings.TrimSpace(text) == "" {
		return fmt.Errorf("text cannot be empty or contain only whitespace")
	}

	// Show what we're about to type
	fmt.Printf("Typing text: %q", text)
	if delayMs > 0 {
		fmt.Printf(" (with %dms delay between characters)", delayMs)
	}
	fmt.Println()

	// Use appropriate typing function based on delay
	var err error
	if delayMs > 0 {
		err = automation.TypeStringWithDelay(text, delayMs)
	} else {
		err = automation.TypeString(text)
	}

	if err != nil {
		return fmt.Errorf("failed to type text: %w", err)
	}

	// Show success message with character count
	charCount := len(text)
	fmt.Printf("✓ Successfully typed %d character", charCount)
	if charCount != 1 {
		fmt.Print("s")
	}
	fmt.Println()

	return nil
}
