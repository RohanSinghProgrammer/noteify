package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	newFileInput          textinput.Model
	isNewFileInputVisible bool
}

var (
	cursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("200"))
)

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "What you like to call it?"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 50
	ti.Cursor.Style = cursorStyle
	ti.TextStyle = cursorStyle
	ti.PromptStyle = cursorStyle

	return model{
		newFileInput:          ti,
		isNewFileInputVisible: false,
	}
}

func (m model) View() string {
	// Styles
	var headerStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("16")).
		Background(lipgloss.Color("200")).
		Padding(0, 2)

	var helpStyle = lipgloss.NewStyle().
		Italic(true).
		Foreground(lipgloss.Color("16")).
		Background(lipgloss.Color("51")).
		Padding(0, 2)

		// Texts
	header := headerStyle.Render("Welcome to Noteify!")
	view := ""
	help := helpStyle.Render("Ctrl+N: new file | Ctrl+L: list | Ctrl+S: save | Ctrl+Q: quit | Esc: back/save")

	if m.isNewFileInputVisible {
		view = m.newFileInput.View()
	}

	s := fmt.Sprintf("\n%s\n\n%s\n\n%s", header, view, help)

	// Send the UI for rendering
	return s
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "ctrl+n":
			m.isNewFileInputVisible = true
			return m, nil
		}

	}
	if m.isNewFileInputVisible {
		m.newFileInput, cmd = m.newFileInput.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
