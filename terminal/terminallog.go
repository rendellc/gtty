package terminal

import (
	"log"
	"strings"
	"time"

	"rendellc/gtty/command"
	"rendellc/gtty/serial"

	tea "github.com/charmbracelet/bubbletea"
)

type terminalLine struct {
	content string
}

type ReceivedTerminalLineMsg string

type Model struct {
	lines            []terminalLine
	ready            bool
	scrollIndex int
	serialRx         *serial.Receiver
}

func (m *Model) waitForSerialLines() tea.Cmd {
	return func() tea.Msg {
		if m.serialRx == nil {
			// we dont have an active receiver yet. Sleep for a short while,
			// then try again.
			time.Sleep(1*time.Second)
			return m.waitForSerialLines() 
		}

		return ReceivedTerminalLineMsg(m.serialRx.Get())
	}
}

func CreateFlowModel(rx *serial.Receiver) Model {
	lines := make([]terminalLine, 0)

	return Model{
		lines:            lines,
		ready:            false,
		serialRx:         rx,
	}
}

func (m *Model) SetReceiver(rx *serial.Receiver) {
	m.serialRx = rx
}

func (m Model) Init() tea.Cmd {
	log.Printf("Init TerminalDisplay\n")
	return tea.Batch(
		m.waitForSerialLines(),
	)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	cmds := []tea.Cmd{}


	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyDelete:
			m.lines = m.lines[:0]
		}
	case command.InputSubmitMsg:
		m.lines = append(m.lines, terminalLine{
			content: msg.Input,
		})
	case ReceivedTerminalLineMsg:

		m.lines = append(m.lines, terminalLine{
			content: string(msg),
		})
		cmds = append(cmds, m.waitForSerialLines())
	}

	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View(lineCount int) string {
	if lineCount == 0 {
		return "Initilizing"
	}

	lines := []string{}
	for i, termline := range m.lines {
		lines[i] = termline.content
	}

	return strings.Join(lines, "\n")
}
