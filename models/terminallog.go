package models

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type terminalLine struct {
	content string
}

type ReceivedTerminalLineMsg string

type TerminalDisplay struct {
	lines            []terminalLine
	loadingIndicator spinner.Model
	serialLine           <-chan string
}

func (sp *TerminalDisplay) waitForSerialLines() tea.Cmd {
	return func() tea.Msg {
		return ReceivedTerminalLineMsg(<-sp.serialLine)
	}
}

func CreateTerminalDisplay(serialLines <-chan string) TerminalDisplay {
	lines := make([]terminalLine, 0)
	loadingIndicator := spinner.New()
	loadingIndicator.Spinner = spinner.Dot
	loadingIndicator.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return TerminalDisplay{
		lines:            lines,
		loadingIndicator: loadingIndicator,
		serialLine: serialLines,
	}
}

func (sp TerminalDisplay) Init() tea.Cmd {
	log.Printf("Init TerminalDisplay\n")
	return tea.Batch(
		sp.loadingIndicator.Tick,
		sp.waitForSerialLines(),
	)
}

func (sp TerminalDisplay) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}

	var cmd tea.Cmd
	sp.loadingIndicator, cmd = sp.loadingIndicator.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyDelete:
			sp.lines = sp.lines[:0]
		}
	case InputSubmitMsg:
		sp.lines = append(sp.lines, terminalLine{
			content: msg.Input,
		})
	case ReceivedTerminalLineMsg:
		sp.lines = append(sp.lines, terminalLine{
			content: string(msg),
		})
		cmds = append(cmds, sp.waitForSerialLines())
	}

	return sp, tea.Batch(cmds...)
}

func (sp TerminalDisplay) View() string {
	s := ""
	for _, line := range sp.lines {
		s += fmt.Sprintf("%s\n", line.content)
	}
	s += sp.loadingIndicator.View()
	return s
}
