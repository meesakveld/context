package progress

import (
	"golang.org/x/term"
	"os"
)

func IsTerminal() bool {
	return term.IsTerminal(int(os.Stdout.Fd()))
}
