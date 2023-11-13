package style

import (
	"github.com/charmbracelet/lipgloss"
)

func col(x string) lipgloss.Color { return lipgloss.Color(x) }

var CommandFooter = func() lipgloss.Style {
	b := lipgloss.RoundedBorder()
	b.Right = "├"
	return lipgloss.NewStyle().BorderStyle(b).Padding(0,1)
}()

var InfoFooter = func() lipgloss.Style {
	b := lipgloss.RoundedBorder()
	b.Left = "┤"
	return CommandFooter.Copy().BorderStyle(b)
}()

var ViewFooter = func() lipgloss.Style {
	b := lipgloss.RoundedBorder()
	b.Right = "├"
	b.Left = "┤"
	return CommandFooter.Copy().BorderStyle(b)
}()

var MainView = func() lipgloss.Style {
	b := lipgloss.RoundedBorder()
	return lipgloss.NewStyle().BorderStyle(b).Margin(0, 10)
}()

var ConfItem = lipgloss.NewStyle()
var ConfItemLabel = lipgloss.NewStyle().Width(10)
var ConfItemValue = lipgloss.NewStyle().Width(10).Align(lipgloss.Right)

