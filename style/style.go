package style

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

var (
	Color = termenv.EnvColorProfile().Color
	Keyword = termenv.Style{}.Foreground(Color("204")).Background(Color("235")).Styled
	Help = termenv.Style{}.Foreground(Color("241")).Styled
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

var MainView = lipgloss.NewStyle()

func HelpLine() string {
	options := []string{
		"esc: exit",
		"enter: send",
		"tab: menu",
		"delete: clear",
	}
	sep := " • "
	return Help(strings.Join(options, sep))
}
