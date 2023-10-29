package models

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

type terminalLine struct {
	content string
}

type TerminalLog struct {
	lines []terminalLine
}

func CreateTerminalLog() TerminalLog {
	lines := make([]terminalLine, 0)
	return TerminalLog{
		lines: lines,
	}
}

func (sp TerminalLog) Init() tea.Cmd {
	log.Println("Initialize terminalLog")
	return nil
}

func (sp TerminalLog) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case InputSubmitMsg:
		log.Printf("TerminalLog got input submit message %+v\n", msg)

		sp.lines = append(sp.lines, terminalLine{
			content: msg.Input,
		})
		return sp, nil
	}

	return sp, nil
}

func (sp TerminalLog) View() string {
	if len(sp.lines) == 0 {
		return "<no terminal output received yet>"
	}

	s := ""
	for _, line := range sp.lines {
		s += fmt.Sprintf("%s\n", line.content)
	}
	return s
}
