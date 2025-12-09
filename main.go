package main

import (
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	cursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("200"))
	vault       string
)

type model struct {
	newFileInput          textinput.Model
	isNewFileInputVisible bool
	currentFile           *os.File
	contentTextarea       textarea.Model
}

func initialModel() model {

	// create base directory to save notes
	err := os.MkdirAll(vault, 0750)
	if err != nil {
		log.Fatal(err.Error())
	}

	// initialize new file input
	ti := textinput.New()
	ti.Placeholder = "What you like to call it?"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 50
	ti.Cursor.Style = cursorStyle
	ti.TextStyle = cursorStyle
	ti.PromptStyle = cursorStyle

	// initialize content textarea
	ta := textarea.New()
	ta.Placeholder = "Write something inside your file"
	ta.ShowLineNumbers = false
	ta.Focus()

	return model{
		newFileInput:          ti,
		isNewFileInputVisible: false,
		contentTextarea:       ta,
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

	if m.currentFile != nil {
		view = m.contentTextarea.View()
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
		case "ctrl+c", "esc":
			return m, tea.Quit
			
		case "ctrl+n":
			m.isNewFileInputVisible = true
			return m, nil

		case "ctrl+s":
			if m.currentFile == nil {
				break
			}

			if err := m.currentFile.Truncate(0); err != nil {
				fmt.Println("Unable to truncate file")
				return m, nil
			}
			if _, err :=  m.currentFile.Seek(0,0); err != nil {
				fmt.Println("Unable to seek file")
				return m, nil
			}
			if _, err := m.currentFile.WriteString(m.contentTextarea.Value()); err != nil {
				fmt.Println("Unable to write in file")
				return m, nil
			}
			if err := m.currentFile.Close(); err != nil {
				fmt.Println("Unable to close file")
			}
			m.currentFile = nil
			m.contentTextarea.SetValue("")
			return m, nil

		case "enter":
			// preserve default behavior if textarea is open
			if m.currentFile != nil {
				break
			}

			filename := m.newFileInput.Value()
			filepath := fmt.Sprintf("%s/%s.md", vault, filename)

			if _, err := os.Stat(filepath); err == nil {
				return m, nil
			}

			f, err := os.Create(filepath)
			if err != nil {
				log.Fatal(err)
			}
			m.currentFile = f
			m.newFileInput.SetValue("")
			m.isNewFileInputVisible = false
			return m, nil
		}

	}
	if m.isNewFileInputVisible {
		m.newFileInput, cmd = m.newFileInput.Update(msg)
		return m, cmd
	}
	if m.currentFile != nil {
		m.contentTextarea, cmd = m.contentTextarea.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Error getting home directory", err.Error())
	}
	vault = fmt.Sprintf("%s/.noteify", homeDir)
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
