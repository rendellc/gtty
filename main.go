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

type appView int

const (
	appViewTerminal appView = iota
	appViewOptions
	numberOfAppViews // note: keep at end
)

type appConfig struct {
	SerialConfig   serial.Config
	SimulateSerial bool
}

type app struct {
	commandInput    command.CommandInputModel
	terminalDisplay flow.Model
	config          *appConfig
	appView         appView
	width           int
	height          int
	ready           bool
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
		a.ready = true

		a.commandInput.SetMaxWidth(a.width / 6)
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return a, tea.Quit
		case tea.KeyTab:
			a.appView = appView((int(a.appView) + 1) % int(numberOfAppViews))
		}
	}

	// log.Printf("Main: got msg type: %T %v", msg, msg)
	inputModel, cmd1 := a.commandInput.Update(msg)
	logModel, cmd2 := a.terminalDisplay.Update(msg)

	a.commandInput = inputModel
	a.terminalDisplay = logModel

	return a, tea.Batch(cmd1, cmd2)
}

func (a app) getViewString() string {
	switch a.appView {
	case appViewTerminal:
		return "terminal"
	case appViewOptions:
		return "options"
	}

	return "unknown"
}

func (a app) footerView(commandStr string, scrollPercent int) string {
	if !a.ready {
		return ""
	}
	command := style.CommandFooter.Render(commandStr)
	view := style.ViewFooter.Render(a.getViewString())
	info := style.InfoFooter.Render(fmt.Sprintf("%d%%", scrollPercent))

	// line is split across two regions
	// <cmdbox> ---l--- <viewbox> ---r--- <infobox>
	viewWidth := lipgloss.Width(view)
	lineWidthLeft := max(0, a.width/2-lipgloss.Width(command)-viewWidth/2)
	lineWidthRight := max(0, a.width/2-lipgloss.Width(info)-viewWidth/2)
	if len(view)%2 == 1 {
		lineWidthRight -= 1
	}
	lineLeft := strings.Repeat("─", lineWidthLeft)
	lineRight := strings.Repeat("─", lineWidthRight)

	return lipgloss.JoinHorizontal(lipgloss.Center, command, lineLeft, view, lineRight, info)
}

func (a app) View() string {
	mainStyle := style.MainView.Width(a.width).Height(a.height - 3).Render

	mainView := ""
	if a.appView == appViewTerminal {
		mainView = a.terminalDisplay.View()
	} else if a.appView == appViewOptions {
		mainView = "option viewer"
	}

	return fmt.Sprintf("%s\n%s",
		mainStyle(mainView),
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
		ready: false,
	}

	if _, err := tea.NewProgram(app).Run(); err != nil {
		log.Fatal(err)
	}
}
