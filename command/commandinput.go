package command

import (
	"rendellc/gtty/serial"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type CommandInputModel struct {
	textInput textinput.Model
	serialTx *serial.Transmitter
}

type InputSubmitMsg struct {
	Input string
}

func (m *CommandInputModel) inputSubmitCmd(cmd string) tea.Cmd {
	return func() tea.Msg {
		m.serialTx.Send(cmd)
		return InputSubmitMsg{
			Input: cmd,
		}
	}
}

func CreateCommandInput(tx *serial.Transmitter) CommandInputModel {
	ti := textinput.New()
	ti.Placeholder = "e.g help"
	ti.Prompt = "Command: "
	ti.Focus()
	ti.CharLimit = 0
	ti.Width = 0
	ti.Cursor.BlinkSpeed = 1 * time.Second

	return CommandInputModel{
		textInput: ti,
		serialTx: tx,
	}
}

func (m CommandInputModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m CommandInputModel) Update(msg tea.Msg) (CommandInputModel, tea.Cmd) {
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

func (m CommandInputModel) View() string {
	return m.textInput.View()
}
