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
	visibleLines int
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

func (m *Model) SetVisibleLines(lineCount int) {
	m.visibleLines = lineCount
}

func (m Model) AtMostRecent() bool {
	return m.scrollIndex + m.visibleLines >= len(m.lines)
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

	atMostRecent := m.AtMostRecent()

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
	if atMostRecent {
		m.scrollIndex = max(0, len(m.lines) - m.visibleLines)
	}


	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if len(m.lines) < m.visibleLines {
		m.visibleLines = len(m.lines)
	}
	if m.visibleLines == 0 {
		return "<Waiting for data>"
	}

	// Know now that lineCount <= len(m.lines)
	// Need to enforce for each 
	//    i in [m.scrollIndex, m.scrollIndex + lineCount)
	//    0 <= i < len(m.lines)
	lastIndex := m.scrollIndex + m.visibleLines - 1
	maxIndex := len(m.lines) - 1
	if lastIndex > maxIndex {
		countTooMany := lastIndex - maxIndex
		m.scrollIndex -= countTooMany
	}

	lines := []string{}
	for i := m.scrollIndex; i < m.scrollIndex + m.visibleLines; i++ {
		lines = append(lines, m.lines[i].content)
	}

	return strings.Join(lines, "\n")
}

func min(a, b int) int {
	if a < b{
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b{
		return a
	}
	return b
}
