package flow

import (
	"fmt"
	"log"
	"strings"

	"rendellc/gtty/command"
	"rendellc/gtty/serial"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type terminalLine struct {
	content string
}

type ReceivedTerminalLineMsg string

type Model struct {
	lines            []terminalLine
	loadingIndicator spinner.Model
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
	loadingIndicator := spinner.New()
	loadingIndicator.Spinner = spinner.Dot
	loadingIndicator.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return Model{
		lines:            lines,
		loadingIndicator: loadingIndicator,
		ready:            false,
		serialRx:         rx,
	}
}

func (m Model) Init() tea.Cmd {
	log.Printf("Init TerminalDisplay\n")
	return tea.Batch(
		m.loadingIndicator.Tick,
		m.waitForSerialLines(),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}

	var cmd tea.Cmd
	m.loadingIndicator, cmd = m.loadingIndicator.Update(msg)
	cmds = append(cmds, cmd)

	updateContent := false

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerView()) + 5
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight

		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height - verticalMarginHeight)
			m.viewport.YPosition = headerHeight
			m.viewport.HighPerformanceRendering = false
			m.viewport.SetContent("waiting for content to load")
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
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

		m.viewport.SetContent(s.String())
		// m.viewport.GotoBottom()
	}


	m.viewport, cmd = m.viewport.Update(msg)

	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if !m.ready {
		return "Initilizing"
	}

	return fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.viewport.View(), m.footerView())
	// s := ""
	// for _, line := range m.lines {
	// 	s += fmt.Sprintf("%s\n", line.content)
	// }
	// s += m.loadingIndicator.View()
	// return s
}

func (m Model) headerView() string {
	line := strings.Repeat("-", max(0, m.viewport.Width))
	return line
}

func (m Model) footerView() string {
	info := fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100)
	line := strings.Repeat("-", max(0, m.viewport.Width - lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}
