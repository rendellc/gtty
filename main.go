package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"rendellc/gtty/command"
	"rendellc/gtty/serial"
	"rendellc/gtty/style"
	"rendellc/gtty/terminal"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type appView int

const (
	appViewTerminal appView = iota
	appViewConnection
	appViewOptions
	appViewHelp
	numberOfAppViews // note: keep at end
)

type appConfig struct {
	SerialConfig   serial.Config
	SimulateSerial bool
	Leader         string
}

type app struct {
	help         help.Model
	keys         keyMap
	command      command.Model
	terminal     terminal.Model
	connection   serial.Connection
	configurator configurator
	config       *appConfig
	appView      appView
	width        int
	height       int
	ready        bool
}

func (a app) Init() tea.Cmd {
	log.Println("Initialize app")
	cmds := []tea.Cmd{}
	cmds = append(cmds, a.command.Init())
	cmds = append(cmds, a.terminal.Init())
	cmds = append(cmds, tea.EnterAltScreen)
	return tea.Batch(cmds...)
}

func (a app) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		footerView := a.footerView(a.command.View())
		footerHeight := lipgloss.Height(footerView)

		a.command.SetMaxWidth(a.width / 6)
		a.terminal.SetVisibleLines(a.height - footerHeight)
		a.help.Width = a.width
		a.help.ShowAll = true
		a.ready = true

		return a, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, a.keys.Quit):
			return a, tea.Quit
		case key.Matches(msg, a.keys.CycleView):
			a.appView = appView((int(a.appView) + 1) % int(numberOfAppViews))
			return a, nil
		}
	case command.InputSubmitMsg:
		var cmd tea.Cmd
		if strings.HasPrefix(msg.Input, a.config.Leader) {
			a.handleLeaderCommandMsg(msg)
		} else if a.appView == appViewTerminal {
			a.connection.GetTransmitter().Send(msg.Input)
			a.terminal, cmd = a.terminal.Update(msg)
		} else if a.appView == appViewConnection {
		} else if a.appView == appViewOptions {
		}

		return a, cmd
	}

	// log.Printf("Main: got msg type: %T %v", msg, msg)
	var cmd tea.Cmd
	var cmds []tea.Cmd
	a.command, cmd = a.command.Update(msg)
	cmds = append(cmds, cmd)
	a.terminal, cmd = a.terminal.Update(msg)
	cmds = append(cmds, cmd)

	return a, tea.Batch(cmds...)
}

func (a *app) handleLeaderCommandMsg(msg command.InputSubmitMsg) {
	cmd := strings.TrimPrefix(msg.Input, a.config.Leader)

	if strings.HasPrefix(cmd, "connect") {
		a.connection.Start()
	} else if strings.HasPrefix(cmd, "disconnect") {
		a.connection.Close()
	} else {
		a.configurator.DoCommand(cmd, a.config)
	}

}

func (a app) getViewString() string {
	switch a.appView {
	case appViewTerminal:
		return "terminal"
	case appViewConnection:
		return "connection"
	case appViewOptions:
		return "options"
	case appViewHelp:
		return "help"
	}

	return "unknown"
}

func (a app) footerView(commandStr string) string {
	if !a.ready {
		return ""
	}
	command := style.CommandFooter.Render(commandStr)
	view := style.ViewFooter.Render(a.getViewString())
	info := style.InfoFooter.Render(fmt.Sprintf("%d%%", 0))

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
	footerView := a.footerView(a.command.View())
	footerHeight := lipgloss.Height(footerView)
	mainWidth := int(0.7 * float32(a.width))
	mainMargin := int((a.width - mainWidth) / 2)
	mainStyle := style.MainView.Width(mainWidth).Height(a.height-footerHeight-2).Margin(0, mainMargin).Render
	terminalLines := a.height - footerHeight - 2
	a.terminal.SetVisibleLines(terminalLines)

	mainView := ""
	if a.appView == appViewTerminal {
		mainView = a.terminal.View()
	} else if a.appView == appViewConnection {
		mainView = viewConnectionConfig(a.config)
	} else if a.appView == appViewOptions {
		mainView = "option viewer"
	} else if a.appView == appViewHelp {
		mainView = a.help.View(a.keys)
	}

	return fmt.Sprintf("%s\n%s",
		mainStyle(mainView),
		footerView,
	)
}

func main() {
	config := appConfig{}
	flag.BoolVar(&config.SimulateSerial, "sim", false, "Simulate serial data")
	flag.StringVar(&config.SerialConfig.Device, "device", "COM8", "Serial device")
	flag.IntVar(&config.SerialConfig.BaudRate, "baud", 9600, "Serial connection baud rate")
	flag.IntVar(&config.SerialConfig.DataBits, "databits", 8, "Serial connection data bits")
	flag.IntVar(&config.SerialConfig.StopBits, "stopbits", 1, "Serial connection stop bits")
	flag.StringVar(&config.SerialConfig.Parity, "parity", "odd", "Serial connection parity")
	flag.StringVar(&config.Leader, "leader", ":", "Leader key for use with internal commands")
	flag.Parse()

	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()

	config.SerialConfig.Timeout = 5 * time.Second
	config.SerialConfig.TransmitNewline = "\r\n"

	log.Printf("Config is %+v", config)

	var connection serial.Connection
	if config.SimulateSerial {
		connection = serial.SimulateConnection(500 * time.Millisecond)
	} else {
		connection = serial.CreateConnection(config.SerialConfig)
	}

	rx := connection.GetReceiver()
	tx := connection.GetTransmitter()

	defer connection.Close()

	app := app{
		config:     &config,
		help:       help.New(),
		keys:       keys,
		connection: connection,
		command:    command.CreateCommandInput(&tx),
		terminal:   terminal.CreateFlowModel(&rx),
		ready:      false,
		configurator: NewConfigurator(),
	}

	if _, err := tea.NewProgram(app).Run(); err != nil {
		log.Fatal(err)
	}
}
