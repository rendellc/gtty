package models

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type CommandInputModel struct {
	textInput textinput.Model
}

type InputSubmitMsg struct {
	Input string
}
func InputSubmitCmd(input string) tea.Cmd {
	return func() tea.Msg {
		return InputSubmitMsg{
			Input: input,
		};
	}
}


func CreateCommandInput() CommandInputModel {
	ti := textinput.New()
	ti.Placeholder = "e.g help"
	ti.Prompt = "Command: "
	ti.Focus()
	ti.CharLimit = 128
	ti.Width = 20

	return CommandInputModel{
		textInput: ti,
	}
}

func (m CommandInputModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m CommandInputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	log.Printf("Update CommandInputModel: %+v\n", msg)
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			input := m.textInput.Value()
			cmd := InputSubmitCmd(input)
			return m, cmd

		}
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}



func (m CommandInputModel) View() string {
	return fmt.Sprintf(
		"%s\n\n(esc to quit)",
		m.textInput.View(),
	)
}

