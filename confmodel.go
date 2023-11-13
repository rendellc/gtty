package main

import (
	"fmt"
	"log"
	"strings"
	"strconv"

	"rendellc/gtty/style"
)


type configurator struct{}

func NewConfigurator() configurator {
	return configurator{}
}

func (c configurator) doSetCommand(params []string, config *appConfig) {
	if len(params) < 2 {
		return
	}

	target := params[0]
	value := params[1]
	valueInt, intParseErr := strconv.Atoi(value)
	switch {
	case target == "Device":
		log.Printf("Updating Device to " + value)
		config.SerialConfig.Device = value
	case target == "BaudRate" && intParseErr == nil:
		log.Printf("Updating BaudRate to " + value)
		config.SerialConfig.BaudRate = valueInt
	}
}

func (c configurator) DoCommand(cmd string, config *appConfig) {
	fields := strings.Fields(cmd)
	log.Printf("Command: %v", fields)

	if len(fields) == 0 {
		return
	}

	if fields[0] == "set" {
		c.doSetCommand(fields[1:], config)
	}
}


type ConfItem struct {
	Label string
	Getter func() string
	Setter func(string)
}

func renderConfItemStr(label, value string) string {
	return style.ConfItem.Render(
		fmt.Sprintf("%s%s", 
			style.ConfItemLabel.Render(label),
			style.ConfItemValue.Render(value),
		),
	)
}

func renderConfItemInt(label string, value int) string {
	return style.ConfItem.Render(
		fmt.Sprintf("%s%s", 
			style.ConfItemLabel.Render(label),
			style.ConfItemValue.Render(fmt.Sprintf("%d", value)),
		),
	)
}

func (c ConfItem) render() string {
	return style.ConfItem.Render(
		fmt.Sprintf("%s%s", 
			style.ConfItemLabel.Render(c.Label),
			style.ConfItemValue.Render(c.Getter()),
		),
	)
}

func viewConnectionConfig(config *appConfig) string {
	confItems := []ConfItem{
		ConfItem{
			Label: "Device",
			Getter: func() string {
				return config.SerialConfig.Device
			},
			Setter: func(value string) {
				config.SerialConfig.Device = value
			},
		},
	}
	out := strings.Builder{}
	
	for _, ci := range confItems {
		out.WriteString(ci.render())
		out.WriteString("\n")
	}
	out.WriteString(renderConfItemInt("Baud rate", config.SerialConfig.BaudRate))
	out.WriteString("\n")
	out.WriteString(renderConfItemInt("Data bits", config.SerialConfig.DataBits))
	out.WriteString("\n")
	out.WriteString(renderConfItemStr("Parity", config.SerialConfig.Parity))
	out.WriteString("\n")
	out.WriteString(renderConfItemInt("Stop bits", config.SerialConfig.StopBits))
	out.WriteString("\n")


	return out.String()
}


