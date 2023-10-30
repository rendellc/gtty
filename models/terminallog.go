package models

import (
	"fmt"
	"log"
	"time"

	// "log"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dsyx/serialport-go"
)

type terminalLine struct {
	content string
}

type ReceivedTerminalLineMsg string

type TerminalDisplay struct {
	lines            []terminalLine
	loadingIndicator spinner.Model
	portName         string
	serialConfig     serialport.Config
	rxChan           chan string
}

func listenForActivity(sub chan string, portName string, serialConfig serialport.Config) tea.Cmd {
	connection, err := serialport.Open(portName, serialConfig)
	if err != nil {
		log.Fatalf("Unable to open serial port: %s", err)
	}
	defer connection.Close()

	return func() tea.Msg {
		buf := make([]byte, 1024)
		line := ""
		for {
			n, _ := connection.Read(buf)
			if n == 0 {
				d := 1000 * time.Millisecond
				log.Printf("No data. Sleeping for %v", d)
				time.Sleep(d)
				continue
			}
			readString := string(buf[:n])
			log.Printf("Received %s\n", readString)
			// Split readString on newlines
			line += readString
			sub <- "<data here?>"
			line = ""
		}
	}
}

func waitForActivity(sub chan string) tea.Cmd {
	log.Println("Starting waitForActivity")
	return func() tea.Msg {
		return ReceivedTerminalLineMsg(<-sub)
	}
}

func CreateTerminalDisplay(name string, config serialport.Config) TerminalDisplay {
	log.Printf("Create terminal display for %s\n", name)

	lines := make([]terminalLine, 0)
	loadingIndicator := spinner.New()
	loadingIndicator.Spinner = spinner.Dot
	loadingIndicator.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return TerminalDisplay{
		lines:            lines,
		loadingIndicator: loadingIndicator,
		rxChan:           make(chan string),
		portName:         name,
		serialConfig:     config,
	}
}

func (sp TerminalDisplay) Init() tea.Cmd {
	log.Printf("Starting TerminalDisplay\n")
	return tea.Batch(
		sp.loadingIndicator.Tick,
		listenForActivity(sp.rxChan, sp.portName, sp.serialConfig),
		waitForActivity(sp.rxChan),
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
		cmds = append(cmds, waitForActivity(sp.rxChan))
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
