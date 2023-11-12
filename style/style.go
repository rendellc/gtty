package style

import (
	"github.com/charmbracelet/lipgloss"
)

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
	return lipgloss.NewStyle().BorderStyle(b)
}()
