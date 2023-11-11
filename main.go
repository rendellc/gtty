package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"rendellc/gtty/command"
	"rendellc/gtty/flow"
	"rendellc/gtty/serial"
	"rendellc/gtty/style"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type appConfig struct {
	SerialConfig   serial.Config
	SimulateSerial bool
}

type app struct {
	commandInput    command.CommandInputModel
	terminalDisplay flow.Model
	width           int
	height          int
}

func (a app) Init() tea.Cmd {
	log.Println("Initialize app")
	cmds := []tea.Cmd{}
	cmds = append(cmds, a.commandInput.Init())
	cmds = append(cmds, a.terminalDisplay.Init())
	cmds = append(cmds, tea.EnterAltScreen)
	return tea.Batch(cmds...)
}

func (a app) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return a, tea.Quit
		}
	}

	// log.Printf("Main: got msg type: %T %v", msg, msg)
	inputModel, cmd1 := a.commandInput.Update(msg)
	logModel, cmd2 := a.terminalDisplay.Update(msg)

	a.commandInput = inputModel
	a.terminalDisplay = logModel

	return a, tea.Batch(cmd1, cmd2)
}

func (a app) footerView(commandStr string, scrollPercent int) string {
	command := style.CommandFooter.Render(commandStr)
	info := style.InfoFooter.Render(fmt.Sprintf("%d%%", scrollPercent))

	line := strings.Repeat("─", max(0, a.width - lipgloss.Width(command)-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, command, line, info)
}

func (a app) View() string {
	return fmt.Sprintf("%s\n%s",
		a.terminalDisplay.View("hello"),
		a.footerView(a.commandInput.View(), a.terminalDisplay.ScrollPercent()),
		// a.commandInput.View(),
		// style.HelpLine(),
	)
}

func main() {
	config := appConfig{}
	flag.BoolVar(&config.SimulateSerial, "sim", false, "Simulate serial data")
	flag.Parse()

	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()

	config.SerialConfig.Device = "COM8"
	config.SerialConfig.BaudRate = 9600
	config.SerialConfig.DataBits = 8
	config.SerialConfig.StopBits = 1
	config.SerialConfig.Parity = "odd"
	config.SerialConfig.Timeout = 5 * time.Second
	config.SerialConfig.TransmitNewline = "\r\n"

	log.Printf("Config is %+v", config)

	var connection serial.Connection
	if config.SimulateSerial {
		connection = serial.SimulateConnection(500 * time.Millisecond)
	} else {
		connection = serial.CreateConnection(config.SerialConfig)
	}
	defer connection.Close()

	rx, tx, err := connection.Start()
	if err != nil {
		log.Printf("Error starting listener: %v", err.Error())
	}

	commandInput := command.CreateCommandInput(&tx)
	terminalDisplay := flow.CreateFlowModel(&rx)
	app := app{
		commandInput:    commandInput,
		terminalDisplay: terminalDisplay,
	}

	if _, err := tea.NewProgram(app).Run(); err != nil {
		log.Fatal(err)
	}
}
