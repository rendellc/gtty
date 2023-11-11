package flow

import (
	"log"
	"strings"

	"rendellc/gtty/command"
	"rendellc/gtty/serial"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type terminalLine struct {
	content string
}

type ReceivedTerminalLineMsg string

type Model struct {
	lines            []terminalLine
	ready            bool
	viewport         viewport.Model
	serialRx         *serial.Receiver
}

func (m *Model) waitForSerialLines() tea.Cmd {
	return func() tea.Msg {
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

func (m Model) Init() tea.Cmd {
	log.Printf("Init TerminalDisplay\n")
	return tea.Batch(
		m.waitForSerialLines(),
	)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	cmds := []tea.Cmd{}

	updateContent := false
	ypos := 2

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height - ypos)
			m.viewport.HighPerformanceRendering = false
			m.viewport.SetContent("waiting for content to load")
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - ypos
		}

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyDelete:
			m.lines = m.lines[:0]
			updateContent = true
		}
	case command.InputSubmitMsg:
		m.lines = append(m.lines, terminalLine{
			content: msg.Input,
		})
		updateContent = true
	case ReceivedTerminalLineMsg:

		m.lines = append(m.lines, terminalLine{
			content: string(msg),
		})
		cmds = append(cmds, m.waitForSerialLines())
		updateContent = true
	}

	if updateContent {
		s := strings.Builder{}
		for _, line := range m.lines {
			s.WriteString(line.content)
			s.WriteString("\n")
		}

		atBottom := m.viewport.AtBottom()
		m.viewport.SetContent(s.String())
		
		if atBottom {
			m.viewport.GotoBottom()
		}
	}


	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if !m.ready {
		return "Initilizing"
	}

	return m.viewport.View()
}

func (m Model) headerView() string {
	line := strings.Repeat("─", max(0, m.viewport.Width))
	return line
}


func (m Model) ScrollPercent() int {
	return int(100*m.viewport.ScrollPercent())
}
