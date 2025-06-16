package commands

import (
	"fmt"
	"strconv"

	"github.com/dmahlow/desktop-automation/internal/automation"
	"github.com/spf13/cobra"
)

// NewClickCommand creates the click command
func NewClickCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "click <x> <y>",
		Short: "Click at a specific screen coordinate",
		Long: `Click at a specific screen coordinate.

This command simulates a mouse click at the specified x and y coordinates on the screen.
The coordinates are measured in pixels from the top-left corner of the screen (0,0).`,
		Example: `  # Click at coordinates (100, 200)
  desktop-automation click 100 200

  # Click at the center of a 1920x1080 screen
  desktop-automation click 960 540

  # Click at the top-left corner
  desktop-automation click 0 0`,
		Args: cobra.ExactArgs(2),
		RunE: runClickCommand,
	}

	return cmd
}

// runClickCommand handles the click command execution
func runClickCommand(cmd *cobra.Command, args []string) error {
	// Parse X coordinate
	x, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid x coordinate '%s': must be a valid integer", args[0])
	}

	// Parse Y coordinate
	y, err := strconv.Atoi(args[1])
	if err != nil {
		return fmt.Errorf("invalid y coordinate '%s': must be a valid integer", args[1])
	}

	// Validate coordinates are non-negative
	if x < 0 {
		return fmt.Errorf("x coordinate cannot be negative: %d", x)
	}
	if y < 0 {
		return fmt.Errorf("y coordinate cannot be negative: %d", y)
	}

	// Show current mouse position before click
	currentX, currentY := automation.GetPosition()
	fmt.Printf("Current mouse position: (%d, %d)\n", currentX, currentY)

	// Perform the click using our automation
	fmt.Printf("Clicking at coordinates (%d, %d)...\n", x, y)

	err = automation.Click(x, y)
	if err != nil {
		return fmt.Errorf("failed to click at (%d, %d): %w", x, y, err)
	}

	// Confirm success with coordinates
	fmt.Printf("✓ Successfully clicked at coordinates (%d, %d)\n", x, y)

	return nil
}
