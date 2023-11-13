package main

import (
	"fmt"
	"strings"

	"rendellc/gtty/style"
)

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


