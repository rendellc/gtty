package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"rendellc/gtty/models"
	"rendellc/gtty/style"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dsyx/serialport-go"
)

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

type SerialHandler struct {
	portName string
	config serialport.Config
	connection *serialport.SerialPort
	buf []byte
}

func (h SerialHandler) ReadSerial() string {
	if h.connection == nil {
		log.Fatalf("ReadSerial called, but connection is nil")
	}

	n, _ := h.connection.Read(h.buf)
	return string(h.buf[:n])
}

func (h *SerialHandler) Open() error {
	var err error
	h.connection, err = serialport.Open(h.portName, h.config)
	if err != nil {
		return fmt.Errorf("Unable to open serial port: %s", err)
	}
	return nil
}

func (h *SerialHandler) Close() {
	if h.connection == nil {
		return
	}

	h.connection.Close()
}

func (h *SerialHandler) readSerial(out chan<- string) {
	for {
		readString := h.ReadSerial()
		if len(readString) == 0 {
			d := 1 * time.Millisecond
			log.Printf("No data. Sleeping for %v", d)
			time.Sleep(d)
			continue
		}
		log.Printf("Received %s\n", readString)
		out <- readString
	}
}

func mergeSerialDataToLines(serialData <-chan string, serialLine chan<- string) {
	partialLines := ""
	for {
		select {
		case s := <-serialData:
			partialLines += s
		}

		lines := strings.Split(partialLines, "\n")
		for i, line := range lines {
			if i < len(lines) - 1 {
				line = strings.Trim(line, "\r")
				serialLine <- line
			}
		}
		partialLines = lines[len(lines) - 1]
}
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
	serialHandler := new(SerialHandler)
	serialHandler.portName = "COM8"
	serialHandler.config = config
	serialHandler.connection = nil
	serialHandler.buf = make([]byte, 64)

	if err := serialHandler.Open(); err != nil {
		log.Fatalf(err.Error())
	}
	defer func() {
		log.Println("Closing connection")
		serialHandler.Close()
	}()

	serialChannel := make(chan string)
	serialLineChannel := make(chan string)
	go serialHandler.readSerial(serialChannel)
	go mergeSerialDataToLines(serialChannel, serialLineChannel)

	app := app{
		commandInput:    models.CreateCommandInput(),
		terminalDisplay: models.CreateTerminalDisplay(serialLineChannel),
	}

	if _, err := tea.NewProgram(app).Run(); err != nil {
		log.Fatal(err)
	}
}
