package main

import (
	"fmt"
	"log"
	"os"

	"rendellc/gtty/models"

	tea "github.com/charmbracelet/bubbletea"
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
	commandInput models.CommandInputModel
	terminalLog  models.TerminalLog
}

func (a app) Init() tea.Cmd {
	log.Println("Initialize app")
	return nil
}

func (a app) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	log.Printf("Update app component: %T %+v\n", msg, msg)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return a, tea.Quit
		}

	}

	inputModel, cmd1 := a.commandInput.Update(msg)
	logModel, cmd2 := a.terminalLog.Update(msg)

	a.commandInput = inputModel.(models.CommandInputModel)
	a.terminalLog = logModel.(models.TerminalLog)

	return a, tea.Batch(cmd1, cmd2)
}

func (a app) View() string {
	return fmt.Sprintf("%s\n\n%s",
		a.commandInput.View(),
		a.terminalLog.View(),
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

	app := app{
		commandInput: models.CreateCommandInput(),
		terminalLog:  models.CreateTerminalLog(),
	}
	if _, err := tea.NewProgram(app).Run(); err != nil {
		log.Fatal(err)
	}
}
