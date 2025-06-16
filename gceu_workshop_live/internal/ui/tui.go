package ui

import (
	"fmt"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dmahlow/desktop-automation/internal/automation"
)

type model struct {
	choices  []string
	cursor   int
	selected map[int]struct{}
	state    string
	input    string
	result   string
}

func initialModel() model {
	return model{
		choices: []string{
			"Move Mouse",
			"Click Mouse",
			"Type Text",
			"Get Mouse Position",
			"Get Screen Size",
			"Quit",
		},
		selected: make(map[int]struct{}),
		state:    "menu",
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.state {
		case "menu":
			return m.updateMenu(msg)
		case "input":
			return m.updateInput(msg)
		case "result":
			return m.updateResult(msg)
		}
	}
	return m, nil
}

func (m model) updateMenu(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(m.choices)-1 {
			m.cursor++
		}
	case "enter", " ":
		switch m.cursor {
		case 5: // Quit
			return m, tea.Quit
		case 3: // Get Mouse Position
			x, y := automation.GetMousePos()
			m.result = fmt.Sprintf("Mouse position: %d, %d", x, y)
			m.state = "result"
		case 4: // Get Screen Size
			w, h := automation.GetScreenSize()
			m.result = fmt.Sprintf("Screen size: %dx%d", w, h)
			m.state = "result"
		default:
			m.state = "input"
			m.input = ""
		}
	}
	return m, nil
}

func (m model) updateInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc":
		m.state = "menu"
		m.input = ""
	case "enter":
		m.result = m.executeAction()
		m.state = "result"
	case "backspace":
		if len(m.input) > 0 {
			m.input = m.input[:len(m.input)-1]
		}
	default:
		m.input += msg.String()
	}
	return m, nil
}

func (m model) updateResult(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "enter", "esc", " ":
		m.state = "menu"
		m.result = ""
	}
	return m, nil
}

func (m model) executeAction() string {
	switch m.cursor {
	case 0: // Move Mouse
		coords := strings.Fields(m.input)
		if len(coords) != 2 {
			return "Error: Enter X Y coordinates (e.g., 100 200)"
		}
		x, err1 := strconv.Atoi(coords[0])
		y, err2 := strconv.Atoi(coords[1])
		if err1 != nil || err2 != nil {
			return "Error: Invalid coordinates"
		}
		automation.MoveMouse(x, y)
		return fmt.Sprintf("Moved mouse to %d, %d", x, y)

	case 1: // Click Mouse
		coords := strings.Fields(m.input)
		if len(coords) != 2 {
			return "Error: Enter X Y coordinates (e.g., 100 200)"
		}
		x, err1 := strconv.Atoi(coords[0])
		y, err2 := strconv.Atoi(coords[1])
		if err1 != nil || err2 != nil {
			return "Error: Invalid coordinates"
		}
		automation.Click(x, y)
		return fmt.Sprintf("Clicked at %d, %d", x, y)

	case 2: // Type Text
		if m.input == "" {
			return "Error: Enter text to type"
		}
		automation.TypeText(m.input)
		return fmt.Sprintf("Typed: %s", m.input)
	}
	return "Unknown action"
}

func (m model) View() string {
	switch m.state {
	case "menu":
		return m.viewMenu()
	case "input":
		return m.viewInput()
	case "result":
		return m.viewResult()
	}
	return ""
}

func (m model) viewMenu() string {
	s := "Desktop Automation CLI\n\n"

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}

	s += "\nPress q to quit, arrow keys to navigate, enter to select.\n"
	return s
}

func (m model) viewInput() string {
	action := m.choices[m.cursor]
	s := fmt.Sprintf("Selected: %s\n\n", action)

	switch m.cursor {
	case 0, 1: // Move/Click Mouse
		s += "Enter X Y coordinates (e.g., 100 200):\n"
	case 2: // Type Text
		s += "Enter text to type:\n"
	}

	s += fmt.Sprintf("> %s\n\n", m.input)
	s += "Press Enter to execute, Esc to cancel, Ctrl+C to quit.\n"
	return s
}

func (m model) viewResult() string {
	s := "Result:\n\n"
	s += m.result + "\n\n"
	s += "Press Enter or Space to continue, Ctrl+C to quit.\n"
	return s
}

// StartTUI starts the terminal user interface
func StartTUI() error {
	p := tea.NewProgram(initialModel())
	_, err := p.Run()
	return err
}
