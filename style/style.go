package style

import (
	"strings"

	"github.com/muesli/termenv"
)

var (
	Color = termenv.EnvColorProfile().Color
	Keyword = termenv.Style{}.Foreground(Color("204")).Background(Color("235")).Styled
	Help = termenv.Style{}.Foreground(Color("241")).Styled
)

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
