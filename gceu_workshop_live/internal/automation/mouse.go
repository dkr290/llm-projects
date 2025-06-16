package automation

import (
	"fmt"
	"time"

	"github.com/go-vgo/robotgo"
)

// Click performs a mouse click at the specified coordinates with validation
func Click(x, y int) error {
	// Validate coordinates are non-negative
	if x < 0 {
		return fmt.Errorf("x coordinate cannot be negative: %d", x)
	}
	if y < 0 {
		return fmt.Errorf("y coordinate cannot be negative: %d", y)
	}

	// Get screen dimensions for validation
	screenWidth, screenHeight := robotgo.GetScreenSize()
	if x > screenWidth {
		return fmt.Errorf("x coordinate %d exceeds screen width %d", x, screenWidth)
	}
	if y > screenHeight {
		return fmt.Errorf("y coordinate %d exceeds screen height %d", y, screenHeight)
	}

	// Move to position and click
	robotgo.Move(x, y)
	robotgo.Click()

	return nil
}

// GetPosition returns the current mouse position
func GetPosition() (x, y int) {
	x, y = robotgo.GetMousePos()
	return x, y
}

// MoveMouse moves the mouse cursor to the specified coordinates
func MoveMouse(x, y int) error {
	// Validate coordinates are non-negative
	if x < 0 {
		return fmt.Errorf("x coordinate cannot be negative: %d", x)
	}
	if y < 0 {
		return fmt.Errorf("y coordinate cannot be negative: %d", y)
	}

	robotgo.Move(x, y)
	return nil
}

// DoubleClick performs a double click at the specified coordinates
func DoubleClick(x, y int) error {
	// Validate coordinates are non-negative
	if x < 0 {
		return fmt.Errorf("x coordinate cannot be negative: %d", x)
	}
	if y < 0 {
		return fmt.Errorf("y coordinate cannot be negative: %d", y)
	}

	robotgo.Move(x, y)
	robotgo.Click("left", true)
	return nil
}

// RightClick performs a right click at the specified coordinates
func RightClick(x, y int) error {
	// Validate coordinates are non-negative
	if x < 0 {
		return fmt.Errorf("x coordinate cannot be negative: %d", x)
	}
	if y < 0 {
		return fmt.Errorf("y coordinate cannot be negative: %d", y)
	}

	robotgo.Move(x, y)
	robotgo.Click("right")
	return nil
}

// GetMousePos returns the current mouse position (legacy function for compatibility)
func GetMousePos() (int, int) {
	return GetPosition()
}

// Move moves the mouse cursor to the specified coordinates instantly
func Move(x, y int) error {
	// Validate coordinates are non-negative
	if x < 0 {
		return fmt.Errorf("x coordinate cannot be negative: %d", x)
	}
	if y < 0 {
		return fmt.Errorf("y coordinate cannot be negative: %d", y)
	}

	// Get screen dimensions for validation
	screenWidth, screenHeight := robotgo.GetScreenSize()
	if x > screenWidth {
		return fmt.Errorf("x coordinate %d exceeds screen width %d", x, screenWidth)
	}
	if y > screenHeight {
		return fmt.Errorf("y coordinate %d exceeds screen height %d", y, screenHeight)
	}

	// Use MoveSmooth with very short duration for better reliability on macOS
	robotgo.MoveSmooth(x, y, 0.1, 0.1)

	// Add small delay to ensure movement completes
	time.Sleep(50 * time.Millisecond)

	return nil
}

// SmoothMove moves the mouse cursor to the specified coordinates with smooth animation
func SmoothMove(x, y int, duration float64) error {
	// Validate coordinates are non-negative
	if x < 0 {
		return fmt.Errorf("x coordinate cannot be negative: %d", x)
	}
	if y < 0 {
		return fmt.Errorf("y coordinate cannot be negative: %d", y)
	}

	// Validate duration is positive
	if duration <= 0 {
		return fmt.Errorf("duration must be positive: %f", duration)
	}

	// Get screen dimensions for validation
	screenWidth, screenHeight := robotgo.GetScreenSize()
	if x > screenWidth {
		return fmt.Errorf("x coordinate %d exceeds screen width %d", x, screenWidth)
	}
	if y > screenHeight {
		return fmt.Errorf("y coordinate %d exceeds screen height %d", y, screenHeight)
	}

	robotgo.MoveSmooth(x, y, duration, duration)

	// Add small delay to ensure movement completes
	time.Sleep(100 * time.Millisecond)

	return nil
}
