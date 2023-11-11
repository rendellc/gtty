package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"rendellc/gtty/command"
	"rendellc/gtty/flow"
	"rendellc/gtty/style"
	"rendellc/gtty/serial"

	tea "github.com/charmbracelet/bubbletea"
)

type app struct {
	commandInput    command.CommandInputModel
	terminalDisplay flow.Model
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
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return a, tea.Quit
		}
	}

	// log.Printf("Main: got msg type: %T %v", msg, msg)
	inputModel, cmd1 := a.commandInput.Update(msg)
	logModel, cmd2 := a.terminalDisplay.Update(msg)

	a.commandInput = inputModel.(command.CommandInputModel)
	a.terminalDisplay = logModel.(flow.Model)

	return a, tea.Batch(cmd1, cmd2)
}

func (a app) View() string {
	return fmt.Sprintf("%s\n\n%s\n\n%s",
		a.commandInput.View(),
		style.HelpLine(),
		a.terminalDisplay.View(),
	)
}

type AppConfig struct {
	SerialConfig   serial.Config
	SimulateSerial bool
}

func main() {
	config := AppConfig{}

	flag.BoolVar(&config.SimulateSerial, "sim", false, "Simulate serial data")
	flag.Parse()

	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()



	// config.SerialConfig.Device = "COM8"
	config.SerialConfig.Device = "/dev/random"
	config.SerialConfig.BaudRate = 9600
	config.SerialConfig.DataBits = 8
	config.SerialConfig.StopBits = 1
	config.SerialConfig.Parity = "odd"
	config.SerialConfig.Timeout = 5*time.Second
	config.SerialConfig.TransmitNewline = "\r\n"

	connection := serial.CreateConnection(config.SerialConfig)
	defer connection.Close()

	rx, tx, err := connection.Start()
	if err != nil {
		log.Printf("Error starting listener: %v", err.Error())
	}
	time.Sleep(1*time.Second)
	log.Printf("Rx: %s\n", rx.Get())
	time.Sleep(1*time.Second)
	log.Printf("Rx: %s\n", rx.Get())
	time.Sleep(1*time.Second)
	log.Printf("Rx: %s\n", rx.Get())

	app := app{
		commandInput:    command.CreateCommandInput(&tx),
		terminalDisplay: flow.CreateFlowModel(&rx),
	}

	if _, err := tea.NewProgram(app).Run(); err != nil {
		log.Fatal(err)
	}
}
