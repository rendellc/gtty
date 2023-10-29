package models

import (
	"fmt"
	// "log"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type terminalLine struct {
	content string
}

type TerminalLog struct {
	lines   []terminalLine
	spinner spinner.Model
}

func CreateTerminalLog() TerminalLog {
	lines := make([]terminalLine, 0)

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return TerminalLog{
		lines:   lines,
		spinner: s,
	}
}

func (sp TerminalLog) Init() tea.Cmd {
	return sp.spinner.Tick
}

func (sp TerminalLog) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	sp.spinner, cmd = sp.spinner.Update(msg)

	switch msg := msg.(type) {
	case InputSubmitMsg:
		// log.Printf("TerminalLog got input submit message %+v\n", msg)
		sp.lines = append(sp.lines, terminalLine{
			content: msg.Input,
		})
	}

	return sp, cmd
}

func (sp TerminalLog) View() string {
	s := ""
	for _, line := range sp.lines {
		s += fmt.Sprintf("%s\n", line.content)
	}
	s += sp.spinner.View()
	return s
}
