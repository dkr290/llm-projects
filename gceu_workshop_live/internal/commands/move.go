package commands

import (
	"fmt"
	"strconv"
	"time"

	"github.com/dmahlow/desktop-automation/internal/automation"
	"github.com/spf13/cobra"
)

// NewMoveCommand creates the move command
func NewMoveCommand() *cobra.Command {
	var smooth bool
	var duration float64

	cmd := &cobra.Command{
		Use:   "move <x> <y>",
		Short: "Move the mouse cursor to coordinates",
		Long: `Move the mouse cursor to specific screen coordinates.

This command moves the mouse cursor to the specified x and y coordinates on the screen
without clicking. The coordinates are measured in pixels from the top-left corner of
the screen (0,0).

Use the --smooth flag for animated movement, and --duration to control the animation speed.`,
		Example: `  # Move cursor instantly to coordinates (800, 600)
  desktop-automation move 800 600

  # Move cursor smoothly to coordinates (800, 600) over 1 second
  desktop-automation move --smooth 800 600

  # Move cursor smoothly to coordinates (800, 600) over 5 seconds
  desktop-automation move --smooth --duration 5.0 800 600

  # Move cursor to the center of a 1920x1080 screen
  desktop-automation move 960 540

  # Move cursor to the top-left corner
  desktop-automation move 0 0`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMoveCommand(cmd, args, smooth, duration)
		},
	}

	// Add flags
	cmd.Flags().BoolVar(&smooth, "smooth", false, "Enable smooth animated movement")
	cmd.Flags().Float64Var(&duration, "duration", 1.0, "Duration in seconds for smooth movement (default: 1.0)")

	return cmd
}

// runMoveCommand handles the move command execution
func runMoveCommand(cmd *cobra.Command, args []string, smooth bool, duration float64) error {
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

	// Validate coordinates are not negative
	if x < 0 {
		return fmt.Errorf("x coordinate cannot be negative: %d", x)
	}
	if y < 0 {
		return fmt.Errorf("y coordinate cannot be negative: %d", y)
	}

	// Get current mouse position
	currentX, currentY := automation.GetPosition()
	fmt.Printf("Current position: (%d, %d)\n", currentX, currentY)
	fmt.Printf("Target position: (%d, %d)\n", x, y)

	// Check if we're already at the target position
	if currentX == x && currentY == y {
		fmt.Println("Already at target position!")
		return nil
	}

	// Perform the movement
	if smooth {
		fmt.Printf("Moving smoothly over %.1f seconds...\n", duration)

		// Start a goroutine to show progress
		done := make(chan bool)
		go func() {
			ticker := time.NewTicker(100 * time.Millisecond)
			defer ticker.Stop()

			for {
				select {
				case <-done:
					return
				case <-ticker.C:
					fmt.Print(".")
				}
			}
		}()

		// Perform smooth movement
		err = automation.SmoothMove(x, y, duration)
		close(done)
		fmt.Println() // New line after dots

		if err != nil {
			return fmt.Errorf("failed to move mouse smoothly: %v", err)
		}
	} else {
		fmt.Println("Moving...")
		err = automation.Move(x, y)
		if err != nil {
			return fmt.Errorf("failed to move mouse: %v", err)
		}
	}

	// Confirm final position
	finalX, finalY := automation.GetPosition()
	fmt.Printf("Final position: (%d, %d)\n", finalX, finalY)

	// Check if we reached the target (allow small tolerance for smooth movement)
	tolerance := 2
	if abs(finalX-x) <= tolerance && abs(finalY-y) <= tolerance {
		fmt.Println("✓ Successfully moved to target position!")
	} else {
		fmt.Printf("⚠ Position may not be exact (target: %d,%d, actual: %d,%d)\n", x, y, finalX, finalY)
	}

	return nil
}

// abs returns the absolute value of an integer
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
