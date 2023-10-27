package main

import (
	"fmt"
	"os"
	//"time"
	//"log"

	//"github.com/dsyx/serialport-go"
	tea "github.com/charmbracelet/bubbletea"
)

//type SerialConnection struct {
//	Device string
//	Config serialport.Config
//}

type model struct {
	choices []string
	cursor int
	selected map[int]struct{}
}

func initializeModel() model {
	return model{
		choices: []string{"Buy carrots", "Buy celery", "Buy kohlrabi"},
		selected: make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String(){
			case "ctrl+c", "q":
				return m, tea.Quit

			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if m.cursor < len(m.choices) - 1 {
					m.cursor++
				}
			case "enter", " ":
				_, ok := m.selected[m.cursor]
				if ok {
					delete(m.selected, m.cursor)
				} else {
					m.selected[m.cursor] = struct{}{}

				}
		}
	}

	return m, nil
}

func (m model) View() string { 
	s := "What should we buy at the market?\n\n"

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "x"
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}


	s += "\n"
	s += "Press q to quit.\n"

	return s
}

func main(){
	p := tea.NewProgram(initializeModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there has been an error: %v", err)
		os.Exit(1)
	}

	//conn := SerialConnection{
	//	Device: "COM8",
	//	Config: serialport.Config{
	//		BaudRate: serialport.BR9600,
	//		DataBits: serialport.DB8,
	//		StopBits: serialport.SB1,
	//		Parity: serialport.PO,
	//		Timeout: 100 * time.Millisecond,
	//	},
	//}
	//sp, err := serialport.Open(conn.Device, conn.Config)
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//defer sp.Close()


	//fmt.Println("Device is open")

	//buf := make([]byte, 1024)
	//for {
	//	n, _ := sp.Read(buf)
	//	fmt.Printf("%s", string(buf[:n]))
	//}
}
