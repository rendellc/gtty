package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"rendellc/gtty/style"
	"rendellc/gtty/models"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dsyx/serialport-go"
)

// func createReaderRoutine(filename string) (<-chan string, chan<- int)  {
// 	content := make(chan string)
// 	commands := make(chan int)
//
// 	go func(filename string, content chan<- string, commands <-chan int) {
// 		file, err := os.Open(filename)
// 		if err != nil {
// 			log.Panic(err)
// 		}
// 		defer file.Close()
//
// 		scanner := bufio.NewScanner(file)
// 		for scanner.Scan() {
// 			content <- scanner.Text()
//
// 			select {
// 			case <-commands:
// 				log.Println("Command received. Quitting")
// 				close(content)
// 				return
// 			default:
// 			}
// 		}
// 	}(filename, content, commands)
//
// 	return content, commands
// }

type app struct {
	commandInput    models.CommandInputModel
	terminalDisplay models.TerminalDisplay
}

func (a app) Init() tea.Cmd {
	log.Println("Initialize app")
	cmds := []tea.Cmd{}
	cmds = append(cmds, a.commandInput.Init())
	cmds = append(cmds, a.terminalDisplay.Init())
	cmds = append(cmds, tea.EnterAltScreen)
	cmds = append(cmds, nil)
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

	a.commandInput = inputModel.(models.CommandInputModel)
	a.terminalDisplay = logModel.(models.TerminalDisplay)

	return a, tea.Batch(cmd1, cmd2)
}

func (a app) View() string {
	return fmt.Sprintf("%s\n\n%s\n\n%s",
		a.commandInput.View(),
		style.HelpLine(),
		a.terminalDisplay.View(),
	)
}

func main() {
	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		defer f.Close()
	}

	config := serialport.Config{
		BaudRate: serialport.BR9600,
		DataBits: serialport.DB8,
		StopBits: serialport.SB1,
		Parity:   serialport.PO,
		Timeout:  5 * time.Second,
	}

	app := app{
		commandInput:    models.CreateCommandInput(),
		terminalDisplay: models.CreateTerminalDisplay("COM8", config),
	}

	if _, err := tea.NewProgram(app).Run(); err != nil {
		log.Fatal(err)
	}
}
