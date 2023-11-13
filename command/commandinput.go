package command

import (
	"log"
	"rendellc/gtty/serial"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	textInput textinput.Model
	serialTx  *serial.Transmitter
}

type InputSubmitMsg struct {
	Input string
}

func (m *Model) inputSubmitCmd(cmd string) tea.Cmd {
	return func() tea.Msg {
		m.serialTx.Send(cmd)
		return InputSubmitMsg{
			Input: cmd,
		}
	}
}

func CreateCommandInput(tx *serial.Transmitter) Model {
	ti := textinput.New()
	//ti.Placeholder = ""
	ti.Prompt = "Command: "
	ti.CharLimit = 0
	ti.Width = 0
	ti.Cursor.BlinkSpeed = 1 * time.Second
	ti.Focus()

	return Model{
		textInput: ti,
		serialTx:  tx,
	}
}

func (m Model) Init() tea.Cmd {
	log.Println("Initialize command.Model")
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			input := m.textInput.Value()
			cmd := m.inputSubmitCmd(input)
			m.textInput.SetValue("")
			return m, cmd
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return m.textInput.View()
}

func (m *Model) SetMaxWidth(width int) {
	m.textInput.Width = width
}
